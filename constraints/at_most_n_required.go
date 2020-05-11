package constraints

import (
	"reflect"

	"github.com/seeruk/go-validation"
)

// AtMostNRequired ...
func AtMostNRequired(n int, fields ...string) validation.ConstraintFunc {
	if n < 1 {
		// At most 0 required is saying that all of them must not be set, that's not what this
		// constraint is for. Negative values also don't make any sense.
		panic("constraints: value of n given to AtMostNRequired must be at least 1")
	}

	if n >= len(fields) {
		// - If n < len(fields) that's fine, e.g. at most 2 of 4 fields must be set.
		// - If n == len(fields) that's not fine, e.g. at most 3 of 3 fields must be set, doesn't
		//   make sense because you'd just be saying that at most all of the fields can be set.
		//   You'd never get more than all of the fields being set, so never any violation.
		// - If n > len(fields) that's not fine, e.g. at most 4 of 2 fields must be set, makes no
		//   sense because again, you'd never be able to set more fields than there are fields, so
		//   the constraint would never return any violations.
		panic("constraints: value of n given to AtMostNRequired must be less than the number of fields")
	}

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

		if len(nonEmpty) > n {
			return []validation.ConstraintViolation{
				ctx.Violation("maximum number of required fields exceeded", map[string]interface{}{
					"actual":  len(nonEmpty),
					"maximum": n,
					"fields":  fieldNames,
				}),
			}
		}

		return nil
	}
}
