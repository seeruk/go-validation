package constraints

import (
	"reflect"
	"regexp"

	"github.com/seeruk/go-validation"
)

// Regexp ...
func Regexp(pattern *regexp.Regexp) validation.ConstraintFunc {
	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		rtyp := validation.UnwrapType(rval.Type())

		violations := validation.ShouldBe(ctx, rtyp, reflect.String)
		if len(violations) > 0 {
			return violations
		}

		if !pattern.MatchString(rval.String()) {
			return []validation.ConstraintViolation{
				ctx.Violation("value must match regular expression", map[string]interface{}{
					// TODO: Include actual value?
					"regexp": pattern.String(),
				}),
			}
		}

		return nil
	}
}
