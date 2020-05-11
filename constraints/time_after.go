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
			if !v.After(after) {
				return []validation.ConstraintViolation{
					ctx.Violation("value must be after time", map[string]interface{}{
						// TODO: after_time and actual_time?
						"time": after.Format(time.RFC3339),
					}),
				}
			}
		}

		return nil
	}
}
