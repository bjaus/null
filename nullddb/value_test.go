package nullddb

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bjaus/null"
)

// --- Constructor Tests ---

func TestNew(t *testing.T) {
	v := New("hello")
	if !v.IsValid() || v.Get() != "hello" {
		t.Errorf("expected valid 'hello', got %v", v)
	}
}

func TestNewNull(t *testing.T) {
	v := NewNull[string]()
	if !v.IsNull() {
		t.Error("expected null")
	}
}

func TestNewPtr(t *testing.T) {
	s := "hello"
	v := NewPtr(&s)
	if !v.IsValid() || v.Get() != "hello" {
		t.Errorf("expected valid 'hello', got %v", v)
	}

	vn := NewPtr[string](nil)
	if !vn.IsNull() {
		t.Error("expected null for nil pointer")
	}
}

func TestFrom(t *testing.T) {
	nv := null.New("hello")
	v := From(nv)
	if !v.IsValid() || v.Get() != "hello" {
		t.Errorf("expected valid 'hello', got %v", v)
	}
}

// --- Marshal Tests ---

func TestMarshal_String(t *testing.T) {
	v := New("hello")
	av, err := v.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := av.(*types.AttributeValueMemberS)
	if !ok || s.Value != "hello" {
		t.Errorf("expected S='hello', got %#v", av)
	}
}

func TestMarshal_Null(t *testing.T) {
	v := NewNull[string]()
	av, err := v.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := av.(*types.AttributeValueMemberNULL); !ok {
		t.Errorf("expected NULL, got %#v", av)
	}
}

func TestMarshal_Unset(t *testing.T) {
	var v Value[string]
	av, err := v.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := av.(*types.AttributeValueMemberNULL); !ok {
		t.Errorf("expected NULL for unset, got %#v", av)
	}
}

func TestMarshal_IntTypes(t *testing.T) {
	tests := []struct {
		name string
		val  any
	}{
		{"int", New(42)},
		{"int64", New(int64(42))},
		{"int32", New(int32(42))},
		{"int16", New(int16(42))},
		{"int8", New(int8(42))},
		{"uint", New(uint(42))},
		{"uint64", New(uint64(42))},
		{"uint32", New(uint32(42))},
		{"uint16", New(uint16(42))},
		{"uint8", New(uint8(42))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var av types.AttributeValue
			var err error

			switch v := tt.val.(type) {
			case Value[int]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[int64]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[int32]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[int16]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[int8]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[uint]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[uint64]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[uint32]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[uint16]:
				av, err = v.MarshalDynamoDBAttributeValue()
			case Value[uint8]:
				av, err = v.MarshalDynamoDBAttributeValue()
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			n, ok := av.(*types.AttributeValueMemberN)
			if !ok {
				t.Errorf("expected N, got %#v", av)
			}
			if n.Value != "42" {
				t.Errorf("expected '42', got %q", n.Value)
			}
		})
	}
}

func TestMarshal_Float(t *testing.T) {
	v64 := New(3.14)
	av, err := v64.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := av.(*types.AttributeValueMemberN); !ok {
		t.Errorf("expected N, got %#v", av)
	}

	v32 := New(float32(3.14))
	av, err = v32.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := av.(*types.AttributeValueMemberN); !ok {
		t.Errorf("expected N, got %#v", av)
	}
}

func TestMarshal_Bool(t *testing.T) {
	v := New(true)
	av, err := v.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, ok := av.(*types.AttributeValueMemberBOOL)
	if !ok || !b.Value {
		t.Errorf("expected BOOL=true, got %#v", av)
	}
}

func TestMarshal_Bytes(t *testing.T) {
	v := New([]byte("hello"))
	av, err := v.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, ok := av.(*types.AttributeValueMemberB)
	if !ok || string(b.Value) != "hello" {
		t.Errorf("expected B='hello', got %#v", av)
	}
}

func TestMarshal_Time(t *testing.T) {
	now := time.Now()
	v := New(now)
	av, err := v.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := av.(*types.AttributeValueMemberS)
	if !ok {
		t.Errorf("expected S, got %#v", av)
	}
	parsed, err := time.Parse(time.RFC3339Nano, s.Value)
	if err != nil {
		t.Fatalf("cannot parse time: %v", err)
	}
	if !parsed.Equal(now) {
		t.Errorf("time mismatch: got %v, want %v", parsed, now)
	}
}

func TestMarshal_UnsupportedType(t *testing.T) {
	type Custom struct{ X int }
	v := New(Custom{X: 42})
	_, err := v.MarshalDynamoDBAttributeValue()
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

// --- Unmarshal Tests ---

func TestUnmarshal_String(t *testing.T) {
	av := &types.AttributeValueMemberS{Value: "hello"}
	var v Value[string]
	if err := v.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.IsValid() || v.Get() != "hello" {
		t.Errorf("expected valid 'hello', got %v", v)
	}

	// N to string
	avn := &types.AttributeValueMemberN{Value: "42"}
	var v2 Value[string]
	if err := v2.UnmarshalDynamoDBAttributeValue(avn); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v2.Get() != "42" {
		t.Errorf("expected '42', got %q", v2.Get())
	}

	// Error case
	avb := &types.AttributeValueMemberBOOL{Value: true}
	var v3 Value[string]
	if err := v3.UnmarshalDynamoDBAttributeValue(avb); err == nil {
		t.Error("expected error unmarshaling BOOL into string")
	}
}

func TestUnmarshal_Null(t *testing.T) {
	av := &types.AttributeValueMemberNULL{Value: true}
	var v Value[string]
	if err := v.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.IsNull() {
		t.Error("expected null")
	}
}

func TestUnmarshal_IntTypes(t *testing.T) {
	av := &types.AttributeValueMemberN{Value: "42"}

	var vi Value[int]
	if err := vi.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("int: %v", err)
	}
	if vi.Get() != 42 {
		t.Errorf("int: expected 42, got %d", vi.Get())
	}

	var vi64 Value[int64]
	if err := vi64.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("int64: %v", err)
	}

	var vi32 Value[int32]
	if err := vi32.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("int32: %v", err)
	}

	var vi16 Value[int16]
	if err := vi16.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("int16: %v", err)
	}

	var vi8 Value[int8]
	if err := vi8.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("int8: %v", err)
	}

	// Error cases
	avs := &types.AttributeValueMemberS{Value: "not a number"}
	var vie Value[int64]
	if err := vie.UnmarshalDynamoDBAttributeValue(avs); err == nil {
		t.Error("expected error unmarshaling S into int64")
	}

	avbad := &types.AttributeValueMemberN{Value: "not a number"}
	var vie2 Value[int64]
	if err := vie2.UnmarshalDynamoDBAttributeValue(avbad); err == nil {
		t.Error("expected error parsing invalid number")
	}
}

func TestUnmarshal_UintTypes(t *testing.T) {
	av := &types.AttributeValueMemberN{Value: "42"}

	var vu Value[uint]
	if err := vu.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("uint: %v", err)
	}
	if vu.Get() != 42 {
		t.Errorf("uint: expected 42, got %d", vu.Get())
	}

	var vu64 Value[uint64]
	if err := vu64.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("uint64: %v", err)
	}

	var vu32 Value[uint32]
	if err := vu32.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("uint32: %v", err)
	}

	var vu16 Value[uint16]
	if err := vu16.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("uint16: %v", err)
	}

	var vu8 Value[uint8]
	if err := vu8.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("uint8: %v", err)
	}

	// Error cases
	avs := &types.AttributeValueMemberS{Value: "not a number"}
	var vue Value[uint64]
	if err := vue.UnmarshalDynamoDBAttributeValue(avs); err == nil {
		t.Error("expected error unmarshaling S into uint64")
	}

	avbad := &types.AttributeValueMemberN{Value: "not a number"}
	var vue2 Value[uint64]
	if err := vue2.UnmarshalDynamoDBAttributeValue(avbad); err == nil {
		t.Error("expected error parsing invalid number")
	}
}

func TestUnmarshal_Float(t *testing.T) {
	av := &types.AttributeValueMemberN{Value: "3.14"}

	var vf64 Value[float64]
	if err := vf64.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("float64: %v", err)
	}
	if vf64.Get() != 3.14 {
		t.Errorf("float64: expected 3.14, got %f", vf64.Get())
	}

	var vf32 Value[float32]
	if err := vf32.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("float32: %v", err)
	}

	// Error cases
	avs := &types.AttributeValueMemberS{Value: "not a number"}
	var vfe Value[float64]
	if err := vfe.UnmarshalDynamoDBAttributeValue(avs); err == nil {
		t.Error("expected error unmarshaling S into float64")
	}

	avbad := &types.AttributeValueMemberN{Value: "not a number"}
	var vfe2 Value[float64]
	if err := vfe2.UnmarshalDynamoDBAttributeValue(avbad); err == nil {
		t.Error("expected error parsing invalid number")
	}
}

func TestUnmarshal_Bool(t *testing.T) {
	av := &types.AttributeValueMemberBOOL{Value: true}
	var v Value[bool]
	if err := v.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.Get() {
		t.Error("expected true")
	}

	// Error case
	avs := &types.AttributeValueMemberS{Value: "not a bool"}
	var ve Value[bool]
	if err := ve.UnmarshalDynamoDBAttributeValue(avs); err == nil {
		t.Error("expected error unmarshaling S into bool")
	}
}

func TestUnmarshal_Bytes(t *testing.T) {
	av := &types.AttributeValueMemberB{Value: []byte("hello")}
	var v Value[[]byte]
	if err := v.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(v.Get()) != "hello" {
		t.Errorf("expected 'hello', got %q", v.Get())
	}

	// S to []byte
	avs := &types.AttributeValueMemberS{Value: "world"}
	var v2 Value[[]byte]
	if err := v2.UnmarshalDynamoDBAttributeValue(avs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(v2.Get()) != "world" {
		t.Errorf("expected 'world', got %q", v2.Get())
	}

	// Error case
	avn := &types.AttributeValueMemberN{Value: "123"}
	var ve Value[[]byte]
	if err := ve.UnmarshalDynamoDBAttributeValue(avn); err == nil {
		t.Error("expected error unmarshaling N into []byte")
	}
}

func TestUnmarshal_Time(t *testing.T) {
	timeStr := "2024-01-15T10:30:00.123456789Z"
	av := &types.AttributeValueMemberS{Value: timeStr}
	var v Value[time.Time]
	if err := v.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Get().IsZero() {
		t.Error("expected non-zero time")
	}

	// Error: not S type
	avn := &types.AttributeValueMemberN{Value: "123"}
	var ve Value[time.Time]
	if err := ve.UnmarshalDynamoDBAttributeValue(avn); err == nil {
		t.Error("expected error unmarshaling N into time.Time")
	}

	// Error: invalid time
	avbad := &types.AttributeValueMemberS{Value: "not a time"}
	var ve2 Value[time.Time]
	if err := ve2.UnmarshalDynamoDBAttributeValue(avbad); err == nil {
		t.Error("expected error parsing invalid time")
	}
}

func TestUnmarshal_UnsupportedType(t *testing.T) {
	type Custom struct{ X int }
	av := &types.AttributeValueMemberS{Value: "hello"}
	var v Value[Custom]
	if err := v.UnmarshalDynamoDBAttributeValue(av); err == nil {
		t.Error("expected error for unsupported type")
	}
}

// --- Integration Tests ---

func TestIntegration_MarshalUnmarshal(t *testing.T) {
	type Item struct {
		ID       string        `dynamodbav:"id"`
		Name     Value[string] `dynamodbav:"name"`
		Age      Value[int]    `dynamodbav:"age"`
		Verified Value[bool]   `dynamodbav:"verified"`
	}

	original := Item{
		ID:       "123",
		Name:     New("Alice"),
		Age:      NewNull[int](),
		Verified: New(true),
	}

	av, err := attributevalue.MarshalMap(original)
	if err != nil {
		t.Fatalf("MarshalMap failed: %v", err)
	}

	var decoded Item
	if err := attributevalue.UnmarshalMap(av, &decoded); err != nil {
		t.Fatalf("UnmarshalMap failed: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID mismatch")
	}
	if !decoded.Name.IsValid() || decoded.Name.Get() != "Alice" {
		t.Errorf("Name mismatch")
	}
	if !decoded.Age.IsNull() {
		t.Errorf("Age should be null")
	}
	if !decoded.Verified.IsValid() || !decoded.Verified.Get() {
		t.Errorf("Verified mismatch")
	}
}

func TestIntegration_Time(t *testing.T) {
	now := time.Now().Truncate(time.Nanosecond).UTC()
	v := New(now)

	av, err := v.MarshalDynamoDBAttributeValue()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Value[time.Time]
	if err := decoded.UnmarshalDynamoDBAttributeValue(av); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !decoded.IsValid() {
		t.Error("expected valid")
	}
	if !decoded.Get().Equal(now) {
		t.Errorf("got %v, want %v", decoded.Get(), now)
	}
}
