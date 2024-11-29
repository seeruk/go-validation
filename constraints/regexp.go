package constraints

import (
	"reflect"
	"regexp"

	"github.com/seeruk/go-validation"
)

// Regexp ...
func Regexp(pattern *regexp.Regexp) validation.ConstraintFunc {
	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		if !pattern.MatchString(rval.String()) {
			return []validation.ConstraintViolation{
				ctx.Violation("value must match regular expression", map[string]any{
					// TODO: Include actual value?
					"regexp": pattern.String(),
				}),
			}
		}

		return nil
	}, reflect.String)
}
