package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MaxLength ...
func MaxLength(max int) validation.ConstraintFunc {
	allowed := []reflect.Kind{reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String}

	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		if rval.Len() > max {
			return []validation.ConstraintViolation{
				ctx.Violation("maximum length exceeded", map[string]any{
					"actual":  rval.Len(),
					"maximum": max,
				}),
			}
		}

		return nil
	}, allowed...)
}
