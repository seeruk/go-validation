package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Length ...
func Length(length int) validation.ConstraintFunc {
	allowed := []reflect.Kind{reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String}

	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		if rval.Len() != length {
			return []validation.ConstraintViolation{
				ctx.Violation("exact length not met", map[string]any{
					"actual":   rval.Len(),
					"expected": length,
				}),
			}
		}

		return nil
	}, allowed...)
}
