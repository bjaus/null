package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"math"
)

var _ Nuller = new(Float64)

// Float64 is a type that can be null or a float64.
type Float64 struct {
	float64
	valid bool
	set   bool
}

// NewFloat64 creates a new Float64 value.
func NewFloat64(f float64, valid ...bool) Float64 {
	var v bool
	if len(valid) == 0 {
		v = true
	} else {
		v = valid[0]
	}
	return Float64{float64: f, valid: v, set: true}
}

func NewFloat64Null() Float64 {
	return Float64{set: true}
}

func (n Float64) Equal(o Float64) bool {
	result := n.set == o.set && n.valid == o.valid
	if !result {
		return false
	}
	tolerance := 0.001
	diff := math.Abs(n.float64 - o.float64)
	return diff < tolerance
}

// Float64 is the underlying float64 value.
func (n Float64) Float64() float64 {
	return n.float64
}

// Valid implements the Nuller interface.
func (n Float64) Valid() bool {
	return n.valid
}

// Set implements the Nuller interface.
func (n Float64) Set() bool {
	return n.set
}

// IsNull indicates whether a value was explicitly null.
func (n Float64) Null() bool {
	return n.Set() && !n.Valid()
}

// MarshalJSON implements the json.Marshaler interface.
func (n Float64) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Float64) UnmarshalJSON(b []byte) error {
	n.set = true
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	return n.Scan(i)
}

// Value implements the driver.Valuer interface.
func (n Float64) Value() (driver.Value, error) {
	return sql.NullFloat64{Float64: n.float64, Valid: n.valid}.Value()
}

// Scan implements the sql.Scanner interface.
func (n *Float64) Scan(i interface{}) error {
	var nf sql.NullFloat64
	if err := nf.Scan(i); err != nil {
		return err
	}
	*n = Float64{
		float64: nf.Float64,
		valid:   nf.Valid,
		set:     true,
	}
	return nil
}
