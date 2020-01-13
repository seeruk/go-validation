package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Max ...
func Max(max float64) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsNillable(rval) && rval.IsNil() {
			return nil
		}

		validation.MustBe(validation.UnwrapType(rval.Type()),
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
		)

		var exceeded bool

		switch rval.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			exceeded = float64(rval.Int()) > max
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			exceeded = float64(rval.Uint()) > max
		case reflect.Float32, reflect.Float64:
			exceeded = rval.Float() > max
		}

		if exceeded {
			return []validation.ConstraintViolation{
				ctx.Violation("maximum value exceeded", map[string]interface{}{
					"maximum": max,
				}),
			}
		}

		return nil
	}
}
