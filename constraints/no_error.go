package constraints

import (
	"fmt"

	"github.com/seeruk/go-validation"
)

// NoError attempts to validate a value by passing it to a provided function that returns an error
// if the value is invalid, for example, doing something like parsing a URL with the stdlib.
func NoError[V any](fn func(V) error, message string) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		v, ok := rval.Interface().(V)
		if !ok {
			var vt V
			return []validation.ConstraintViolation{
				ctx.Violation("value does not match expected type", map[string]any{
					"expected": fmt.Sprintf("%T", vt),
					"actual":   fmt.Sprintf("%T", rval.Interface()),
				}),
			}
		}

		if err := fn(v); err != nil {
			return []validation.ConstraintViolation{
				ctx.Violation(message, map[string]any{
					"error": err.Error(),
				}),
			}
		}

		return nil
	}
}
