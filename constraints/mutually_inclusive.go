package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MutuallyInclusive ...
func MutuallyInclusive(fields ...string) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		validation.MustBe(validation.UnwrapType(rval.Type()), reflect.Struct)

		if validation.IsEmpty(rval) {
			return nil
		}

		fieldNames := make([]string, 0, len(fields))

		var nonEmpty []string
		for _, field := range fields {
			// We need to get all of the aliased field names, not the fields arg.
			fieldName := validation.FieldName(ctx, field)
			fieldNames = append(fieldNames, fieldName)

			f := rval.FieldByName(field)
			if !validation.IsEmpty(f) {
				nonEmpty = append(nonEmpty, fieldName)
			}
		}

		if len(nonEmpty) > 1 && len(nonEmpty) != len(fields) {
			return []validation.ConstraintViolation{
				ctx.Violation("fields are mutually inclusive", map[string]interface{}{
					"fields": fieldNames,
				}),
			}
		}

		return nil
	}
}
