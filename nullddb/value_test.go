package nullddb

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bjaus/null"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// --- Constructor Tests ---

type ConstructorSuite struct {
	suite.Suite
}

func TestConstructorSuite(t *testing.T) {
	suite.Run(t, new(ConstructorSuite))
}

func (s *ConstructorSuite) TestNew() {
	v := New("hello")
	s.True(v.IsValid())
	s.Equal("hello", v.Get())
}

func (s *ConstructorSuite) TestNewNull() {
	v := NewNull[string]()
	s.True(v.IsNull())
}

func (s *ConstructorSuite) TestNewPtr_NonNil() {
	str := "hello"
	v := NewPtr(&str)
	s.True(v.IsValid())
	s.Equal("hello", v.Get())
}

func (s *ConstructorSuite) TestNewPtr_Nil() {
	vn := NewPtr[string](nil)
	s.True(vn.IsNull())
}

func (s *ConstructorSuite) TestFrom() {
	nv := null.New("hello")
	v := From(nv)
	s.True(v.IsValid())
	s.Equal("hello", v.Get())
}

// --- Marshal Tests ---

type MarshalSuite struct {
	suite.Suite
}

func TestMarshalSuite(t *testing.T) {
	suite.Run(t, new(MarshalSuite))
}

func (s *MarshalSuite) TestMarshal_String() {
	v := New("hello")
	av, err := v.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	str, ok := av.(*types.AttributeValueMemberS)
	s.Require().True(ok)
	s.Equal("hello", str.Value)
}

func (s *MarshalSuite) TestMarshal_Null() {
	v := NewNull[string]()
	av, err := v.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	_, ok := av.(*types.AttributeValueMemberNULL)
	s.True(ok)
}

func (s *MarshalSuite) TestMarshal_Unset() {
	var v Value[string]
	av, err := v.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	_, ok := av.(*types.AttributeValueMemberNULL)
	s.True(ok)
}

func (s *MarshalSuite) TestMarshal_IntTypes() {
	tests := map[string]struct {
		val any
	}{
		"int":    {New(42)},
		"int64":  {New(int64(42))},
		"int32":  {New(int32(42))},
		"int16":  {New(int16(42))},
		"int8":   {New(int8(42))},
		"uint":   {New(uint(42))},
		"uint64": {New(uint64(42))},
		"uint32": {New(uint32(42))},
		"uint16": {New(uint16(42))},
		"uint8":  {New(uint8(42))},
	}

	for name, tt := range tests {
		s.Run(name, func() {
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

			s.Require().NoError(err)
			n, ok := av.(*types.AttributeValueMemberN)
			s.Require().True(ok)
			s.Equal("42", n.Value)
		})
	}
}

func (s *MarshalSuite) TestMarshal_Float() {
	v64 := New(3.14)
	av, err := v64.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	_, ok := av.(*types.AttributeValueMemberN)
	s.True(ok)

	v32 := New(float32(3.14))
	av, err = v32.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	_, ok = av.(*types.AttributeValueMemberN)
	s.True(ok)
}

func (s *MarshalSuite) TestMarshal_Bool() {
	v := New(true)
	av, err := v.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	b, ok := av.(*types.AttributeValueMemberBOOL)
	s.Require().True(ok)
	s.True(b.Value)
}

func (s *MarshalSuite) TestMarshal_Bytes() {
	v := New([]byte("hello"))
	av, err := v.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	b, ok := av.(*types.AttributeValueMemberB)
	s.Require().True(ok)
	s.Equal("hello", string(b.Value))
}

func (s *MarshalSuite) TestMarshal_Time() {
	now := time.Now()
	v := New(now)
	av, err := v.MarshalDynamoDBAttributeValue()
	s.Require().NoError(err)
	str, ok := av.(*types.AttributeValueMemberS)
	s.Require().True(ok)
	parsed, err := time.Parse(time.RFC3339Nano, str.Value)
	s.Require().NoError(err)
	s.True(parsed.Equal(now))
}

func (s *MarshalSuite) TestMarshal_UnsupportedType() {
	type Custom struct{ X int }
	v := New(Custom{X: 42})
	_, err := v.MarshalDynamoDBAttributeValue()
	s.Error(err)
}

// --- Unmarshal Tests ---

type UnmarshalSuite struct {
	suite.Suite
}

func TestUnmarshalSuite(t *testing.T) {
	suite.Run(t, new(UnmarshalSuite))
}

func (s *UnmarshalSuite) TestUnmarshal_String() {
	av := &types.AttributeValueMemberS{Value: "hello"}
	var v Value[string]
	err := v.UnmarshalDynamoDBAttributeValue(av)
	s.Require().NoError(err)
	s.True(v.IsValid())
	s.Equal("hello", v.Get())
}

func (s *UnmarshalSuite) TestUnmarshal_String_FromN() {
	avn := &types.AttributeValueMemberN{Value: "42"}
	var v Value[string]
	err := v.UnmarshalDynamoDBAttributeValue(avn)
	s.Require().NoError(err)
	s.Equal("42", v.Get())
}

func (s *UnmarshalSuite) TestUnmarshal_String_Error() {
	avb := &types.AttributeValueMemberBOOL{Value: true}
	var v Value[string]
	err := v.UnmarshalDynamoDBAttributeValue(avb)
	s.Error(err)
}

func (s *UnmarshalSuite) TestUnmarshal_Null() {
	av := &types.AttributeValueMemberNULL{Value: true}
	var v Value[string]
	err := v.UnmarshalDynamoDBAttributeValue(av)
	s.Require().NoError(err)
	s.True(v.IsNull())
}

func (s *UnmarshalSuite) TestUnmarshal_IntTypes() {
	av := &types.AttributeValueMemberN{Value: "42"}

	var vi Value[int]
	s.Require().NoError(vi.UnmarshalDynamoDBAttributeValue(av))
	s.Equal(42, vi.Get())

	var vi64 Value[int64]
	s.Require().NoError(vi64.UnmarshalDynamoDBAttributeValue(av))

	var vi32 Value[int32]
	s.Require().NoError(vi32.UnmarshalDynamoDBAttributeValue(av))

	var vi16 Value[int16]
	s.Require().NoError(vi16.UnmarshalDynamoDBAttributeValue(av))

	var vi8 Value[int8]
	s.Require().NoError(vi8.UnmarshalDynamoDBAttributeValue(av))
}

func (s *UnmarshalSuite) TestUnmarshal_Int_Errors() {
	avs := &types.AttributeValueMemberS{Value: "not a number"}
	var v1 Value[int64]
	s.Error(v1.UnmarshalDynamoDBAttributeValue(avs))

	avbad := &types.AttributeValueMemberN{Value: "not a number"}
	var v2 Value[int64]
	s.Error(v2.UnmarshalDynamoDBAttributeValue(avbad))
}

func (s *UnmarshalSuite) TestUnmarshal_UintTypes() {
	av := &types.AttributeValueMemberN{Value: "42"}

	var vu Value[uint]
	s.Require().NoError(vu.UnmarshalDynamoDBAttributeValue(av))
	s.Equal(uint(42), vu.Get())

	var vu64 Value[uint64]
	s.Require().NoError(vu64.UnmarshalDynamoDBAttributeValue(av))

	var vu32 Value[uint32]
	s.Require().NoError(vu32.UnmarshalDynamoDBAttributeValue(av))

	var vu16 Value[uint16]
	s.Require().NoError(vu16.UnmarshalDynamoDBAttributeValue(av))

	var vu8 Value[uint8]
	s.Require().NoError(vu8.UnmarshalDynamoDBAttributeValue(av))
}

func (s *UnmarshalSuite) TestUnmarshal_Uint_Errors() {
	avs := &types.AttributeValueMemberS{Value: "not a number"}
	var v1 Value[uint64]
	s.Error(v1.UnmarshalDynamoDBAttributeValue(avs))

	avbad := &types.AttributeValueMemberN{Value: "not a number"}
	var v2 Value[uint64]
	s.Error(v2.UnmarshalDynamoDBAttributeValue(avbad))
}

func (s *UnmarshalSuite) TestUnmarshal_Float() {
	av := &types.AttributeValueMemberN{Value: "3.14"}

	var vf64 Value[float64]
	s.Require().NoError(vf64.UnmarshalDynamoDBAttributeValue(av))
	s.Equal(3.14, vf64.Get())

	var vf32 Value[float32]
	s.Require().NoError(vf32.UnmarshalDynamoDBAttributeValue(av))
}

func (s *UnmarshalSuite) TestUnmarshal_Float_Errors() {
	avs := &types.AttributeValueMemberS{Value: "not a number"}
	var v1 Value[float64]
	s.Error(v1.UnmarshalDynamoDBAttributeValue(avs))

	avbad := &types.AttributeValueMemberN{Value: "not a number"}
	var v2 Value[float64]
	s.Error(v2.UnmarshalDynamoDBAttributeValue(avbad))
}

func (s *UnmarshalSuite) TestUnmarshal_Bool() {
	av := &types.AttributeValueMemberBOOL{Value: true}
	var v Value[bool]
	s.Require().NoError(v.UnmarshalDynamoDBAttributeValue(av))
	s.True(v.Get())
}

func (s *UnmarshalSuite) TestUnmarshal_Bool_Error() {
	avs := &types.AttributeValueMemberS{Value: "not a bool"}
	var v Value[bool]
	s.Error(v.UnmarshalDynamoDBAttributeValue(avs))
}

func (s *UnmarshalSuite) TestUnmarshal_Bytes() {
	av := &types.AttributeValueMemberB{Value: []byte("hello")}
	var v Value[[]byte]
	s.Require().NoError(v.UnmarshalDynamoDBAttributeValue(av))
	s.Equal("hello", string(v.Get()))
}

func (s *UnmarshalSuite) TestUnmarshal_Bytes_FromS() {
	avs := &types.AttributeValueMemberS{Value: "world"}
	var v Value[[]byte]
	s.Require().NoError(v.UnmarshalDynamoDBAttributeValue(avs))
	s.Equal("world", string(v.Get()))
}

func (s *UnmarshalSuite) TestUnmarshal_Bytes_Error() {
	avn := &types.AttributeValueMemberN{Value: "123"}
	var v Value[[]byte]
	s.Error(v.UnmarshalDynamoDBAttributeValue(avn))
}

func (s *UnmarshalSuite) TestUnmarshal_Time() {
	timeStr := "2024-01-15T10:30:00.123456789Z"
	av := &types.AttributeValueMemberS{Value: timeStr}
	var v Value[time.Time]
	s.Require().NoError(v.UnmarshalDynamoDBAttributeValue(av))
	s.False(v.Get().IsZero())
}

func (s *UnmarshalSuite) TestUnmarshal_Time_Errors() {
	avn := &types.AttributeValueMemberN{Value: "123"}
	var v1 Value[time.Time]
	s.Error(v1.UnmarshalDynamoDBAttributeValue(avn))

	avbad := &types.AttributeValueMemberS{Value: "not a time"}
	var v2 Value[time.Time]
	s.Error(v2.UnmarshalDynamoDBAttributeValue(avbad))
}

func (s *UnmarshalSuite) TestUnmarshal_UnsupportedType() {
	type Custom struct{ X int }
	av := &types.AttributeValueMemberS{Value: "hello"}
	var v Value[Custom]
	s.Error(v.UnmarshalDynamoDBAttributeValue(av))
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
	require.NoError(t, err)

	var decoded Item
	require.NoError(t, attributevalue.UnmarshalMap(av, &decoded))

	require.Equal(t, original.ID, decoded.ID)
	require.True(t, decoded.Name.IsValid())
	require.Equal(t, "Alice", decoded.Name.Get())
	require.True(t, decoded.Age.IsNull())
	require.True(t, decoded.Verified.IsValid())
	require.True(t, decoded.Verified.Get())
}

func TestIntegration_Time(t *testing.T) {
	now := time.Now().Truncate(time.Nanosecond).UTC()
	v := New(now)

	av, err := v.MarshalDynamoDBAttributeValue()
	require.NoError(t, err)

	var decoded Value[time.Time]
	require.NoError(t, decoded.UnmarshalDynamoDBAttributeValue(av))

	require.True(t, decoded.IsValid())
	require.True(t, decoded.Get().Equal(now))
}
