package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

var _ Nuller = new(Int16)

// Int16 is a type that can be null or an Int16.
type Int16 struct {
	int16
	valid bool
	set   bool
}

// NewInt16 creates a new Int16 value.
func NewInt16(i int16, valid ...bool) Int16 {
	var v bool
	if len(valid) == 0 {
		v = true
	} else {
		v = valid[0]
	}
	return Int16{int16: i, valid: v, set: true}
}

func NewInt16Null() Int16 {
	return Int16{set: true}
}

func (n Int16) Equal(o Int16) bool {
	return n.set == o.set && n.valid == o.valid && n.int16 == o.int16
}

// Int16 is the underlying int16 value.
func (n Int16) Int16() int16 {
	return n.int16
}

// Valid implements the Nuller interface.
func (n Int16) Valid() bool {
	return n.valid
}

// Set implements the Nuller interface.
func (n Int16) Set() bool {
	return n.set
}

// IsNull indicates whether a value was explicitly null.
func (n Int16) Null() bool {
	return n.Set() && !n.Valid()
}

// MarshalJSON implements the json.Marshaler interface.
func (n Int16) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.int16)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Int16) UnmarshalJSON(b []byte) error {
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
func (n Int16) Value() (driver.Value, error) {
	return sql.NullInt16{Int16: n.int16, Valid: n.valid}.Value()
}

// Scan implements the sql.Scanner interface.
func (n *Int16) Scan(i interface{}) error {
	var ni sql.NullInt16
	if err := ni.Scan(i); err != nil {
		return err
	}
	*n = Int16{
		int16: ni.Int16,
		valid: ni.Valid,
		set:   true,
	}
	return nil
}
