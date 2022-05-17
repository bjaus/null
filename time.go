package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

var _ Nuller = new(Time)

// Time is a type that can be null or a time.
type Time struct {
	time  time.Time
	valid bool
	set   bool
}

// NewTime creates a new Time value.
func NewTime(t time.Time, valid ...bool) Time {
	var v bool
	if len(valid) == 0 {
		v = true
	} else {
		v = valid[0]
	}
	return Time{time: t, valid: v, set: true}
}

func NewTimeNull() Time {
	return Time{set: true}
}

func (n Time) Equal(o Time) bool {
	return n.set == o.set && n.valid == o.valid && n.time == o.time
}

// Time is the underlying time.Time value.
func (n Time) Time() time.Time {
	return n.time
}

// Valid implements the Nuller interface.
func (n Time) Valid() bool {
	return n.valid && !n.time.IsZero()
}

// Set implements the Nuller interface.
func (n Time) Set() bool {
	return n.set
}

// IsNull indicates whether a value was explicitly null.
func (n Time) Null() bool {
	return n.Set() && !n.Valid()
}

// MarshalJSON implements the json.Marshaler interface.
func (n Time) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON deserializes a Time from JSON.
func (n *Time) UnmarshalJSON(b []byte) error {
	n.set = true
	// scan for null
	if bytes.Equal(b, []byte("null")) {
		return n.Scan(nil)
	}
	// scan for JSON timestamp
	var t time.Time
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	return n.Scan(t)
}

// Value implements the driver driver.Valuer interface.
func (n Time) Value() (driver.Value, error) {
	if n.Valid() {
		return n.time, nil
	}
	return nil, nil
}

// Scan implements the sql.Scanner interface.
//
// The value type must be time.Time or string / []byte (formatted time-string), otherwise Scan fails.
func (n *Time) Scan(value interface{}) error {
	n.set = true
	var err error
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		if v.IsZero() {
			return nil
		}
		n.time, n.valid = v, true
		return nil
	case []byte:
		n.time, err = parseDateTime(string(v), time.UTC)
		n.valid = (err == nil)
		return err
	case string:
		n.time, err = parseDateTime(v, time.UTC)
		n.valid = (err == nil)
		return err
	}

	n.valid = false
	return nil
}

func parseDateTime(str string, loc *time.Location) (time.Time, error) {
	var t time.Time
	var err error

	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26:
		if str == base[:len(str)] {
			return t, err
		}
		const timeFormat = "2006-01-02 15:04:05.000000"
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		return t, fmt.Errorf("invalid timestamp string: %q", str)
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return t, err
}
