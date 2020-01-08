package constraints

import "github.com/seeruk/go-validation"

// Required ...
var Required validation.ConstraintFunc = func(ctx validation.Context) []validation.ConstraintViolation {
	if ctx.Value.IsZero() {
		return []validation.ConstraintViolation{{
			Message: "a value is required",
		}}
	}
	// TODO: Do we want to do length checks too?
	return nil
}
