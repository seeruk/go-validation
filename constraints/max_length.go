package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MaxLength ...
func MaxLength(max int) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		// Value must be able to be passed to 'len'.
		validation.MustBe(validation.UnwrapType(rval.Type()),
			reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String)

		if rval.Len() > max {
			return []validation.ConstraintViolation{
				ctx.Violation("maximum length exceeded", map[string]interface{}{
					"maximum": max,
				}),
			}
		}

		return nil
	}
}
