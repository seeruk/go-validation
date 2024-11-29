package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Max ...
func Max(max float64) validation.ConstraintFunc {
	allowed := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
	}

	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		var actual float64

		switch rval.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actual = float64(rval.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			actual = float64(rval.Uint())
		case reflect.Float32, reflect.Float64:
			actual = rval.Float()
		}

		if actual > max {
			return []validation.ConstraintViolation{
				ctx.Violation("maximum value exceeded", map[string]any{
					"actual":  actual,
					"maximum": max,
				}),
			}
		}

		return nil
	}, allowed...)
}
