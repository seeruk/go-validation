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

		validation.MustBe(validation.UnwrapType(rval.Type()), reflect.String)

		if !pattern.MatchString(rval.String()) {
			return []validation.ConstraintViolation{
				ctx.Violation("value must match regular expression", map[string]interface{}{
					"regexp": pattern.String(),
				}),
			}
		}

		return nil
	}
}
