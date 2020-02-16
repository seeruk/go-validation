package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// ExactlyNRequired ...
func ExactlyNRequired(n int, fields ...string) validation.ConstraintFunc {
	if n < 1 {
		// Exactly 0 required is saying that all of them must not be set, that's not what this
		// constraint is for. Negative values also don't make any sense.
		panic("constraints: value of n given to ExactlyNRequired must be at least 1")
	}

	if n >= len(fields) {
		// - If n < len(fields) that's fine, e.g. exactly 2 of 4 fields must be set.
		// - If n == len(fields) that's not fine, e.g. exactly 3 of 3 fields must be set, doesn't
		//   make sense because you'd just be saying that exactly all of the fields must be set.
		//   You should just use the Required constraint instead.
		// - If n > len(fields) that's not fine, e.g. exactly 4 of 2 fields must be set, makes no
		//   sense because again, you'd never be able to set more fields than there are fields, so
		//   the constraint would always return any violations.
		panic("constraints: value of n given to ExactlyNRequired must be less than the number of fields")
	}

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

		if len(nonEmpty) != n {
			return []validation.ConstraintViolation{
				ctx.Violation("exact number of required fields not met", map[string]interface{}{
					"actual":   len(nonEmpty),
					"expected": n,
					"fields":   fieldNames,
				}),
			}
		}

		return nil
	}
}
