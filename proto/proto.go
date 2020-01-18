package proto

import (
	"fmt"

	"github.com/gogo/protobuf/types"
)

// MapToStruct converts a map of string to interface{} (e.g. a JSON object) to a ProtoBuf 'Struct'
// type, ready to be consumed by things like gRPC services.
func MapToStruct(in map[string]interface{}) *types.Struct {
	inLen := len(in)
	if inLen == 0 {
		return nil
	}

	fields := make(map[string]*types.Value, inLen)
	for k, v := range in {
		fields[k] = toValue(v)
	}

	return &types.Struct{Fields: fields}
}

// toValue converts a Go value to a ProtoBuf 'Value', if possible.
func toValue(v interface{}) *types.Value {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return &types.Value{
			Kind: &types.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &types.Value{
			Kind: &types.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		return &types.Value{
			Kind: &types.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	case map[string]interface{}:
		return &types.Value{
			Kind: &types.Value_StructValue{
				StructValue: MapToStruct(v),
			},
		}
	default:
		return &types.Value{
			Kind: &types.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}
