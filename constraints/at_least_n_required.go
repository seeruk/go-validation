package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// AtLeastNRequired ...
func AtLeastNRequired(n int, fields ...string) validation.ConstraintFunc {
	if n < 1 {
		// At least 0 required is saying that at least none of the fields must be set, which is the
		// same as not using this constraint. Negative values also don't make sense.
		panic("constraints: value of n given to AtLeastNRequired must be at least 1")
	}

	if n >= len(fields) {
		// - If n < len(fields) that's fine, e.g. at least 2 of 4 fields must be set.
		// - If n == len(fields) that's not fine, e.g. at least 3 of 3 fields must be set, doesn't
		//   make sense because you'd just be saying that all of the given fields are required. You
		//   should just use the Required constraint on all of the fields instead then.
		// - If n > len(fields) that's not fine, e.g. at least 4 of 2 fields must be set, makes no
		//   sense, because you'd never be able to satisfy that requirement, and the constraint
		//   would always return a violation.
		panic("constraints: value of n given to AtLeastNRequired must be less than the number of fields")
	}

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

		if len(nonEmpty) < n {
			return []validation.ConstraintViolation{
				ctx.Violation("minimum number of required fields not met", map[string]interface{}{
					"minimum": n,
					"fields":  fieldNames,
				}),
			}
		}

		return nil
	}
}
