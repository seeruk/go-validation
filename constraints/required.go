package constraints

import "github.com/seeruk/go-validation"

// Required ...
var Required validation.ConstraintFunc = func(ctx validation.Context) []validation.ConstraintViolation {
	if validation.IsEmpty(ctx.Value().Node) {
		return []validation.ConstraintViolation{
			ctx.Violation("a value is required", nil),
		}
	}
	return nil
}
