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

func TestNewInt16(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    int16
		p2    bool
		int16 int16
		valid bool
		set   bool
	}{
		{p1: 0, p2: false, int16: 0, valid: false, set: true},
		{p1: 0, p2: true, int16: 0, valid: true, set: true},
		{p1: 1, p2: false, int16: 1, valid: false, set: true},
		{p1: 1, p2: true, int16: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%d+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			n := null.NewInt16(tc.p1, tc.p2)
			key := fmt.Sprintf("null.NewInt16(%d, %t)", tc.p1, tc.p2)
			if n.Int16() != tc.int16 {
				t.Errorf("%s.Int16(): got %d, want %d", key, n.Int16(), tc.int16)
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

func TestNewInt16Null(t *testing.T) {
	n := null.NewInt16Null()

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

func TestInt16_Equal(t *testing.T) {
	cases := map[string]struct {
		got   null.Int16
		want  null.Int16
		equal bool
	}{
		"empty":                      {got: null.Int16{}, want: null.Int16{}, equal: true},
		"zero":                       {got: null.NewInt16(0), want: null.NewInt16(0), equal: true},
		"1 == 1":                     {got: null.NewInt16(1), want: null.NewInt16(1), equal: true},
		"0 != 1":                     {got: null.NewInt16(0), want: null.NewInt16(1), equal: false},
		"1 [invalid] != 1":           {got: null.NewInt16(1, false), want: null.NewInt16(1), equal: false},
		"1 [invalid] == 1 [invalid]": {got: null.NewInt16(1, false), want: null.NewInt16(1, false), equal: true},
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

func TestInt16_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    int16
		p2    bool
		bytes []byte
	}{
		{p1: 0, p2: false, bytes: []byte("null")},
		{p1: 0, p2: true, bytes: []byte(`0`)},
		{p1: 1, p2: false, bytes: []byte("null")},
		{p1: 1, p2: true, bytes: []byte(`1`)},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%d+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.NewInt16(%d, %t)", tc.p1, tc.p2)
			n := null.NewInt16(tc.p1, tc.p2)
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

func TestInt16_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		bytes []byte
		int16 int16
		valid bool
		set   bool
	}{
		{bytes: []byte("null"), int16: 0, valid: false, set: true},
		{bytes: []byte(`1`), int16: 1, valid: true, set: true},
		{bytes: []byte(`0`), int16: 0, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.bytes)
		t.Run(name, func(t *testing.T) {
			var n null.Int16
			err := n.UnmarshalJSON(tc.bytes)
			if err != nil {
				t.Fatal(err)
			}
			key := fmt.Sprintf("null.Int16.UnmarshalJSON(%v)", tc.bytes)
			if n.Int16() != tc.int16 {
				t.Errorf("%s.Int16(): got %q, want %q", key, n.Int16(), tc.int16)
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

func TestInt16_Scan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input interface{}
		int16 int16
		valid bool
		set   bool
	}{
		{input: nil, int16: 0, valid: false, set: true},
		{input: "0", int16: 0, valid: true, set: true},
		{input: "1", int16: 1, valid: true, set: true},
		{input: 0, int16: 0, valid: true, set: true},
		{input: 1, int16: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.input)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.Int16.Scan(%v)", tc.input)
			var n null.Int16
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if n.Int16() != tc.int16 {
				t.Errorf("%s.Int16(): got %q, want %q", key, n.Int16(), tc.int16)
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
				var n null.Int16
				err := n.Scan(b)
				if err != nil {
					t.Fatal(err)
				}
				if n.Int16() != tc.int16 {
					t.Errorf("%s.Int16(): got %q, want %q", key, n.Int16(), tc.int16)
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

func TestInt16_Value(t *testing.T) {
	cases := []struct {
		n   null.Int16
		val driver.Value
	}{
		{n: null.Int16{}, val: nil},
		{n: null.NewInt16(0, false), val: nil},
		{n: null.NewInt16(1, false), val: nil},
		{n: null.NewInt16(0, true), val: 0},
		{n: null.NewInt16(1, true), val: 1},
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
					t.Errorf("null.Int16.Val(): got %v, want %v", v, tc.val)
				}
			} else if reflect.DeepEqual(v, tc.val) {
				t.Errorf("null.Int16.Val(): got %v, want %v", v, tc.val)
			}
		})
	}
}
