package null_test

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/bjaus/null"
)

func TestNewFloat64(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1      float64
		p2      bool
		float64 float64
		valid   bool
		set     bool
	}{
		{p1: 0, p2: false, float64: 0, valid: false, set: true},
		{p1: 0, p2: true, float64: 0, valid: true, set: true},
		{p1: 1, p2: false, float64: 1, valid: false, set: true},
		{p1: 1, p2: true, float64: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%f+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			n := null.NewFloat64(tc.p1, tc.p2)
			key := fmt.Sprintf("null.NewFloat64(%f, %t)", tc.p1, tc.p2)
			if n.Float64() != tc.float64 {
				t.Errorf("%s.Float64(): got %f, want %f", key, n.Float64(), tc.float64)
			}
			if n.Valid() != tc.valid {
				t.Errorf("%s.Valid(): got %t, want %t", key, n.Valid(), tc.valid)
			}
			if n.Set() != tc.set {
				t.Errorf("%s.Set(): got %t, want %t", key, n.Set(), tc.set)
			}
		})
	}
}

func TestNewFloat64Null(t *testing.T) {
	n := null.NewFloat64Null()

	if !n.Set() {
		t.Error("should be set")
	}
	if !n.Null() {
		t.Error("should be null")
	}
	if n.Valid() {
		t.Error("should not be valid")
	}
}

func TestFloat64_Equal(t *testing.T) {
	cases := map[string]struct {
		got   null.Float64
		want  null.Float64
		equal bool
	}{
		"empty": {got: null.Float64{}, want: null.Float64{}, equal: true},
		"zero":  {got: null.NewFloat64(0), want: null.NewFloat64(0), equal: true},
		"1.234": {got: null.NewFloat64(1.23456789), want: null.NewFloat64(1.23467890123456789), equal: true},
		"2.34":  {got: null.NewFloat64(2.3456789), want: null.NewFloat64(2.3467890), equal: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := tc.got.Equal(tc.want)
			if got != tc.equal {
				t.Errorf("%q equal: got %t, want %t", name, got, tc.equal)
			}
		})
	}
}
func TestFloat64_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    float64
		p2    bool
		bytes []byte
	}{
		{p1: 0, p2: false, bytes: []byte("null")},
		{p1: 0, p2: true, bytes: []byte(`0`)},
		{p1: 1, p2: false, bytes: []byte("null")},
		{p1: 1, p2: true, bytes: []byte(`1`)},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%f+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.NewFloat64(%f, %t)", tc.p1, tc.p2)
			n := null.NewFloat64(tc.p1, tc.p2)
			b, err := json.Marshal(n)
			if err != nil {
				t.Fatal(err)
			}
			if num := bytes.Compare(b, tc.bytes); num != 0 {
				t.Errorf("%s.MarshalJSON(): got %s, want %s", key, b, tc.bytes)
			}
		})
	}
}

func TestFloat64_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		bytes   []byte
		float64 float64
		valid   bool
		set     bool
	}{
		{bytes: []byte("null"), float64: 0, valid: false, set: true},
		{bytes: []byte(`1`), float64: 1, valid: true, set: true},
		{bytes: []byte(`0`), float64: 0, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.bytes)
		t.Run(name, func(t *testing.T) {
			var n null.Float64
			err := n.UnmarshalJSON(tc.bytes)
			if err != nil {
				t.Fatal(err)
			}
			key := fmt.Sprintf("null.Float64.UnmarshalJSON(%v)", tc.bytes)
			if n.Float64() != tc.float64 {
				t.Errorf("%s.Float64(): got %f, want %f", key, n.Float64(), tc.float64)
			}
			if n.Valid() != tc.valid {
				t.Errorf("%s.Valid(): got %t, want %t", key, n.Valid(), tc.valid)
			}
			if n.Set() != tc.set {
				t.Errorf("%s.Set(): got %t, want %t", key, n.Set(), tc.set)
			}
		})
	}
}

func TestFloat64_Scan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input   interface{}
		float64 float64
		valid   bool
		set     bool
	}{
		{input: nil, float64: 0, valid: false, set: true},
		{input: "0", float64: 0, valid: true, set: true},
		{input: "1", float64: 1, valid: true, set: true},
		{input: 0, float64: 0, valid: true, set: true},
		{input: 1, float64: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.input)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.Float64.Scan(%v)", tc.input)
			var n null.Float64
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if n.Float64() != tc.float64 {
				t.Errorf("%s.Float64(): got %f, want %f", key, n.Float64(), tc.float64)
			}
			if n.Valid() != tc.valid {
				t.Errorf("%s.Valid(): got %t, want %t", key, n.Valid(), tc.valid)
			}
			if n.Set() != tc.set {
				t.Errorf("%s.Set(): got %t, want %t", key, n.Set(), tc.set)
			}

			// Easy way to also test []byte
			switch v := tc.input.(type) {
			case string:
				b := []byte(v)
				var n null.Float64
				err := n.Scan(b)
				if err != nil {
					t.Fatal(err)
				}
				if n.Float64() != tc.float64 {
					t.Errorf("%s.Float64(): got %f, want %f", key, n.Float64(), tc.float64)
				}
				if n.Valid() != tc.valid {
					t.Errorf("%s.Valid(): got %t, want %t", key, n.Valid(), tc.valid)
				}
				if n.Set() != tc.set {
					t.Errorf("%s.Set(): got %t, want %t", key, n.Set(), tc.set)
				}
			}
		})
	}
}

func TestFloat64_Value(t *testing.T) {
	cases := []struct {
		n   null.Float64
		val driver.Value
	}{
		{n: null.Float64{}, val: nil},
		{n: null.NewFloat64(0, false), val: nil},
		{n: null.NewFloat64(1, false), val: nil},
		{n: null.NewFloat64(0, true), val: 0},
		{n: null.NewFloat64(1, true), val: 1},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%+v", tc.n)
		t.Run(name, func(t *testing.T) {
			v, err := tc.n.Value()
			if err != nil {
				t.Fatal(err)
			}
			if tc.val == nil {
				if v != tc.val {
					t.Errorf("null.Float64.Val(): got %v, want %v", v, tc.val)
				}
			} else if reflect.DeepEqual(v, tc.val) {
				t.Errorf("null.Float64.Val(): got %v, want %v", v, tc.val)
			}
		})
	}
}
