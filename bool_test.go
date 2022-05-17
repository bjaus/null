package null_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bjaus/null"
)

func TestNewBool(t *testing.T) {
	cases := []struct {
		p1    bool
		p2    bool
		bool  bool
		valid bool
		set   bool
	}{
		{p1: false, p2: false, bool: false, valid: false, set: true},
		{p1: true, p2: false, bool: true, valid: false, set: true},
		{p1: false, p2: true, bool: false, valid: true, set: true},
		{p1: true, p2: true, bool: true, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%t+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			n := null.NewBool(tc.p1, tc.p2)
			key := fmt.Sprintf("null.NewBool(%t, %t)", tc.p1, tc.p2)
			if n.Bool() != tc.bool {
				t.Errorf("%s.Bool(): got %t, want %t", key, n.Bool(), tc.bool)
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

func TestNewBoolNull(t *testing.T) {
	n := null.NewBoolNull()

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

func TestBool_Equal(t *testing.T) {
	cases := map[string]struct {
		got   null.Bool
		want  null.Bool
		equal bool
	}{
		"empty":                         {got: null.Bool{}, want: null.Bool{}, equal: true},
		"false [valid]":                 {got: null.NewBool(false), want: null.NewBool(false), equal: true},
		"false [invalid]":               {got: null.NewBool(false, false), want: null.NewBool(false, false), equal: true},
		"true [valid]":                  {got: null.NewBool(true), want: null.NewBool(true), equal: true},
		"true [invalid]":                {got: null.NewBool(true, true), want: null.NewBool(true, true), equal: true},
		"false true":                    {got: null.NewBool(false), want: null.NewBool(true), equal: false},
		"true false":                    {got: null.NewBool(true), want: null.NewBool(false), equal: false},
		"false [invalid] false [valid]": {got: null.NewBool(false, false), want: null.NewBool(false), equal: false},
		"true [invalid] true [valid]":   {got: null.NewBool(true, false), want: null.NewBool(true), equal: false},
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

func TestBool_MarshalJSON(t *testing.T) {
	cases := []struct {
		p1    bool
		p2    bool
		bytes []byte
	}{
		{p1: false, p2: false, bytes: []byte("null")},
		{p1: false, p2: true, bytes: []byte("false")},
		{p1: true, p2: false, bytes: []byte("null")},
		{p1: true, p2: true, bytes: []byte("true")},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%t+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.NewBool(%t, %t)", tc.p1, tc.p2)
			n := null.NewBool(tc.p1, tc.p2)
			b, err := json.Marshal(n)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, tc.bytes) {
				t.Errorf("%s.MarshalJSON(): got %s, want %s", key, b, tc.bytes)
			}
		})
	}
}

func TestBool_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		bytes []byte
		bool  bool
		valid bool
		set   bool
	}{
		{bytes: []byte("null"), bool: false, valid: false, set: true},
		{bytes: []byte("false"), bool: false, valid: true, set: true},
		{bytes: []byte("true"), bool: true, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.bytes)
		t.Run(name, func(t *testing.T) {
			var n null.Bool
			err := n.UnmarshalJSON(tc.bytes)
			if err != nil {
				t.Fatal(err)
			}
			key := fmt.Sprintf("null.Bool.UnmarshalJSON(%v)", tc.bytes)
			if n.Bool() != tc.bool {
				t.Errorf("%s.Bool(): got %t, want %t", key, n.Bool(), tc.bool)
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

func TestBool_Scan(t *testing.T) {
	cases := []struct {
		input interface{}
		bool  bool
		valid bool
		set   bool
	}{
		{input: nil, bool: false, valid: false, set: true},
		{input: false, bool: false, valid: true, set: true},
		{input: true, bool: true, valid: true, set: true},
		{input: 0, bool: false, valid: true, set: true},
		{input: 1, bool: true, valid: true, set: true},
		{input: "0", bool: false, valid: true, set: true},
		{input: "1", bool: true, valid: true, set: true},
		{input: "FALSE", bool: false, valid: true, set: true},
		{input: "False", bool: false, valid: true, set: true},
		{input: "T", bool: true, valid: true, set: true},
		{input: "TRUE", bool: true, valid: true, set: true},
		{input: "True", bool: true, valid: true, set: true},
		{input: "f", bool: false, valid: true, set: true},
		{input: "f", bool: false, valid: true, set: true},
		{input: "false", bool: false, valid: true, set: true},
		{input: "t", bool: true, valid: true, set: true},
		{input: "true", bool: true, valid: true, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.input)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.Bool.Scan(%v)", tc.input)
			var n null.Bool
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if n.Bool() != tc.bool {
				t.Errorf("%s.Bool(): got %t, want %t", key, n.Bool(), tc.bool)
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
				var n null.Bool
				err := n.Scan(b)
				if err != nil {
					t.Fatal(err)
				}
				if n.Bool() != tc.bool {
					t.Errorf("%s.Bool(): got %t, want %t", key, n.Bool(), tc.bool)
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

func TestBool_Value(t *testing.T) {
	cases := []struct {
		n   null.Bool
		val interface{}
	}{
		{n: null.Bool{}, val: nil},
		{n: null.NewBool(false, false), val: nil},
		{n: null.NewBool(true, false), val: nil},
		{n: null.NewBool(false, true), val: false},
		{n: null.NewBool(true, true), val: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%+v", tc.n)
		t.Run(name, func(t *testing.T) {
			v, err := tc.n.Value()
			if err != nil {
				t.Fatal(err)
			}
			if v != tc.val {
				t.Errorf("null.Bool.Val(): got %v, want %v", v, tc.val)
			}
		})
	}
}
