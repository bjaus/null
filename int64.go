package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

var _ Nuller = new(Int64)

// Int64 is a type that can be null or an Int64.
type Int64 struct {
	int64
	valid bool
	set   bool
}

// NewInt64 creates an new Int64 value.
func NewInt64(i int64, valid ...bool) Int64 {
	var v bool
	if len(valid) == 0 {
		v = true
	} else {
		v = valid[0]
	}
	return Int64{int64: i, valid: v, set: true}
}

func NewInt64Null() Int64 {
	return Int64{set: true}
}

func (n Int64) Equal(o Int64) bool {
	return n.set == o.set && n.valid == o.valid && n.int64 == o.int64
}

// Int64 is the underlying int64 value.
func (n Int64) Int64() int64 {
	return n.int64
}

// Valid implements the Nuller interface.
func (n Int64) Valid() bool {
	return n.valid
}

// Set implements the Nuller interface.
func (n Int64) Set() bool {
	return n.set
}

// IsNull indicates whether a value was explicitly null.
func (n Int64) Null() bool {
	return n.Set() && !n.Valid()
}

// MarshalJSON implements the json.Marshaler interface.
func (n Int64) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Int64) UnmarshalJSON(b []byte) error {
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
func (n Int64) Value() (driver.Value, error) {
	return sql.NullInt64{Int64: n.int64, Valid: n.valid}.Value()
}

// Scan implements the sql.Scanner interface.
func (n *Int64) Scan(i interface{}) error {
	var ni sql.NullInt64
	if err := ni.Scan(i); err != nil {
		return err
	}
	*n = Int64{
		int64: ni.Int64,
		valid: ni.Valid,
		set:   true,
	}
	return nil
}
