package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

var _ Nuller = new(String)

// String is a type that can be null or a string.
type String struct {
	string
	valid bool
	set   bool
}

// NewString creates a new String value.
func NewString(s string, valid ...bool) String {
	var v bool
	if len(valid) == 0 {
		v = true
	} else {
		v = valid[0]
	}
	return String{string: s, valid: v, set: true}
}

func NewStringNull() String {
	return String{set: true}
}

func (n String) Equal(o String) bool {
	return n.set == o.set && n.valid == o.valid && n.string == o.string
}

// String is the underlying string value.
func (n String) String() string {
	return n.string
}

// Valid implements the Nuller interface.
func (n String) Valid() bool {
	return n.valid
}

// Set implements the Nuller interface.
func (n String) Set() bool {
	return n.set
}

// IsNull indicates whether a value was explicitly null.
func (n String) Null() bool {
	return n.Set() && !n.Valid()
}

// MarshalJSON implements the json.Marshaler interface.
func (n String) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.string)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *String) UnmarshalJSON(b []byte) error {
	n.set = true
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	return n.Scan(i)
}

// Value implements the driver.Valuer interface.
func (n String) Value() (driver.Value, error) {
	return sql.NullString{String: n.string, Valid: n.valid}.Value()
}

// Scan implements the sql.Scanner interface.
func (n *String) Scan(i interface{}) error {
	var ns sql.NullString
	if err := ns.Scan(i); err != nil {
		return err
	}
	*n = String{
		string: ns.String,
		valid:  ns.Valid,
		set:    true,
	}
	return nil
}
