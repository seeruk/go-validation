package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MutuallyExclusive ...
// TODO: Support maps.
func MutuallyExclusive(fields ...string) validation.ConstraintFunc {
	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		var nonEmpty []string
		for _, field := range fields {
			f := rval.FieldByName(field)
			if !validation.IsEmpty(f) {
				nonEmpty = append(nonEmpty, validation.FieldName(ctx, field))
			}
		}

		if len(nonEmpty) > 1 {
			return []validation.ConstraintViolation{
				ctx.Violation("fields are mutually exclusive", map[string]any{
					"fields": nonEmpty,
				}),
			}
		}

		return nil
	}, reflect.Struct)
}
