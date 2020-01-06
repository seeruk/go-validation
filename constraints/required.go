package constraints

import "github.com/seeruk/go-validation"

// Required ...
var Required validation.Constraint = func() []validation.ConstraintViolation {
	return nil
}
