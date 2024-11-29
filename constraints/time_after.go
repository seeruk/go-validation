package constraints

import (
	"reflect"
	"time"

	"github.com/seeruk/go-validation"
)

// TimeAfter ...
func TimeAfter(after time.Time) validation.ConstraintFunc {
	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		switch v := rval.Interface().(type) {
		case time.Time:
			if !v.After(after) {
				return []validation.ConstraintViolation{
					ctx.Violation("value must be after time", map[string]any{
						// TODO: after_time and actual_time?
						"time": after.Format(time.RFC3339),
					}),
				}
			}
		}

		return nil
	}, reflect.Struct)
}
