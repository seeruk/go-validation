package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MinLength ...
func MinLength(min int) validation.ConstraintFunc {
	allowed := []reflect.Kind{reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String}

	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		if rval.Len() < min {
			return []validation.ConstraintViolation{
				ctx.Violation("minimum length not met", map[string]any{
					"actual":  rval.Len(),
					"minimum": min,
				}),
			}
		}

		return nil
	}, allowed...)
}
