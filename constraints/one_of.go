package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// OneOf ...
func OneOf(allowed ...any) validation.ConstraintFunc {
	if len(allowed) < 2 {
		panic("constraints: OneOf must be given at least 2 allowed values")
	}

	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		var found bool
		for _, a := range allowed {
			if reflect.DeepEqual(rval.Interface(), a) {
				found = true
				break
			}
		}

		if !found {
			return []validation.ConstraintViolation{
				ctx.Violation("value must be one of the allowed values", map[string]any{
					"allowed": allowed,
				}),
			}
		}

		return nil
	}
}
