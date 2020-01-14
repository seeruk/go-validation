package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MinLength ...
func MinLength(min int) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)

		// Value must be able to be passed to 'len'.
		validation.MustBe(validation.UnwrapType(rval.Type()),
			reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String)

		if validation.IsEmpty(rval) {
			return nil
		}

		if rval.Len() < min {
			return []validation.ConstraintViolation{
				ctx.Violation("minimum length not met", map[string]interface{}{
					"minimum": min,
				}),
			}
		}

		return nil
	}
}
