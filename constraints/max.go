package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Max ...
func Max(max float64) validation.ConstraintFunc {
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

		var violated bool

		switch rval.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			violated = float64(rval.Int()) > max
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			violated = float64(rval.Uint()) > max
		case reflect.Float32, reflect.Float64:
			violated = rval.Float() > max
		}

		if violated {
			return []validation.ConstraintViolation{
				ctx.Violation("maximum value exceeded", map[string]interface{}{
					"maximum": max,
				}),
			}
		}

		return nil
	}
}
