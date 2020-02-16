package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Kind ...
func Kind(allowed ...reflect.Kind) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		rtyp := validation.UnwrapType(rval.Type())

		violations := validation.ShouldBe(ctx, rtyp, allowed...)
		if len(violations) > 0 {
			return violations
		}

		return nil
	}
}
