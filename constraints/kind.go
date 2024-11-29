package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// Kind ...
func Kind(allowed ...reflect.Kind) validation.ConstraintFunc {
	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		return nil
	}, allowed...)
}
