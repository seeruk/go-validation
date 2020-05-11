package constraints

import (
	"reflect"
	"time"

	"github.com/seeruk/go-validation"
)

// TimeBefore ...
func TimeBefore(before time.Time) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		rtyp := validation.UnwrapType(rval.Type())

		violations := validation.ShouldBe(ctx, rtyp, reflect.Struct)
		if len(violations) > 0 {
			return violations
		}

		switch v := rval.Interface().(type) {
		case time.Time:
			if !v.Before(before) {
				return []validation.ConstraintViolation{
					ctx.Violation("value must be before time", map[string]interface{}{
						// TODO: before_time and actual_time?
						"time": before.Format(time.RFC3339),
					}),
				}
			}
		}

		return nil
	}
}
