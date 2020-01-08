package constraints

import "github.com/seeruk/go-validation"

// MutuallyExclusive ...
func MutuallyExclusive(fields ...string) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		var nonEmpty []string
		for _, field := range fields {
			f := ctx.Value.FieldByName(field)
			if !validation.IsEmpty(f) {
				nonEmpty = append(nonEmpty, field)
			}
		}

		if len(nonEmpty) > 1 {
			return []validation.ConstraintViolation{{
				Message: "fields are mutually exclusive",
			}}
		}

		return nil
	}
}
