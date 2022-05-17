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

func TestNewInt64(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    int64
		p2    bool
		int64 int64
		valid bool
		set   bool
	}{
		{p1: 0, p2: false, int64: 0, valid: false, set: true},
		{p1: 0, p2: true, int64: 0, valid: true, set: true},
		{p1: 1, p2: false, int64: 1, valid: false, set: true},
		{p1: 1, p2: true, int64: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%d+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			n := null.NewInt64(tc.p1, tc.p2)
			key := fmt.Sprintf("null.NewInt64(%d, %t)", tc.p1, tc.p2)
			if n.Int64() != tc.int64 {
				t.Errorf("%s.Int64(): got %d, want %d", key, n.Int64(), tc.int64)
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

func TestNewInt64Null(t *testing.T) {
	n := null.NewInt64Null()

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

func TestInt64_Equal(t *testing.T) {
	cases := map[string]struct {
		got   null.Int64
		want  null.Int64
		equal bool
	}{
		"empty":                      {got: null.Int64{}, want: null.Int64{}, equal: true},
		"zero":                       {got: null.NewInt64(0), want: null.NewInt64(0), equal: true},
		"1 == 1":                     {got: null.NewInt64(1), want: null.NewInt64(1), equal: true},
		"0 != 1":                     {got: null.NewInt64(0), want: null.NewInt64(1), equal: false},
		"1 [invalid] != 1":           {got: null.NewInt64(1, false), want: null.NewInt64(1), equal: false},
		"1 [invalid] == 1 [invalid]": {got: null.NewInt64(1, false), want: null.NewInt64(1, false), equal: true},
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

func TestInt64_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    int64
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
			key := fmt.Sprintf("null.NewInt64(%d, %t)", tc.p1, tc.p2)
			n := null.NewInt64(tc.p1, tc.p2)
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

func TestInt64_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		bytes []byte
		int64 int64
		valid bool
		set   bool
	}{
		{bytes: []byte("null"), int64: 0, valid: false, set: true},
		{bytes: []byte(`1`), int64: 1, valid: true, set: true},
		{bytes: []byte(`0`), int64: 0, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.bytes)
		t.Run(name, func(t *testing.T) {
			var n null.Int64
			err := n.UnmarshalJSON(tc.bytes)
			if err != nil {
				t.Fatal(err)
			}
			key := fmt.Sprintf("null.Int64.UnmarshalJSON(%v)", tc.bytes)
			if n.Int64() != tc.int64 {
				t.Errorf("%s.Int64(): got %q, want %q", key, n.Int64(), tc.int64)
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

func TestInt64_Scan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input interface{}
		int64 int64
		valid bool
		set   bool
	}{
		{input: nil, int64: 0, valid: false, set: true},
		{input: "0", int64: 0, valid: true, set: true},
		{input: "1", int64: 1, valid: true, set: true},
		{input: 0, int64: 0, valid: true, set: true},
		{input: 1, int64: 1, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.input)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.Int64.Scan(%v)", tc.input)
			var n null.Int64
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if n.Int64() != tc.int64 {
				t.Errorf("%s.Int64(): got %q, want %q", key, n.Int64(), tc.int64)
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
				var n null.Int64
				err := n.Scan(b)
				if err != nil {
					t.Fatal(err)
				}
				if n.Int64() != tc.int64 {
					t.Errorf("%s.Int64(): got %q, want %q", key, n.Int64(), tc.int64)
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

func TestInt64_Value(t *testing.T) {
	cases := []struct {
		n   null.Int64
		val driver.Value
	}{
		{n: null.Int64{}, val: nil},
		{n: null.NewInt64(0, false), val: nil},
		{n: null.NewInt64(1, false), val: nil},
		{n: null.NewInt64(0, true), val: 0},
		{n: null.NewInt64(1, true), val: 1},
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
					t.Errorf("null.Int64.Val(): got %v, want %v", v, tc.val)
				}
			} else if reflect.DeepEqual(v, tc.val) {
				t.Errorf("null.Int64.Val(): got %v, want %v", v, tc.val)
			}
		})
	}
}
