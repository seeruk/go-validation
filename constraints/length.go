package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Length ...
func Length(length int) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)

		// Value must be able to be passed to 'len'.
		validation.MustBe(validation.UnwrapType(rval.Type()),
			reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String)

		if rval.Len() != length {
			return []validation.ConstraintViolation{
				ctx.Violation("exact length not met", map[string]interface{}{
					"expected": length,
				}),
			}
		}

		return nil
	}
}
