package constraints

import "github.com/seeruk/go-validation"

// Empty ...
var Empty validation.ConstraintFunc = func(ctx validation.Context) []validation.ConstraintViolation {
	if !validation.IsEmpty(ctx.Value().Node) {
		return []validation.ConstraintViolation{
			ctx.Violation("a value must not be provided", nil),
		}
	}
	return nil
}
