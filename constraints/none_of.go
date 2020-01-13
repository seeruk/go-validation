package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// NoneOf ...
func NoneOf(disallowed ...interface{}) validation.ConstraintFunc {
	if len(disallowed) < 2 {
		panic("constraint: NoneOf must be given at least 2 disallowed values")
	}

	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		var found bool
		for _, a := range disallowed {
			if reflect.DeepEqual(ctx.Value().Node.Interface(), a) {
				found = true
				break
			}
		}

		if found {
			return []validation.ConstraintViolation{
				ctx.Violation("value must not be one of the disallowed values", map[string]interface{}{
					"disallowed": disallowed,
				}),
			}
		}

		return nil
	}
}
