package protobuf

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

// MapToStruct converts a map of string to any (e.g. a JSON object) to a ProtoBuf 'Struct'
// type, ready to be consumed by things like gRPC services.
func MapToStruct(in map[string]any) *structpb.Struct {
	inLen := len(in)
	if inLen == 0 {
		return nil
	}

	fields := make(map[string]*structpb.Value, inLen)
	for k, v := range in {
		fields[k] = toValue(v)
	}

	return &structpb.Struct{Fields: fields}
}

// toValue converts a Go value to a ProtoBuf 'Value', if possible.
func toValue(v any) *structpb.Value {
	switch v := v.(type) {
	case nil:
		return &structpb.Value{
			Kind: &structpb.Value_NullValue{
				NullValue: structpb.NullValue_NULL_VALUE,
			},
		}
	case bool:
		return &structpb.Value{
			Kind: &structpb.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int16:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint16:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	case map[string]any:
		return &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: MapToStruct(v),
			},
		}
	default:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}
