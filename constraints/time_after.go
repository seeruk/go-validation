package constraints

import (
	"reflect"
	"time"

	"github.com/seeruk/go-validation"
)

// TimeAfter ...
func TimeAfter(after time.Time) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		rtyp := validation.UnwrapType(rval.Type())

		if ctx.StrictTypes {
			validation.MustBe(rtyp, reflect.Struct)
		} else {
			violations := validation.ShouldBe(ctx, rtyp, reflect.Struct)
			if len(violations) > 0 {
				return violations
			}
		}

		if validation.IsEmpty(rval) {
			return nil
		}

		switch v := rval.Interface().(type) {
		case time.Time:
			if !v.After(after) {
				return []validation.ConstraintViolation{
					ctx.Violation("value must be after time", map[string]interface{}{
						// TODO: after_time and actual_time?
						"time": after.Format(time.RFC3339),
					}),
				}
			}
		default:
			panic("constraints: value given to TimeAfter must be a time.Time (or pointer to)")
		}

		return nil
	}
}
