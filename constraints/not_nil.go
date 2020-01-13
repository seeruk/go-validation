package constraints

import "github.com/seeruk/go-validation"

// NotNil ...
var NotNil validation.ConstraintFunc = func(ctx validation.Context) []validation.ConstraintViolation {
	rval := validation.UnwrapValue(ctx.Value().Node)
	if validation.IsNillable(rval) && rval.IsNil() {
		return []validation.ConstraintViolation{
			ctx.Violation("value must not be nil", nil),
		}
	}

	return nil
}
