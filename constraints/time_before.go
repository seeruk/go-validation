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
			if !v.Before(before) {
				return []validation.ConstraintViolation{
					ctx.Violation("value must be before time", map[string]interface{}{
						// TODO: before_time and actual_time?
						"time": before.Format(time.RFC3339),
					}),
				}
			}
		default:
			panic("constraints: value given to TimeBefore must be a time.Time (or pointer to)")
		}

		return nil
	}
}
