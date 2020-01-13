package constraints

import "github.com/seeruk/go-validation"

// Nil ...
var Nil validation.ConstraintFunc = func(ctx validation.Context) []validation.ConstraintViolation {
	rval := validation.UnwrapValue(ctx.Value().Node)
	if validation.IsNillable(rval) && !rval.IsNil() {
		return []validation.ConstraintViolation{
			ctx.Violation("value must be nil", nil),
		}
	}

	return nil
}
