package null_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bjaus/null"
)

func TestNewString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1     string
		p2     bool
		string string
		valid  bool
		set    bool
	}{
		{p1: "", p2: false, string: "", valid: false, set: true},
		{p1: "", p2: true, string: "", valid: true, set: true},
		{p1: "test", p2: false, string: "test", valid: false, set: true},
		{p1: "test", p2: true, string: "test", valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%q+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			n := null.NewString(tc.p1, tc.p2)
			key := fmt.Sprintf("null.NewString(%q, %t)", tc.p1, tc.p2)
			if n.String() != tc.string {
				t.Errorf("%s.String(): got %q, want %q", key, n.String(), tc.string)
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

func TestNewStringNull(t *testing.T) {
	n := null.NewStringNull()

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

func TestString_Equal(t *testing.T) {
	cases := map[string]struct {
		got   null.String
		want  null.String
		equal bool
	}{
		"empty":                          {got: null.String{}, want: null.String{}, equal: true},
		"abc":                            {got: null.NewString("abc"), want: null.NewString("abc"), equal: true},
		"abc != abc [invalid]":           {got: null.NewString("abc"), want: null.NewString("abc", false), equal: false},
		"abc [invalid] == abc [invalid]": {got: null.NewString("abc", false), want: null.NewString("abc", false), equal: true},
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

func TestString_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		p1    string
		p2    bool
		bytes []byte
	}{
		{p1: "", p2: false, bytes: []byte("null")},
		{p1: "", p2: true, bytes: []byte(`""`)},
		{p1: "test", p2: false, bytes: []byte("null")},
		{p1: "test", p2: true, bytes: []byte(`"test"`)},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%q+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.NewString(%q, %t)", tc.p1, tc.p2)
			nb := null.NewString(tc.p1, tc.p2)
			b, err := json.Marshal(nb)
			if err != nil {
				t.Fatal(err)
			}
			if num := bytes.Compare(b, tc.bytes); num != 0 {
				t.Errorf("%s.MarshalJSON(): got %s, want %s", key, b, tc.bytes)
			}
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		bytes  []byte
		string string
		valid  bool
		set    bool
	}{
		{bytes: []byte("null"), string: "", valid: false, set: true},
		{bytes: []byte(`"test"`), string: "test", valid: true, set: true},
		{bytes: []byte(`""`), string: "", valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.bytes)
		t.Run(name, func(t *testing.T) {
			var n null.String
			err := n.UnmarshalJSON(tc.bytes)
			if err != nil {
				t.Fatal(err)
			}
			key := fmt.Sprintf("null.String.UnmarshalJSON(%v)", tc.bytes)
			if n.String() != tc.string {
				t.Errorf("%s.String(): got %q, want %q", key, n.String(), tc.string)
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

func TestString_Scan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input  interface{}
		string string
		valid  bool
		set    bool
	}{
		{input: nil, string: "", valid: false, set: true},
		{input: "", string: "", valid: true, set: true},
		{input: "null", string: "null", valid: true, set: true},
		{input: 1, string: "1", valid: true, set: true},
		{input: 1.23, string: "1.23", valid: true, set: true},
		{input: true, string: "true", valid: true, set: true},
		{input: false, string: "false", valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.input)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.String.Scan(%v)", tc.input)
			var n null.String
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if n.String() != tc.string {
				t.Errorf("%s.String(): got %q, want %q", key, n.String(), tc.string)
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
				var n null.String
				err := n.Scan(b)
				if err != nil {
					t.Fatal(err)
				}
				if n.String() != tc.string {
					t.Errorf("%s.String(): got %q, want %q", key, n.String(), tc.string)
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

func TestString_Value(t *testing.T) {
	t.Parallel()

	cases := []struct {
		n   null.String
		val interface{}
	}{
		{n: null.String{}, val: nil},
		{n: null.NewString("", false), val: nil},
		{n: null.NewString("test", false), val: nil},
		{n: null.NewString("", true), val: ""},
		{n: null.NewString("test", true), val: "test"},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%+v", tc.n)
		t.Run(name, func(t *testing.T) {
			v, err := tc.n.Value()
			if err != nil {
				t.Fatal(err)
			}
			if v != tc.val {
				t.Errorf("null.String.Val(): got %v, want %v", v, tc.val)
			}
		})
	}
}
