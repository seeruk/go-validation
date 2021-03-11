package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// MutuallyExclusive ...
// TODO: Support maps.
func MutuallyExclusive(fields ...string) validation.ConstraintFunc {
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

		var nonEmpty []string
		for _, field := range fields {
			f := rval.FieldByName(field)
			if !validation.IsEmpty(f) {
				nonEmpty = append(nonEmpty, validation.FieldName(ctx, field))
			}
		}

		if len(nonEmpty) > 1 {
			return []validation.ConstraintViolation{
				ctx.Violation("fields are mutually exclusive", map[string]interface{}{
					"fields": nonEmpty,
				}),
			}
		}

		return nil
	}
}
