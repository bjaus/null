package null

import (
	"encoding/json"
	"testing"
	"time"
)

// --- State Tests ---

func TestState_String(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{Unset, "unset"},
		{Null, "null"},
		{Valid, "valid"},
		{State(99), "State(99)"}, // unknown state
	}
	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

// --- Constructor Tests ---

func TestNew(t *testing.T) {
	v := New("hello")
	if !v.IsValid() {
		t.Error("New should create valid value")
	}
	if v.Get() != "hello" {
		t.Errorf("expected hello, got %s", v.Get())
	}
	if v.State() != Valid {
		t.Errorf("expected Valid state, got %s", v.State())
	}
}

func TestNewNull(t *testing.T) {
	n := NewNull[string]()
	if !n.IsNull() {
		t.Error("NewNull should create null value")
	}
	if !n.IsSet() {
		t.Error("NewNull should be set")
	}
	if n.IsValid() {
		t.Error("NewNull should not be valid")
	}
	if n.State() != Null {
		t.Errorf("expected Null state, got %s", n.State())
	}
}

func TestNewPtr(t *testing.T) {
	// Non-nil pointer
	s := "world"
	p := NewPtr(&s)
	if !p.IsValid() {
		t.Error("NewPtr with non-nil should be valid")
	}
	if p.Get() != "world" {
		t.Errorf("expected world, got %s", p.Get())
	}

	// Nil pointer
	pn := NewPtr[string](nil)
	if !pn.IsNull() {
		t.Error("NewPtr with nil should be null")
	}
}

func TestZeroValue(t *testing.T) {
	var u Value[string]
	if u.IsSet() {
		t.Error("zero value should be unset")
	}
	if u.IsNull() {
		t.Error("zero value should not be null")
	}
	if u.IsValid() {
		t.Error("zero value should not be valid")
	}
	if u.State() != Unset {
		t.Errorf("expected Unset state, got %s", u.State())
	}
}

// --- Accessor Tests ---

func TestGet(t *testing.T) {
	// Valid
	v := New(42)
	if v.Get() != 42 {
		t.Errorf("expected 42, got %d", v.Get())
	}

	// Null returns zero
	n := NewNull[int]()
	if n.Get() != 0 {
		t.Errorf("expected 0, got %d", n.Get())
	}

	// Unset returns zero
	var u Value[int]
	if u.Get() != 0 {
		t.Errorf("expected 0, got %d", u.Get())
	}
}

func TestGetOr(t *testing.T) {
	// Valid returns value
	v := New(42)
	if v.GetOr(99) != 42 {
		t.Errorf("expected 42, got %d", v.GetOr(99))
	}

	// Null returns default
	n := NewNull[int]()
	if n.GetOr(99) != 99 {
		t.Errorf("expected 99, got %d", n.GetOr(99))
	}

	// Unset returns default
	var u Value[int]
	if u.GetOr(99) != 99 {
		t.Errorf("expected 99, got %d", u.GetOr(99))
	}
}

func TestPtr(t *testing.T) {
	// Valid returns pointer
	v := New(42)
	p := v.Ptr()
	if p == nil {
		t.Fatal("Ptr should not be nil for valid value")
	}
	if *p != 42 {
		t.Errorf("expected *Ptr = 42, got %d", *p)
	}

	// Null returns nil
	n := NewNull[int]()
	if n.Ptr() != nil {
		t.Error("Ptr should be nil for null value")
	}

	// Unset returns nil
	var u Value[int]
	if u.Ptr() != nil {
		t.Error("Ptr should be nil for unset value")
	}
}

// --- JSON Tests ---

func TestJSON_ThreeStates(t *testing.T) {
	type Request struct {
		Name  Value[string] `json:"name"`
		Email Value[string] `json:"email"`
		Age   Value[int]    `json:"age"`
	}

	// All fields present with values
	input1 := `{"name": "Alice", "email": "alice@example.com", "age": 30}`
	var r1 Request
	if err := json.Unmarshal([]byte(input1), &r1); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !r1.Name.IsValid() || r1.Name.Get() != "Alice" {
		t.Errorf("expected valid name Alice, got %v", r1.Name)
	}
	if !r1.Email.IsValid() || r1.Email.Get() != "alice@example.com" {
		t.Errorf("expected valid email, got %v", r1.Email)
	}
	if !r1.Age.IsValid() || r1.Age.Get() != 30 {
		t.Errorf("expected valid age 30, got %v", r1.Age)
	}

	// Some fields explicitly null
	input2 := `{"name": "Bob", "email": null}`
	var r2 Request
	if err := json.Unmarshal([]byte(input2), &r2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !r2.Name.IsValid() {
		t.Error("name should be valid")
	}
	if !r2.Email.IsNull() {
		t.Error("email should be null")
	}
	if r2.Age.IsSet() {
		t.Error("age should be unset (not present in JSON)")
	}

	// Empty string is different from null
	input3 := `{"name": ""}`
	var r3 Request
	if err := json.Unmarshal([]byte(input3), &r3); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !r3.Name.IsValid() {
		t.Error("name should be valid (empty string is a value)")
	}
	if r3.Name.Get() != "" {
		t.Error("name should be empty string")
	}
	if r3.Name.IsNull() {
		t.Error("empty string is not null")
	}
}

func TestJSON_Marshal(t *testing.T) {
	type Response struct {
		Name  Value[string] `json:"name"`
		Email Value[string] `json:"email"`
		Age   Value[int]    `json:"age"`
	}

	r := Response{
		Name:  New("Alice"),
		Email: NewNull[string](),
		// Age is unset (zero value)
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := `{"name":"Alice","email":null,"age":null}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestJSON_UnmarshalError(t *testing.T) {
	var v Value[int]
	err := json.Unmarshal([]byte(`"not a number"`), &v)
	if err == nil {
		t.Error("expected unmarshal error for invalid type")
	}
}

func TestZeroValueIsValid(t *testing.T) {
	// Zero value of T can be valid
	zero := New(0)
	if !zero.IsValid() {
		t.Error("zero int should be valid")
	}
	if zero.Get() != 0 {
		t.Error("zero int should return 0")
	}

	empty := New("")
	if !empty.IsValid() {
		t.Error("empty string should be valid")
	}

	falseVal := New(false)
	if !falseVal.IsValid() {
		t.Error("false should be valid")
	}
}

// --- SQL Valuer Tests ---

func TestValue_SQLValuer(t *testing.T) {
	// Valid string
	vs := New("hello")
	v, err := vs.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "hello" {
		t.Errorf("expected hello, got %v", v)
	}

	// Null returns nil
	ns := NewNull[string]()
	v, err = ns.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}

	// Unset returns nil
	var us Value[string]
	v, err = us.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}

func TestValue_SQLValuer_Types(t *testing.T) {
	// int types
	vi := New(42)
	v, _ := vi.Value()
	if v != int64(42) {
		t.Errorf("int: expected 42, got %v", v)
	}

	vi32 := New(int32(42))
	v, _ = vi32.Value()
	if v != int64(42) {
		t.Errorf("int32: expected 42, got %v", v)
	}

	vi16 := New(int16(42))
	v, _ = vi16.Value()
	if v != int64(42) {
		t.Errorf("int16: expected 42, got %v", v)
	}

	vi8 := New(int8(42))
	v, _ = vi8.Value()
	if v != int64(42) {
		t.Errorf("int8: expected 42, got %v", v)
	}

	vi64 := New(int64(42))
	v, _ = vi64.Value()
	if v != int64(42) {
		t.Errorf("int64: expected 42, got %v", v)
	}

	// uint types
	vu := New(uint(42))
	v, _ = vu.Value()
	if v != int64(42) {
		t.Errorf("uint: expected 42, got %v", v)
	}

	vu64 := New(uint64(42))
	v, _ = vu64.Value()
	if v != int64(42) {
		t.Errorf("uint64: expected 42, got %v", v)
	}

	vu32 := New(uint32(42))
	v, _ = vu32.Value()
	if v != int64(42) {
		t.Errorf("uint32: expected 42, got %v", v)
	}

	vu16 := New(uint16(42))
	v, _ = vu16.Value()
	if v != int64(42) {
		t.Errorf("uint16: expected 42, got %v", v)
	}

	vu8 := New(uint8(42))
	v, _ = vu8.Value()
	if v != int64(42) {
		t.Errorf("uint8: expected 42, got %v", v)
	}

	// float types
	vf64 := New(3.14)
	v, _ = vf64.Value()
	if v != 3.14 {
		t.Errorf("float64: expected 3.14, got %v", v)
	}

	vf32 := New(float32(3.14))
	v, _ = vf32.Value()
	if v != float64(float32(3.14)) {
		t.Errorf("float32: expected %v, got %v", float64(float32(3.14)), v)
	}

	// bool
	vb := New(true)
	v, _ = vb.Value()
	if v != true {
		t.Errorf("bool: expected true, got %v", v)
	}

	// []byte
	vby := New([]byte("hello"))
	v, _ = vby.Value()
	if string(v.([]byte)) != "hello" {
		t.Errorf("[]byte: expected hello, got %v", v)
	}

	// time.Time
	now := time.Now()
	vt := New(now)
	v, _ = vt.Value()
	if v != now {
		t.Errorf("time.Time: expected %v, got %v", now, v)
	}

	// Custom type (falls through to default)
	type Custom struct{ X int }
	vc := New(Custom{X: 42})
	v, _ = vc.Value()
	if c, ok := v.(Custom); !ok || c.X != 42 {
		t.Errorf("Custom: expected {X:42}, got %v", v)
	}
}

// --- SQL Scanner Tests ---

func TestValue_SQLScanner_Nil(t *testing.T) {
	var v Value[string]
	if err := v.Scan(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.IsNull() {
		t.Error("expected null after scanning nil")
	}
}

func TestValue_SQLScanner_String(t *testing.T) {
	var v Value[string]
	if err := v.Scan("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.IsValid() || v.Get() != "hello" {
		t.Errorf("expected valid 'hello', got %v", v)
	}

	// []byte to string
	var v2 Value[string]
	if err := v2.Scan([]byte("world")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v2.IsValid() || v2.Get() != "world" {
		t.Errorf("expected valid 'world', got %v", v2)
	}

	// Error case
	var v3 Value[string]
	if err := v3.Scan(123); err == nil {
		t.Error("expected error scanning int into string")
	}
}

func TestValue_SQLScanner_Int(t *testing.T) {
	// int
	var vi Value[int]
	if err := vi.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi.Get() != 42 {
		t.Errorf("expected 42, got %d", vi.Get())
	}

	// int64
	var vi64 Value[int64]
	if err := vi64.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi64.Get() != 42 {
		t.Errorf("expected 42, got %d", vi64.Get())
	}

	// int32
	var vi32 Value[int32]
	if err := vi32.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi32.Get() != 42 {
		t.Errorf("expected 42, got %d", vi32.Get())
	}

	// int16
	var vi16 Value[int16]
	if err := vi16.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi16.Get() != 42 {
		t.Errorf("expected 42, got %d", vi16.Get())
	}

	// int8
	var vi8 Value[int8]
	if err := vi8.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi8.Get() != 42 {
		t.Errorf("expected 42, got %d", vi8.Get())
	}

	// from int (not int64)
	var vi2 Value[int64]
	if err := vi2.Scan(int(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi2.Get() != 42 {
		t.Errorf("expected 42, got %d", vi2.Get())
	}

	// from int32
	var vi3 Value[int64]
	if err := vi3.Scan(int32(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi3.Get() != 42 {
		t.Errorf("expected 42, got %d", vi3.Get())
	}

	// from float64
	var vi4 Value[int64]
	if err := vi4.Scan(float64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vi4.Get() != 42 {
		t.Errorf("expected 42, got %d", vi4.Get())
	}

	// Error case
	var vi5 Value[int64]
	if err := vi5.Scan("not a number"); err == nil {
		t.Error("expected error scanning string into int64")
	}
}

func TestValue_SQLScanner_Uint(t *testing.T) {
	// uint
	var vu Value[uint]
	if err := vu.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu.Get() != 42 {
		t.Errorf("expected 42, got %d", vu.Get())
	}

	// uint64
	var vu64 Value[uint64]
	if err := vu64.Scan(uint64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu64.Get() != 42 {
		t.Errorf("expected 42, got %d", vu64.Get())
	}

	// uint32
	var vu32 Value[uint32]
	if err := vu32.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu32.Get() != 42 {
		t.Errorf("expected 42, got %d", vu32.Get())
	}

	// uint16
	var vu16 Value[uint16]
	if err := vu16.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu16.Get() != 42 {
		t.Errorf("expected 42, got %d", vu16.Get())
	}

	// uint8
	var vu8 Value[uint8]
	if err := vu8.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu8.Get() != 42 {
		t.Errorf("expected 42, got %d", vu8.Get())
	}

	// from uint
	var vu2 Value[uint64]
	if err := vu2.Scan(uint(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu2.Get() != 42 {
		t.Errorf("expected 42, got %d", vu2.Get())
	}

	// from uint32
	var vu3 Value[uint64]
	if err := vu3.Scan(uint32(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu3.Get() != 42 {
		t.Errorf("expected 42, got %d", vu3.Get())
	}

	// from float64
	var vu4 Value[uint64]
	if err := vu4.Scan(float64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vu4.Get() != 42 {
		t.Errorf("expected 42, got %d", vu4.Get())
	}

	// Error case
	var vu5 Value[uint64]
	if err := vu5.Scan("not a number"); err == nil {
		t.Error("expected error scanning string into uint64")
	}
}

func TestValue_SQLScanner_Float(t *testing.T) {
	// float64
	var vf64 Value[float64]
	if err := vf64.Scan(float64(3.14)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vf64.Get() != 3.14 {
		t.Errorf("expected 3.14, got %f", vf64.Get())
	}

	// float32
	var vf32 Value[float32]
	if err := vf32.Scan(float64(3.14)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vf32.Get() != float32(3.14) {
		t.Errorf("expected %f, got %f", float32(3.14), vf32.Get())
	}

	// from float32
	var vf2 Value[float64]
	if err := vf2.Scan(float32(3.14)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// from int64
	var vf3 Value[float64]
	if err := vf3.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vf3.Get() != 42 {
		t.Errorf("expected 42, got %f", vf3.Get())
	}

	// Error case
	var vf4 Value[float64]
	if err := vf4.Scan("not a number"); err == nil {
		t.Error("expected error scanning string into float64")
	}
}

func TestValue_SQLScanner_Bool(t *testing.T) {
	var vb Value[bool]
	if err := vb.Scan(true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !vb.Get() {
		t.Error("expected true")
	}

	// from int64
	var vb2 Value[bool]
	if err := vb2.Scan(int64(1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !vb2.Get() {
		t.Error("expected true from 1")
	}

	var vb3 Value[bool]
	if err := vb3.Scan(int64(0)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vb3.Get() {
		t.Error("expected false from 0")
	}

	// Error case
	var vb4 Value[bool]
	if err := vb4.Scan("not a bool"); err == nil {
		t.Error("expected error scanning string into bool")
	}
}

func TestValue_SQLScanner_Time(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	var vt Value[time.Time]
	if err := vt.Scan(now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !vt.Get().Equal(now) {
		t.Errorf("expected %v, got %v", now, vt.Get())
	}

	// from string
	timeStr := "2024-01-15T10:30:00Z"
	var vt2 Value[time.Time]
	if err := vt2.Scan(timeStr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// from []byte
	var vt3 Value[time.Time]
	if err := vt3.Scan([]byte(timeStr)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Error: invalid time string
	var vt4 Value[time.Time]
	if err := vt4.Scan("not a time"); err == nil {
		t.Error("expected error scanning invalid time string")
	}

	// Error: invalid time bytes
	var vt5 Value[time.Time]
	if err := vt5.Scan([]byte("not a time")); err == nil {
		t.Error("expected error scanning invalid time bytes")
	}

	// Error: wrong type
	var vt6 Value[time.Time]
	if err := vt6.Scan(123); err == nil {
		t.Error("expected error scanning int into time.Time")
	}
}

func TestValue_SQLScanner_Bytes(t *testing.T) {
	var vb Value[[]byte]
	if err := vb.Scan([]byte("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(vb.Get()) != "hello" {
		t.Errorf("expected hello, got %s", vb.Get())
	}

	// from string
	var vb2 Value[[]byte]
	if err := vb2.Scan("world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(vb2.Get()) != "world" {
		t.Errorf("expected world, got %s", vb2.Get())
	}

	// Error case
	var vb3 Value[[]byte]
	if err := vb3.Scan(123); err == nil {
		t.Error("expected error scanning int into []byte")
	}
}

func TestValue_SQLScanner_Reflect(t *testing.T) {
	// Custom type that uses reflection fallback
	type MyInt int

	// Assignable
	var v Value[MyInt]
	if err := v.Scan(MyInt(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Get() != 42 {
		t.Errorf("expected 42, got %d", v.Get())
	}

	// Convertible
	var v2 Value[MyInt]
	if err := v2.Scan(int(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v2.Get() != 42 {
		t.Errorf("expected 42, got %d", v2.Get())
	}

	// Not convertible
	var v3 Value[MyInt]
	if err := v3.Scan("not convertible"); err == nil {
		t.Error("expected error for non-convertible type")
	}
}
