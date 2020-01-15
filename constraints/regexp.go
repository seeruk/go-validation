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
		validation.MustBe(validation.UnwrapType(rval.Type()), reflect.String)

		if validation.IsEmpty(rval) {
			return nil
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
