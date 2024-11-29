package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Equals ...
func Equals(value any) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		if !reflect.DeepEqual(rval.Interface(), value) {
			return []validation.ConstraintViolation{
				ctx.Violation("value must equal expected value", map[string]any{
					"expected": value,
				}),
			}
		}

		return nil
	}
}
