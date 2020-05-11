package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MinLength ...
func MinLength(min int) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		rtyp := validation.UnwrapType(rval.Type())

		allowed := []reflect.Kind{reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String}

		// Value should be able to be passed to 'len'.
		violations := validation.ShouldBe(ctx, rtyp, allowed...)
		if len(violations) > 0 {
			return violations
		}

		if rval.Len() < min {
			return []validation.ConstraintViolation{
				ctx.Violation("minimum length not met", map[string]interface{}{
					"actual":  rval.Len(),
					"minimum": min,
				}),
			}
		}

		return nil
	}
}
