package constraints

import "github.com/seeruk/go-validation"

// Required ...
var Required validation.ConstraintFunc = func() []validation.ConstraintViolation {
	return nil
}
