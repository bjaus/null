package null_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bjaus/null"
)

func TestNewTime(t *testing.T) {
	now := time.Now().UTC()

	cases := []struct {
		p1    time.Time
		p2    bool
		time  time.Time
		valid bool
		set   bool
	}{
		{p1: now, p2: false, time: now, valid: false, set: true},
		{p1: now, p2: true, time: now, valid: true, set: true},
		{p1: time.Time{}, p2: false, time: time.Time{}, valid: false, set: true},
		{p1: time.Time{}, p2: true, time: time.Time{}, valid: false, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			n := null.NewTime(tc.p1, tc.p2)
			key := fmt.Sprintf("null.NewTime(%v, %t)", tc.p1, tc.p2)
			if n.Time() != tc.time {
				t.Errorf("%s.Time(): got %v, want %v", key, n.Time(), tc.time)
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

func TestNewTimeNull(t *testing.T) {
	n := null.NewTimeNull()

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

func TestTime_Equal(t *testing.T) {
	now := time.Now().UTC()
	<-time.After(time.Microsecond * 3)
	then := time.Now().UTC()

	cases := map[string]struct {
		got   null.Time
		want  null.Time
		equal bool
	}{
		"empty":                          {got: null.Time{}, want: null.Time{}, equal: true},
		"now":                            {got: null.NewTime(now), want: null.NewTime(now), equal: true},
		"now != now [invalid]":           {got: null.NewTime(now), want: null.NewTime(now, false), equal: false},
		"now [invalid] == now [invalid]": {got: null.NewTime(now, false), want: null.NewTime(now, false), equal: true},
		"not != then":                    {got: null.NewTime(now), want: null.NewTime(then), equal: false},
		"same but different instances": {
			got:   null.NewTime(time.Date(1988, time.December, 16, 12, 34, 56, 123456, time.UTC)),
			want:  null.NewTime(time.Date(1988, time.December, 16, 12, 34, 56, 123456, time.UTC)),
			equal: true,
		},
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

func TestTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	nowfmt := fmt.Sprintf("%q", now.Format(time.RFC3339Nano))

	cases := []struct {
		p1    time.Time
		p2    bool
		bytes []byte
	}{
		{p1: time.Time{}, p2: false, bytes: []byte("null")},
		{p1: time.Time{}, p2: true, bytes: []byte("null")},
		{p1: now, p2: false, bytes: []byte("null")},
		{p1: now, p2: true, bytes: []byte(nowfmt)},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v+%t", tc.p1, tc.p2)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.NewTime(%v, %t)", tc.p1, tc.p2)
			n := null.NewTime(tc.p1, tc.p2)
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

func TestTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	// TODO: this could still use some work

	now := time.Now().UTC()
	nowfmt := fmt.Sprintf("%q", now.Format(time.RFC3339Nano))

	cases := []struct {
		bytes []byte
		time  time.Time
		valid bool
		set   bool
	}{
		{bytes: []byte("null"), time: time.Time{}, valid: false, set: true},
		{bytes: []byte(nowfmt), time: now, valid: true, set: true},
	}

	for _, tc := range cases {
		name := string(tc.bytes)
		t.Run(name, func(t *testing.T) {
			var n null.Time
			err := n.UnmarshalJSON(tc.bytes)
			if err != nil {
				t.Fatal(err)
			}
			key := fmt.Sprintf("null.Time.UnmarshalJSON(%s)", tc.bytes)
			if n.Time() != tc.time {
				t.Errorf("%s.Time(): got %v, want %v", key, n.Time(), tc.time)
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

func TestTime_Scan(t *testing.T) {
	t.Parallel()

	// TODO: this could still use some work

	now := time.Now().UTC()

	cases := []struct {
		input interface{}
		time  time.Time
		valid bool
		set   bool
	}{
		{input: nil, time: time.Time{}, valid: false, set: true},
		{input: time.Time{}, time: time.Time{}, valid: false, set: true},
		{input: now, time: now, valid: true, set: true},
		{input: "1988-12-16", time: time.Date(1988, time.December, 16, 0, 0, 0, 0, time.UTC), valid: true, set: true},
		{input: "1988-12-16 12:34:56", time: time.Date(1988, time.December, 16, 12, 34, 56, 0, time.UTC), valid: true, set: true},
		{input: "1988-12-16 12:34:56.1", time: time.Date(1988, time.December, 16, 12, 34, 56, 100_000_000, time.UTC), valid: true, set: true},
		{input: "1988-12-16 12:34:56.01", time: time.Date(1988, time.December, 16, 12, 34, 56, 10_000_000, time.UTC), valid: true, set: true},
		{input: "1988-12-16 12:34:56.001", time: time.Date(1988, time.December, 16, 12, 34, 56, 1_000_000, time.UTC), valid: true, set: true},
		{input: "1988-12-16 12:34:56.0001", time: time.Date(1988, time.December, 16, 12, 34, 56, 100_000, time.UTC), valid: true, set: true},
		{input: "1988-12-16 12:34:56.00001", time: time.Date(1988, time.December, 16, 12, 34, 56, 10_000, time.UTC), valid: true, set: true},
		{input: "1988-12-16 12:34:56.000001", time: time.Date(1988, time.December, 16, 12, 34, 56, 1_000, time.UTC), valid: true, set: true},
		{input: false, time: time.Time{}, valid: false, set: true},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v", tc.input)
		t.Run(name, func(t *testing.T) {
			key := fmt.Sprintf("null.Time.Scan(%v)", tc.input)
			var n null.Time
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if n.Time() != tc.time {
				t.Errorf("%s.Time(): got %v, want %v", key, n.Time(), tc.time)
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
				var n null.Time
				err := n.Scan(b)
				if err != nil {
					t.Fatal(err)
				}
				if n.Time() != tc.time {
					t.Errorf("%s.Time(): got %v, want %v", key, n.Time(), tc.time)
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

func TestTime_Value(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	cases := []struct {
		n   null.Time
		val interface{}
	}{
		{n: null.Time{}, val: nil},
		{n: null.NewTime(time.Time{}, false), val: nil},
		{n: null.NewTime(now, false), val: nil},
		{n: null.NewTime(time.Time{}, true), val: nil},
		{n: null.NewTime(now, true), val: now},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%+v", tc.n)
		t.Run(name, func(t *testing.T) {
			v, err := tc.n.Value()
			if err != nil {
				t.Fatal(err)
			}
			if v != tc.val {
				t.Errorf("null.Time.Val(): got %v, want %v", v, tc.val)
			}
		})
	}
}
