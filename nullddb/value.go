package nullddb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bjaus/null"
)

// Value wraps null.Value[T] and adds DynamoDB marshaling support.
type Value[T any] struct {
	null.Value[T]
}

// --- Constructors ---

// New creates a valid Value containing v.
func New[T any](v T) Value[T] {
	return Value[T]{null.New(v)}
}

// NewNull creates a Value that is explicitly null.
func NewNull[T any]() Value[T] {
	return Value[T]{null.NewNull[T]()}
}

// NewPtr creates a Value from a pointer.
func NewPtr[T any](p *T) Value[T] {
	return Value[T]{null.NewPtr(p)}
}

// From wraps an existing null.Value.
func From[T any](v null.Value[T]) Value[T] {
	return Value[T]{v}
}

// --- DynamoDB Marshaler ---

func (v Value[T]) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	if !v.IsValid() {
		return &types.AttributeValueMemberNULL{Value: true}, nil
	}

	val := any(v.Get())
	switch x := val.(type) {
	case string:
		return &types.AttributeValueMemberS{Value: x}, nil
	case int:
		return &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(x), 10)}, nil
	case int64:
		return &types.AttributeValueMemberN{Value: strconv.FormatInt(x, 10)}, nil
	case int32:
		return &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(x), 10)}, nil
	case int16:
		return &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(x), 10)}, nil
	case int8:
		return &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(x), 10)}, nil
	case uint:
		return &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(x), 10)}, nil
	case uint64:
		return &types.AttributeValueMemberN{Value: strconv.FormatUint(x, 10)}, nil
	case uint32:
		return &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(x), 10)}, nil
	case uint16:
		return &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(x), 10)}, nil
	case uint8:
		return &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(x), 10)}, nil
	case float64:
		return &types.AttributeValueMemberN{Value: strconv.FormatFloat(x, 'f', -1, 64)}, nil
	case float32:
		return &types.AttributeValueMemberN{Value: strconv.FormatFloat(float64(x), 'f', -1, 32)}, nil
	case bool:
		return &types.AttributeValueMemberBOOL{Value: x}, nil
	case []byte:
		return &types.AttributeValueMemberB{Value: x}, nil
	case time.Time:
		return &types.AttributeValueMemberS{Value: x.Format(time.RFC3339Nano)}, nil
	default:
		return nil, fmt.Errorf("nullddb: unsupported type %T", val)
	}
}

// --- DynamoDB Unmarshaler ---

func (v *Value[T]) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	if _, ok := av.(*types.AttributeValueMemberNULL); ok {
		*v = NewNull[T]()
		return nil
	}

	var target any = new(T)
	ptr := target

	switch p := ptr.(type) {
	case *string:
		if err := unmarshalString(p, av); err != nil {
			return err
		}
	case *int:
		var i int64
		if err := unmarshalInt64(&i, av); err != nil {
			return err
		}
		*p = int(i)
	case *int64:
		if err := unmarshalInt64(p, av); err != nil {
			return err
		}
	case *int32:
		var i int64
		if err := unmarshalInt64(&i, av); err != nil {
			return err
		}
		*p = int32(i)
	case *int16:
		var i int64
		if err := unmarshalInt64(&i, av); err != nil {
			return err
		}
		*p = int16(i)
	case *int8:
		var i int64
		if err := unmarshalInt64(&i, av); err != nil {
			return err
		}
		*p = int8(i)
	case *uint:
		var u uint64
		if err := unmarshalUint64(&u, av); err != nil {
			return err
		}
		*p = uint(u)
	case *uint64:
		if err := unmarshalUint64(p, av); err != nil {
			return err
		}
	case *uint32:
		var u uint64
		if err := unmarshalUint64(&u, av); err != nil {
			return err
		}
		*p = uint32(u)
	case *uint16:
		var u uint64
		if err := unmarshalUint64(&u, av); err != nil {
			return err
		}
		*p = uint16(u)
	case *uint8:
		var u uint64
		if err := unmarshalUint64(&u, av); err != nil {
			return err
		}
		*p = uint8(u)
	case *float64:
		if err := unmarshalFloat64(p, av); err != nil {
			return err
		}
	case *float32:
		var f float64
		if err := unmarshalFloat64(&f, av); err != nil {
			return err
		}
		*p = float32(f)
	case *bool:
		if err := unmarshalBool(p, av); err != nil {
			return err
		}
	case *[]byte:
		if err := unmarshalBytes(p, av); err != nil {
			return err
		}
	case *time.Time:
		if err := unmarshalTime(p, av); err != nil {
			return err
		}
	default:
		return fmt.Errorf("nullddb: unsupported type %T", ptr)
	}

	*v = New(*target.(*T))
	return nil
}

func unmarshalString(dst *string, av types.AttributeValue) error {
	switch x := av.(type) {
	case *types.AttributeValueMemberS:
		*dst = x.Value
	case *types.AttributeValueMemberN:
		*dst = x.Value
	default:
		return fmt.Errorf("nullddb: cannot unmarshal %T into string", av)
	}
	return nil
}

func unmarshalInt64(dst *int64, av types.AttributeValue) error {
	n, ok := av.(*types.AttributeValueMemberN)
	if !ok {
		return fmt.Errorf("nullddb: cannot unmarshal %T into int64", av)
	}
	i, err := strconv.ParseInt(n.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("nullddb: cannot parse %q as int64: %w", n.Value, err)
	}
	*dst = i
	return nil
}

func unmarshalUint64(dst *uint64, av types.AttributeValue) error {
	n, ok := av.(*types.AttributeValueMemberN)
	if !ok {
		return fmt.Errorf("nullddb: cannot unmarshal %T into uint64", av)
	}
	u, err := strconv.ParseUint(n.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("nullddb: cannot parse %q as uint64: %w", n.Value, err)
	}
	*dst = u
	return nil
}

func unmarshalFloat64(dst *float64, av types.AttributeValue) error {
	n, ok := av.(*types.AttributeValueMemberN)
	if !ok {
		return fmt.Errorf("nullddb: cannot unmarshal %T into float64", av)
	}
	f, err := strconv.ParseFloat(n.Value, 64)
	if err != nil {
		return fmt.Errorf("nullddb: cannot parse %q as float64: %w", n.Value, err)
	}
	*dst = f
	return nil
}

func unmarshalBool(dst *bool, av types.AttributeValue) error {
	b, ok := av.(*types.AttributeValueMemberBOOL)
	if !ok {
		return fmt.Errorf("nullddb: cannot unmarshal %T into bool", av)
	}
	*dst = b.Value
	return nil
}

func unmarshalBytes(dst *[]byte, av types.AttributeValue) error {
	switch x := av.(type) {
	case *types.AttributeValueMemberB:
		*dst = x.Value
	case *types.AttributeValueMemberS:
		*dst = []byte(x.Value)
	default:
		return fmt.Errorf("nullddb: cannot unmarshal %T into []byte", av)
	}
	return nil
}

func unmarshalTime(dst *time.Time, av types.AttributeValue) error {
	s, ok := av.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("nullddb: cannot unmarshal %T into time.Time", av)
	}
	t, err := time.Parse(time.RFC3339Nano, s.Value)
	if err != nil {
		return fmt.Errorf("nullddb: cannot parse %q as time: %w", s.Value, err)
	}
	*dst = t
	return nil
}
