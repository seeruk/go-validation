package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Min ...
func Min(min float64) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)

		validation.MustBe(validation.UnwrapType(rval.Type()),
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
		)

		if validation.IsEmpty(rval) {
			return nil
		}

		var actual float64

		switch rval.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actual = float64(rval.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			actual = float64(rval.Uint())
		case reflect.Float32, reflect.Float64:
			actual = rval.Float()
		}

		if actual < min {
			return []validation.ConstraintViolation{
				ctx.Violation("minimum value not met", map[string]interface{}{
					"minimum": min,
				}),
			}
		}

		return nil
	}
}
