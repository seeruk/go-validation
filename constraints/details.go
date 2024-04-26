package constraints

import "github.com/seeruk/go-validation"

// Details allows you to provide a custom violation message and details for a constraint. This can
// be used to provide purpose-specific messages and details for constraints, as opposed to the
// generic messaging that constraints typically provide.
func Details(c validation.Constraint, msg string, details ...any) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		violations := c.Violations(ctx)
		if len(violations) == 0 {
			return nil
		}
		return []validation.ConstraintViolation{
			ctx.Violation(msg, detailsMap(details...)),
		}
	}
}

// detailsMap converts a variadic list of key/value pairs into a map.
func detailsMap(details ...any) map[string]any {
	m := make(map[string]any, len(details)/2)
	for i := 0; i < len(details); i += 2 {
		key := details[i].(string)
		value := details[i+1]
		m[key] = value
	}
	return m
}
