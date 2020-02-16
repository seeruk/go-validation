package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MaxLength ...
func MaxLength(max int) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		rtyp := validation.UnwrapType(rval.Type())

		allowed := []reflect.Kind{reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String}

		// Value must be able to be passed to 'len'.
		if ctx.StrictTypes {
			validation.MustBe(rtyp, allowed...)
		} else {
			violations := validation.ShouldBe(ctx, rtyp, allowed...)
			if len(violations) > 0 {
				return violations
			}
		}

		if validation.IsEmpty(rval) {
			return nil
		}

		if rval.Len() > max {
			return []validation.ConstraintViolation{
				ctx.Violation("maximum length exceeded", map[string]interface{}{
					"actual":  rval.Len(),
					"maximum": max,
				}),
			}
		}

		return nil
	}
}
