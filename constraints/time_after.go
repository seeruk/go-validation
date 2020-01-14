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
		validation.MustBe(validation.UnwrapType(rval.Type()), reflect.Struct)

		if validation.IsEmpty(rval) {
			return nil
		}

		switch v := rval.Interface().(type) {
		case time.Time:
			if !v.After(after) {
				return []validation.ConstraintViolation{
					ctx.Violation("value must be after time", map[string]interface{}{
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