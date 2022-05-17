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

func TestNewInt32(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    int32
		p2    bool
		int32 int32
		valid bool
		set   bool
	}{
		{p1: 0, p2: false, int32: 0, valid: false, set: true},
		{p1: 0, p2: true, int32: 0, valid: true, set: true},
		{p1: 1, p2: false, int32: 1, valid: false, set: true},
		{p1: 1, p2: true, int32: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%d+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			n := null.NewInt32(tc.p1, tc.p2)
			key := fmt.Sprintf("null.NewInt32(%d, %t)", tc.p1, tc.p2)
			if n.Int32() != tc.int32 {
				t.Errorf("%s.Int32(): got %d, want %d", key, n.Int32(), tc.int32)
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

func TestNewInt32Null(t *testing.T) {
	n := null.NewInt32Null()

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

func TestInt32_Equal(t *testing.T) {
	cases := map[string]struct {
		got   null.Int32
		want  null.Int32
		equal bool
	}{
		"empty":                      {got: null.Int32{}, want: null.Int32{}, equal: true},
		"zero":                       {got: null.NewInt32(0), want: null.NewInt32(0), equal: true},
		"1 == 1":                     {got: null.NewInt32(1), want: null.NewInt32(1), equal: true},
		"0 != 1":                     {got: null.NewInt32(0), want: null.NewInt32(1), equal: false},
		"1 [invalid] != 1":           {got: null.NewInt32(1, false), want: null.NewInt32(1), equal: false},
		"1 [invalid] == 1 [invalid]": {got: null.NewInt32(1, false), want: null.NewInt32(1, false), equal: true},
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

func TestInt32_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    int32
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
			key := fmt.Sprintf("null.NewInt32(%d, %t)", tc.p1, tc.p2)
			n := null.NewInt32(tc.p1, tc.p2)
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

func TestInt32_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		bytes []byte
		int32 int32
		valid bool
		set   bool
	}{
		{bytes: []byte("null"), int32: 0, valid: false, set: true},
		{bytes: []byte(`1`), int32: 1, valid: true, set: true},
		{bytes: []byte(`0`), int32: 0, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.bytes)
		t.Run(name, func(t *testing.T) {
			var n null.Int32
			err := n.UnmarshalJSON(tc.bytes)
			if err != nil {
				t.Fatal(err)
			}
			key := fmt.Sprintf("null.Int32.UnmarshalJSON(%v)", tc.bytes)
			if n.Int32() != tc.int32 {
				t.Errorf("%s.Int32(): got %q, want %q", key, n.Int32(), tc.int32)
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

func TestInt32_Scan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input interface{}
		int32 int32
		valid bool
		set   bool
	}{
		{input: nil, int32: 0, valid: false, set: true},
		{input: "0", int32: 0, valid: true, set: true},
		{input: "1", int32: 1, valid: true, set: true},
		{input: 0, int32: 0, valid: true, set: true},
		{input: 1, int32: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.input)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.Int32.Scan(%v)", tc.input)
			var n null.Int32
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if n.Int32() != tc.int32 {
				t.Errorf("%s.Int32(): got %q, want %q", key, n.Int32(), tc.int32)
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
				var n null.Int32
				err := n.Scan(b)
				if err != nil {
					t.Fatal(err)
				}
				if n.Int32() != tc.int32 {
					t.Errorf("%s.Int32(): got %q, want %q", key, n.Int32(), tc.int32)
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

func TestInt32_Value(t *testing.T) {
	cases := []struct {
		n   null.Int32
		val driver.Value
	}{
		{n: null.Int32{}, val: nil},
		{n: null.NewInt32(0, false), val: nil},
		{n: null.NewInt32(1, false), val: nil},
		{n: null.NewInt32(0, true), val: 0},
		{n: null.NewInt32(1, true), val: 1},
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
					t.Errorf("null.Int32.Val(): got %v, want %v", v, tc.val)
				}
			} else if reflect.DeepEqual(v, tc.val) {
				t.Errorf("null.Int32.Val(): got %v, want %v", v, tc.val)
			}
		})
	}
}
