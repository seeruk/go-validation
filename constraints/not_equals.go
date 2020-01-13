package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// NotEquals ...
func NotEquals(value interface{}) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		if reflect.DeepEqual(rval.Interface(), value) {
			return []validation.ConstraintViolation{
				// TODO: Better message wording...
				ctx.Violation("value must not equal expected value", map[string]interface{}{
					"expected": value,
				}),
			}
		}

		return nil
	}
}
