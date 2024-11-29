package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// ValueFuncFunc is a function used by ValueFunc to perform validation.
type ValueFuncFunc func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation

// ValueFunc is a helper-constraint that makes it easier to define custom constraints by just taking
// an unwrapped value after verifying it's the correct kind.
func ValueFunc(fn ValueFuncFunc, kinds ...reflect.Kind) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		rtyp := validation.UnwrapType(rval.Type())

		violations := validation.ShouldBe(ctx, rtyp, kinds...)
		if len(violations) > 0 {
			return violations
		}

		return fn(ctx, rval)
	}
}
