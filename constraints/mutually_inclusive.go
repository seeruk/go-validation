package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MutuallyInclusive ...
// TODO: Support maps.
func MutuallyInclusive(fields ...string) validation.ConstraintFunc {
	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		fieldNames := make([]string, 0, len(fields))

		var nonEmpty []string
		for _, field := range fields {
			// We need to get all aliased field names, not the fields arg.
			fieldName := validation.FieldName(ctx, field)
			fieldNames = append(fieldNames, fieldName)

			f := rval.FieldByName(field)
			if !validation.IsEmpty(f) {
				nonEmpty = append(nonEmpty, fieldName)
			}
		}

		if len(nonEmpty) > 1 && len(nonEmpty) != len(fields) {
			return []validation.ConstraintViolation{
				ctx.Violation("fields are mutually inclusive", map[string]any{
					"fields": fieldNames,
				}),
			}
		}

		return nil
	}, reflect.Struct)
}
