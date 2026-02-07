package null

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// --- State Tests ---

func TestState_String(t *testing.T) {
	tests := map[string]struct {
		state State
		want  string
	}{
		"unset":   {Unset, "unset"},
		"null":    {Null, "null"},
		"valid":   {Valid, "valid"},
		"unknown": {State(99), "State(99)"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.String())
		})
	}
}

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
	s.Equal(Valid, v.State())
}

func (s *ConstructorSuite) TestNewNull() {
	n := NewNull[string]()
	s.True(n.IsNull())
	s.True(n.IsSet())
	s.False(n.IsValid())
	s.Equal(Null, n.State())
}

func (s *ConstructorSuite) TestNewPtr_NonNil() {
	str := "world"
	p := NewPtr(&str)
	s.True(p.IsValid())
	s.Equal("world", p.Get())
}

func (s *ConstructorSuite) TestNewPtr_Nil() {
	pn := NewPtr[string](nil)
	s.True(pn.IsNull())
}

func (s *ConstructorSuite) TestZeroValue() {
	var u Value[string]
	s.False(u.IsSet())
	s.False(u.IsNull())
	s.False(u.IsValid())
	s.Equal(Unset, u.State())
}

// --- Accessor Tests ---

type AccessorSuite struct {
	suite.Suite
}

func TestAccessorSuite(t *testing.T) {
	suite.Run(t, new(AccessorSuite))
}

func (s *AccessorSuite) TestGet_Valid() {
	v := New(42)
	s.Equal(42, v.Get())
}

func (s *AccessorSuite) TestGet_Null() {
	n := NewNull[int]()
	s.Equal(0, n.Get())
}

func (s *AccessorSuite) TestGet_Unset() {
	var u Value[int]
	s.Equal(0, u.Get())
}

func (s *AccessorSuite) TestGetOr_Valid() {
	v := New(42)
	s.Equal(42, v.GetOr(99))
}

func (s *AccessorSuite) TestGetOr_Null() {
	n := NewNull[int]()
	s.Equal(99, n.GetOr(99))
}

func (s *AccessorSuite) TestGetOr_Unset() {
	var u Value[int]
	s.Equal(99, u.GetOr(99))
}

func (s *AccessorSuite) TestPtr_Valid() {
	v := New(42)
	p := v.Ptr()
	s.Require().NotNil(p)
	s.Equal(42, *p)
}

func (s *AccessorSuite) TestPtr_Null() {
	n := NewNull[int]()
	s.Nil(n.Ptr())
}

func (s *AccessorSuite) TestPtr_Unset() {
	var u Value[int]
	s.Nil(u.Ptr())
}

// --- JSON Tests ---

type JSONSuite struct {
	suite.Suite
}

func TestJSONSuite(t *testing.T) {
	suite.Run(t, new(JSONSuite))
}

func (s *JSONSuite) TestUnmarshal_AllFieldsPresent() {
	type Request struct {
		Name  Value[string] `json:"name"`
		Email Value[string] `json:"email"`
		Age   Value[int]    `json:"age"`
	}

	input := `{"name": "Alice", "email": "alice@example.com", "age": 30}`
	var r Request
	err := json.Unmarshal([]byte(input), &r)
	s.Require().NoError(err)

	s.True(r.Name.IsValid())
	s.Equal("Alice", r.Name.Get())
	s.True(r.Email.IsValid())
	s.Equal("alice@example.com", r.Email.Get())
	s.True(r.Age.IsValid())
	s.Equal(30, r.Age.Get())
}

func (s *JSONSuite) TestUnmarshal_ExplicitNull() {
	type Request struct {
		Name  Value[string] `json:"name"`
		Email Value[string] `json:"email"`
		Age   Value[int]    `json:"age"`
	}

	input := `{"name": "Bob", "email": null}`
	var r Request
	err := json.Unmarshal([]byte(input), &r)
	s.Require().NoError(err)

	s.True(r.Name.IsValid())
	s.True(r.Email.IsNull())
	s.False(r.Age.IsSet())
}

func (s *JSONSuite) TestUnmarshal_EmptyStringIsNotNull() {
	type Request struct {
		Name Value[string] `json:"name"`
	}

	input := `{"name": ""}`
	var r Request
	err := json.Unmarshal([]byte(input), &r)
	s.Require().NoError(err)

	s.True(r.Name.IsValid())
	s.Equal("", r.Name.Get())
	s.False(r.Name.IsNull())
}

func (s *JSONSuite) TestMarshal() {
	type Response struct {
		Name  Value[string] `json:"name"`
		Email Value[string] `json:"email"`
		Age   Value[int]    `json:"age"`
	}

	r := Response{
		Name:  New("Alice"),
		Email: NewNull[string](),
		// Age is unset
	}

	data, err := json.Marshal(r)
	s.Require().NoError(err)
	s.Equal(`{"name":"Alice","email":null,"age":null}`, string(data))
}

func (s *JSONSuite) TestUnmarshal_InvalidType() {
	var v Value[int]
	err := json.Unmarshal([]byte(`"not a number"`), &v)
	s.Error(err)
}

func (s *JSONSuite) TestZeroValueIsValid() {
	zero := New(0)
	s.True(zero.IsValid())
	s.Equal(0, zero.Get())

	empty := New("")
	s.True(empty.IsValid())

	falseVal := New(false)
	s.True(falseVal.IsValid())
}

// --- SQL Valuer Tests ---

type SQLValuerSuite struct {
	suite.Suite
}

func TestSQLValuerSuite(t *testing.T) {
	suite.Run(t, new(SQLValuerSuite))
}

func (s *SQLValuerSuite) TestValue_String() {
	vs := New("hello")
	v, err := vs.Value()
	s.Require().NoError(err)
	s.Equal("hello", v)
}

func (s *SQLValuerSuite) TestValue_Null() {
	ns := NewNull[string]()
	v, err := ns.Value()
	s.Require().NoError(err)
	s.Nil(v)
}

func (s *SQLValuerSuite) TestValue_Unset() {
	var us Value[string]
	v, err := us.Value()
	s.Require().NoError(err)
	s.Nil(v)
}

func (s *SQLValuerSuite) TestValue_IntTypes() {
	tests := map[string]struct {
		val      any
		expected int64
	}{
		"int":    {New(42), 42},
		"int64":  {New(int64(42)), 42},
		"int32":  {New(int32(42)), 42},
		"int16":  {New(int16(42)), 42},
		"int8":   {New(int8(42)), 42},
		"uint":   {New(uint(42)), 42},
		"uint64": {New(uint64(42)), 42},
		"uint32": {New(uint32(42)), 42},
		"uint16": {New(uint16(42)), 42},
		"uint8":  {New(uint8(42)), 42},
	}

	for name, tt := range tests {
		s.Run(name, func() {
			var v any
			var err error
			switch val := tt.val.(type) {
			case Value[int]:
				v, err = val.Value()
			case Value[int64]:
				v, err = val.Value()
			case Value[int32]:
				v, err = val.Value()
			case Value[int16]:
				v, err = val.Value()
			case Value[int8]:
				v, err = val.Value()
			case Value[uint]:
				v, err = val.Value()
			case Value[uint64]:
				v, err = val.Value()
			case Value[uint32]:
				v, err = val.Value()
			case Value[uint16]:
				v, err = val.Value()
			case Value[uint8]:
				v, err = val.Value()
			}
			s.Require().NoError(err)
			s.Equal(tt.expected, v)
		})
	}
}

func (s *SQLValuerSuite) TestValue_Float() {
	vf64 := New(3.14)
	v, err := vf64.Value()
	s.Require().NoError(err)
	s.Equal(3.14, v)

	vf32 := New(float32(3.14))
	v, err = vf32.Value()
	s.Require().NoError(err)
	s.Equal(float64(float32(3.14)), v)
}

func (s *SQLValuerSuite) TestValue_Bool() {
	vb := New(true)
	v, err := vb.Value()
	s.Require().NoError(err)
	s.Equal(true, v)
}

func (s *SQLValuerSuite) TestValue_Bytes() {
	vby := New([]byte("hello"))
	v, err := vby.Value()
	s.Require().NoError(err)
	s.Equal("hello", string(v.([]byte)))
}

func (s *SQLValuerSuite) TestValue_Time() {
	now := time.Now()
	vt := New(now)
	v, err := vt.Value()
	s.Require().NoError(err)
	s.Equal(now, v)
}

func (s *SQLValuerSuite) TestValue_CustomType() {
	type Custom struct{ X int }
	vc := New(Custom{X: 42})
	v, err := vc.Value()
	s.Require().NoError(err)
	c, ok := v.(Custom)
	s.True(ok)
	s.Equal(42, c.X)
}

// --- SQL Scanner Tests ---

type SQLScannerSuite struct {
	suite.Suite
}

func TestSQLScannerSuite(t *testing.T) {
	suite.Run(t, new(SQLScannerSuite))
}

func (s *SQLScannerSuite) TestScan_Nil() {
	var v Value[string]
	err := v.Scan(nil)
	s.Require().NoError(err)
	s.True(v.IsNull())
}

func (s *SQLScannerSuite) TestScan_String() {
	var v Value[string]
	err := v.Scan("hello")
	s.Require().NoError(err)
	s.True(v.IsValid())
	s.Equal("hello", v.Get())
}

func (s *SQLScannerSuite) TestScan_BytesToString() {
	var v Value[string]
	err := v.Scan([]byte("world"))
	s.Require().NoError(err)
	s.True(v.IsValid())
	s.Equal("world", v.Get())
}

func (s *SQLScannerSuite) TestScan_String_Error() {
	var v Value[string]
	err := v.Scan(123)
	s.Error(err)
}

func (s *SQLScannerSuite) TestScan_IntTypes() {
	tests := map[string]struct {
		target any
		input  any
		want   int64
	}{
		"int":        {new(Value[int]), int64(42), 42},
		"int64":      {new(Value[int64]), int64(42), 42},
		"int32":      {new(Value[int32]), int64(42), 42},
		"int16":      {new(Value[int16]), int64(42), 42},
		"int8":       {new(Value[int8]), int64(42), 42},
		"from_int":   {new(Value[int64]), int(42), 42},
		"from_int32": {new(Value[int64]), int32(42), 42},
		"from_float": {new(Value[int64]), float64(42), 42},
	}

	for name, tt := range tests {
		s.Run(name, func() {
			var got int64
			var err error
			switch v := tt.target.(type) {
			case *Value[int]:
				err = v.Scan(tt.input)
				got = int64(v.Get())
			case *Value[int64]:
				err = v.Scan(tt.input)
				got = v.Get()
			case *Value[int32]:
				err = v.Scan(tt.input)
				got = int64(v.Get())
			case *Value[int16]:
				err = v.Scan(tt.input)
				got = int64(v.Get())
			case *Value[int8]:
				err = v.Scan(tt.input)
				got = int64(v.Get())
			}
			s.Require().NoError(err)
			s.Equal(tt.want, got)
		})
	}
}

func (s *SQLScannerSuite) TestScan_Int_Error() {
	var v Value[int64]
	err := v.Scan("not a number")
	s.Error(err)
}

func (s *SQLScannerSuite) TestScan_UintTypes() {
	tests := map[string]struct {
		target any
		input  any
		want   uint64
	}{
		"uint":        {new(Value[uint]), int64(42), 42},
		"uint64":      {new(Value[uint64]), uint64(42), 42},
		"uint32":      {new(Value[uint32]), int64(42), 42},
		"uint16":      {new(Value[uint16]), int64(42), 42},
		"uint8":       {new(Value[uint8]), int64(42), 42},
		"from_uint":   {new(Value[uint64]), uint(42), 42},
		"from_uint32": {new(Value[uint64]), uint32(42), 42},
		"from_float":  {new(Value[uint64]), float64(42), 42},
	}

	for name, tt := range tests {
		s.Run(name, func() {
			var got uint64
			var err error
			switch v := tt.target.(type) {
			case *Value[uint]:
				err = v.Scan(tt.input)
				got = uint64(v.Get())
			case *Value[uint64]:
				err = v.Scan(tt.input)
				got = v.Get()
			case *Value[uint32]:
				err = v.Scan(tt.input)
				got = uint64(v.Get())
			case *Value[uint16]:
				err = v.Scan(tt.input)
				got = uint64(v.Get())
			case *Value[uint8]:
				err = v.Scan(tt.input)
				got = uint64(v.Get())
			}
			s.Require().NoError(err)
			s.Equal(tt.want, got)
		})
	}
}

func (s *SQLScannerSuite) TestScan_Uint_Error() {
	var v Value[uint64]
	err := v.Scan("not a number")
	s.Error(err)
}

func (s *SQLScannerSuite) TestScan_Float() {
	var vf64 Value[float64]
	err := vf64.Scan(float64(3.14))
	s.Require().NoError(err)
	s.Equal(3.14, vf64.Get())

	var vf32 Value[float32]
	err = vf32.Scan(float64(3.14))
	s.Require().NoError(err)
	s.Equal(float32(3.14), vf32.Get())
}

func (s *SQLScannerSuite) TestScan_Float_FromOtherTypes() {
	var v1 Value[float64]
	err := v1.Scan(float32(3.14))
	s.Require().NoError(err)

	var v2 Value[float64]
	err = v2.Scan(int64(42))
	s.Require().NoError(err)
	s.Equal(float64(42), v2.Get())
}

func (s *SQLScannerSuite) TestScan_Float_Error() {
	var v Value[float64]
	err := v.Scan("not a number")
	s.Error(err)
}

func (s *SQLScannerSuite) TestScan_Bool() {
	var vb Value[bool]
	err := vb.Scan(true)
	s.Require().NoError(err)
	s.True(vb.Get())
}

func (s *SQLScannerSuite) TestScan_Bool_FromInt() {
	var v1 Value[bool]
	err := v1.Scan(int64(1))
	s.Require().NoError(err)
	s.True(v1.Get())

	var v2 Value[bool]
	err = v2.Scan(int64(0))
	s.Require().NoError(err)
	s.False(v2.Get())
}

func (s *SQLScannerSuite) TestScan_Bool_Error() {
	var v Value[bool]
	err := v.Scan("not a bool")
	s.Error(err)
}

func (s *SQLScannerSuite) TestScan_Time() {
	now := time.Now().Truncate(time.Second)
	var vt Value[time.Time]
	err := vt.Scan(now)
	s.Require().NoError(err)
	s.True(vt.Get().Equal(now))
}

func (s *SQLScannerSuite) TestScan_Time_FromString() {
	timeStr := "2024-01-15T10:30:00Z"
	var vt Value[time.Time]
	err := vt.Scan(timeStr)
	s.Require().NoError(err)
	s.False(vt.Get().IsZero())
}

func (s *SQLScannerSuite) TestScan_Time_FromBytes() {
	timeStr := "2024-01-15T10:30:00Z"
	var vt Value[time.Time]
	err := vt.Scan([]byte(timeStr))
	s.Require().NoError(err)
	s.False(vt.Get().IsZero())
}

func (s *SQLScannerSuite) TestScan_Time_Errors() {
	var v1 Value[time.Time]
	s.Error(v1.Scan("not a time"))

	var v2 Value[time.Time]
	s.Error(v2.Scan([]byte("not a time")))

	var v3 Value[time.Time]
	s.Error(v3.Scan(123))
}

func (s *SQLScannerSuite) TestScan_Bytes() {
	var vb Value[[]byte]
	err := vb.Scan([]byte("hello"))
	s.Require().NoError(err)
	s.Equal("hello", string(vb.Get()))
}

func (s *SQLScannerSuite) TestScan_Bytes_FromString() {
	var v Value[[]byte]
	err := v.Scan("world")
	s.Require().NoError(err)
	s.Equal("world", string(v.Get()))
}

func (s *SQLScannerSuite) TestScan_Bytes_Error() {
	var v Value[[]byte]
	err := v.Scan(123)
	s.Error(err)
}

func (s *SQLScannerSuite) TestScan_Reflect() {
	type MyInt int

	// Assignable
	var v Value[MyInt]
	err := v.Scan(MyInt(42))
	s.Require().NoError(err)
	s.Equal(MyInt(42), v.Get())

	// Convertible
	var v2 Value[MyInt]
	err = v2.Scan(int(42))
	s.Require().NoError(err)
	s.Equal(MyInt(42), v2.Get())

	// Not convertible
	var v3 Value[MyInt]
	s.Error(v3.Scan("not convertible"))
}
