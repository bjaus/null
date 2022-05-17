package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

var _ Nuller = new(Bool)

// Bool is a type that can be "null" or a bool.
type Bool struct {
	bool
	valid bool
	set   bool
}

// NewBool creates a new Bool value.
func NewBool(b bool, valid ...bool) Bool {
	var v bool
	if len(valid) == 0 {
		v = true
	} else {
		v = valid[0]
	}
	return Bool{bool: b, valid: v, set: true}
}

func NewBoolNull() Bool {
	return Bool{set: true}
}

func (n Bool) Equal(o Bool) bool {
	return n.set == o.set && n.valid == o.valid && n.bool == o.bool
}

// Bool is the underlying bool value.
func (n Bool) Bool() bool {
	return n.bool
}

// Valid implements the Nuller interface.
func (n Bool) Valid() bool {
	return n.valid
}

// Set implements the Nuller interface.
func (n Bool) Set() bool {
	return n.set
}

// IsNull indicates whether a value was explicitly null.
func (n Bool) Null() bool {
	return n.Set() && !n.Valid()
}

// MarshalJSON implements the json.Marshaler interface.
func (n Bool) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Bool) UnmarshalJSON(b []byte) error {
	n.set = true
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	return n.Scan(i)
}

// Value implements the driver.Valuer interface.
func (n Bool) Value() (driver.Value, error) {
	return sql.NullBool{Bool: n.bool, Valid: n.valid}.Value()
}

// Scan implements the sql.Scanner interface.
func (n *Bool) Scan(i interface{}) error {
	var nb sql.NullBool
	if err := nb.Scan(i); err != nil {
		return err
	}
	*n = Bool{
		bool:  nb.Bool,
		valid: nb.Valid,
		set:   true,
	}
	return nil
}
