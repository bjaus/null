package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

var _ Nuller = new(Int32)

// Int32 is a type that can be null or an Int32.
type Int32 struct {
	int32
	valid bool
	set   bool
}

// NewInt32 creates a new Int32 value.
func NewInt32(i int32, valid ...bool) Int32 {
	var v bool
	if len(valid) == 0 {
		v = true
	} else {
		v = valid[0]
	}
	return Int32{int32: i, valid: v, set: true}
}

func NewInt32Null() Int32 {
	return Int32{set: true}
}

func (n Int32) Equal(o Int32) bool {
	return n.set == o.set && n.valid == o.valid && n.int32 == o.int32
}

// Int32 is the underlying int32 value.
func (n Int32) Int32() int32 {
	return n.int32
}

// Valid implements the Nuller interface.
func (n Int32) Valid() bool {
	return n.valid
}

// Set implements the Nuller interface.
func (n Int32) Set() bool {
	return n.set
}

// IsNull indicates whether a value was explicitly null.
func (n Int32) Null() bool {
	return n.Set() && !n.Valid()
}

// MarshalJSON implements the json.Marshaler interface.
func (n Int32) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Int32) UnmarshalJSON(b []byte) error {
	n.set = true
	var i json.Number
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	if i == "" {
		return n.Scan(nil)
	}
	return n.Scan(i)
}

// Value implements the driver.Valuer interface.
func (n Int32) Value() (driver.Value, error) {
	return sql.NullInt32{Int32: n.int32, Valid: n.valid}.Value()
}

// Scan implements the sql.Scanner interface.
func (n *Int32) Scan(i interface{}) error {
	var ni sql.NullInt32
	if err := ni.Scan(i); err != nil {
		return err
	}
	*n = Int32{
		int32: ni.Int32,
		valid: ni.Valid,
		set:   true,
	}
	return nil
}
