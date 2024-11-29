package constraints

import (
	"reflect"
	"time"

	"github.com/seeruk/go-validation"
)

// TimeBefore ...
func TimeBefore(before time.Time) validation.ConstraintFunc {
	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		switch v := rval.Interface().(type) {
		case time.Time:
			if !v.Before(before) {
				return []validation.ConstraintViolation{
					ctx.Violation("value must be before time", map[string]any{
						// TODO: before_time and actual_time?
						"time": before.Format(time.RFC3339),
					}),
				}
			}
		}

		return nil
	}, reflect.Struct)
}
