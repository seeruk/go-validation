package constraints

import "github.com/seeruk/go-validation"

// MutuallyExclusive ...
func MutuallyExclusive(fields ...string) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := ctx.Value().Node

		var nonEmpty []string
		for _, field := range fields {
			f := rval.FieldByName(field)
			if !validation.IsEmpty(f) {
				nonEmpty = append(nonEmpty, field)
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
