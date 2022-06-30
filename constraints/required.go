package constraints

import "github.com/seeruk/go-validation"

// Required ...
var Required validation.ConstraintFunc = func(ctx validation.Context) []validation.ConstraintViolation {
	rval := validation.UnwrapValue(ctx.Value().Node)
	if validation.IsEmpty(rval) {
		return []validation.ConstraintViolation{
			ctx.Violation("a value is required", nil),
		}
	}
	return nil
}
