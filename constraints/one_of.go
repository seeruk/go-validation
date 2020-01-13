package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// OneOf ...
func OneOf(allowed ...interface{}) validation.ConstraintFunc {
	if len(allowed) < 2 {
		panic("constraint: OneOf must be given at least 2 allowed values")
	}

	return func(ctx validation.Context) []validation.ConstraintViolation {

		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		var found bool
		for _, a := range allowed {
			if reflect.DeepEqual(ctx.Value().Node.Interface(), a) {
				found = true
				break
			}
		}

		if !found {
			return []validation.ConstraintViolation{
				ctx.Violation("value must be one of the allowed values", map[string]interface{}{
					"allowed": allowed,
				}),
			}
		}

		return nil
	}
}
