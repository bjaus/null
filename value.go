package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// State represents the three possible states of a Value.
type State uint8

const (
	// Unset indicates the Value was never set (zero value).
	// This is the default state for uninitialized Values.
	Unset State = iota

	// Null indicates the Value was explicitly set to null.
	Null

	// Valid indicates the Value contains a valid value.
	Valid
)

// String returns a string representation of the state.
func (s State) String() string {
	switch s {
	case Unset:
		return "unset"
	case Null:
		return "null"
	case Valid:
		return "valid"
	default:
		return fmt.Sprintf("State(%d)", s)
	}
}

// Value represents an optional value that distinguishes between
// unset (absent), null (explicit null), and valid (has value).
//
// The zero value is Unset, meaning the field was never set.
type Value[T any] struct {
	v     T
	state State
}

// --- Constructors ---

// New creates a valid Value containing v.
func New[T any](v T) Value[T] {
	return Value[T]{v: v, state: Valid}
}

// NewNull creates a Value that is explicitly null.
func NewNull[T any]() Value[T] {
	return Value[T]{state: Null}
}

// NewPtr creates a Value from a pointer.
// If p is nil, returns a null Value. Otherwise returns a valid Value with *p.
func NewPtr[T any](p *T) Value[T] {
	if p == nil {
		return NewNull[T]()
	}
	return New(*p)
}

// --- State Queries ---

// IsSet returns true if the Value was explicitly set (either to null or a value).
// Returns false only for the zero value (Unset state).
func (v Value[T]) IsSet() bool {
	return v.state != Unset
}

// IsNull returns true if the Value was explicitly set to null.
func (v Value[T]) IsNull() bool {
	return v.state == Null
}

// IsValid returns true if the Value contains a valid value.
func (v Value[T]) IsValid() bool {
	return v.state == Valid
}

// State returns the current state (Unset, Null, or Valid).
func (v Value[T]) State() State {
	return v.state
}

// --- Value Extraction ---

// Get returns the underlying value.
// Returns the zero value of T if not valid.
func (v Value[T]) Get() T {
	if v.state == Valid {
		return v.v
	}
	var zero T
	return zero
}

// GetOr returns the underlying value if valid, otherwise returns def.
func (v Value[T]) GetOr(def T) T {
	if v.state == Valid {
		return v.v
	}
	return def
}

// Ptr returns a pointer to the value if valid, otherwise nil.
func (v Value[T]) Ptr() *T {
	if v.state == Valid {
		return &v.v
	}
	return nil
}

// --- JSON ---

var nullBytes = []byte("null")

// MarshalJSON implements json.Marshaler.
// Valid values are marshaled as their JSON representation.
// Null and unset values are marshaled as null.
func (v Value[T]) MarshalJSON() ([]byte, error) {
	if v.state == Valid {
		return json.Marshal(v.v)
	}
	return nullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// If the JSON value is null, the Value becomes Null.
// Otherwise, the Value becomes Valid with the unmarshaled value.
//
// Note: This method is only called when the field is present in JSON.
// If the field is absent, the Value remains Unset (zero value).
func (v *Value[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*v = NewNull[T]()
		return nil
	}
	var val T
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*v = New(val)
	return nil
}

// --- SQL driver.Valuer ---

// Value implements driver.Valuer for SQL database operations.
// Returns nil for null and unset values.
func (v Value[T]) Value() (driver.Value, error) {
	if v.state != Valid {
		return nil, nil
	}
	val := any(v.v)
	switch x := val.(type) {
	case string, int64, float64, bool, []byte, time.Time:
		return x, nil
	case int:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case int16:
		return int64(x), nil
	case int8:
		return int64(x), nil
	case uint:
		return int64(x), nil
	case uint64:
		return int64(x), nil
	case uint32:
		return int64(x), nil
	case uint16:
		return int64(x), nil
	case uint8:
		return int64(x), nil
	case float32:
		return float64(x), nil
	default:
		return val, nil
	}
}

// --- SQL sql.Scanner ---

// Scan implements sql.Scanner for SQL database operations.
// A nil source results in a Null value.
func (v *Value[T]) Scan(src any) error {
	if src == nil {
		*v = NewNull[T]()
		return nil
	}

	target := any(&v.v)

	switch ptr := target.(type) {
	case *string:
		if err := scanString(ptr, src); err != nil {
			return err
		}
	case *int64:
		if err := scanInt64(ptr, src); err != nil {
			return err
		}
	case *int:
		var i int64
		if err := scanInt64(&i, src); err != nil {
			return err
		}
		*ptr = int(i)
	case *int32:
		var i int64
		if err := scanInt64(&i, src); err != nil {
			return err
		}
		*ptr = int32(i)
	case *int16:
		var i int64
		if err := scanInt64(&i, src); err != nil {
			return err
		}
		*ptr = int16(i)
	case *int8:
		var i int64
		if err := scanInt64(&i, src); err != nil {
			return err
		}
		*ptr = int8(i)
	case *uint:
		var u uint64
		if err := scanUint64(&u, src); err != nil {
			return err
		}
		*ptr = uint(u)
	case *uint64:
		if err := scanUint64(ptr, src); err != nil {
			return err
		}
	case *uint32:
		var u uint64
		if err := scanUint64(&u, src); err != nil {
			return err
		}
		*ptr = uint32(u)
	case *uint16:
		var u uint64
		if err := scanUint64(&u, src); err != nil {
			return err
		}
		*ptr = uint16(u)
	case *uint8:
		var u uint64
		if err := scanUint64(&u, src); err != nil {
			return err
		}
		*ptr = uint8(u)
	case *float64:
		if err := scanFloat64(ptr, src); err != nil {
			return err
		}
	case *float32:
		var f float64
		if err := scanFloat64(&f, src); err != nil {
			return err
		}
		*ptr = float32(f)
	case *bool:
		if err := scanBool(ptr, src); err != nil {
			return err
		}
	case *time.Time:
		if err := scanTime(ptr, src); err != nil {
			return err
		}
	case *[]byte:
		if err := scanBytes(ptr, src); err != nil {
			return err
		}
	default:
		if err := scanReflect(&v.v, src); err != nil {
			return err
		}
	}

	v.state = Valid
	return nil
}

func scanString(dst *string, src any) error {
	switch s := src.(type) {
	case string:
		*dst = s
	case []byte:
		*dst = string(s)
	default:
		return fmt.Errorf("null: cannot scan %T into string", src)
	}
	return nil
}

func scanInt64(dst *int64, src any) error {
	switch s := src.(type) {
	case int64:
		*dst = s
	case int:
		*dst = int64(s)
	case int32:
		*dst = int64(s)
	case float64:
		*dst = int64(s)
	default:
		return fmt.Errorf("null: cannot scan %T into int64", src)
	}
	return nil
}

func scanUint64(dst *uint64, src any) error {
	switch s := src.(type) {
	case uint64:
		*dst = s
	case uint:
		*dst = uint64(s)
	case uint32:
		*dst = uint64(s)
	case int64:
		*dst = uint64(s)
	case float64:
		*dst = uint64(s)
	default:
		return fmt.Errorf("null: cannot scan %T into uint64", src)
	}
	return nil
}

func scanFloat64(dst *float64, src any) error {
	switch s := src.(type) {
	case float64:
		*dst = s
	case float32:
		*dst = float64(s)
	case int64:
		*dst = float64(s)
	default:
		return fmt.Errorf("null: cannot scan %T into float64", src)
	}
	return nil
}

func scanBool(dst *bool, src any) error {
	switch s := src.(type) {
	case bool:
		*dst = s
	case int64:
		*dst = s != 0
	default:
		return fmt.Errorf("null: cannot scan %T into bool", src)
	}
	return nil
}

func scanTime(dst *time.Time, src any) error {
	switch s := src.(type) {
	case time.Time:
		*dst = s
	case string:
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return fmt.Errorf("null: cannot parse time %q: %w", s, err)
		}
		*dst = t
	case []byte:
		t, err := time.Parse(time.RFC3339, string(s))
		if err != nil {
			return fmt.Errorf("null: cannot parse time %q: %w", s, err)
		}
		*dst = t
	default:
		return fmt.Errorf("null: cannot scan %T into time.Time", src)
	}
	return nil
}

func scanBytes(dst *[]byte, src any) error {
	switch s := src.(type) {
	case []byte:
		*dst = s
	case string:
		*dst = []byte(s)
	default:
		return fmt.Errorf("null: cannot scan %T into []byte", src)
	}
	return nil
}

func scanReflect(dst any, src any) error {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src)

	if srcVal.Type().AssignableTo(dstVal.Type()) {
		dstVal.Set(srcVal)
		return nil
	}
	if srcVal.Type().ConvertibleTo(dstVal.Type()) {
		dstVal.Set(srcVal.Convert(dstVal.Type()))
		return nil
	}
	return fmt.Errorf("null: cannot scan %T into %T", src, dst)
}
