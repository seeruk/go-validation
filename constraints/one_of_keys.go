package constraints

import (
	"fmt"
	"reflect"

	"github.com/seeruk/go-validation"
)

// OneOfKeys ...
func OneOfKeys(keys ...interface{}) validation.ConstraintFunc {
	if len(keys) < 1 {
		panic("constraints: OneOfKeys must be given at least 1 allowed value")
	}

	return func(ctx validation.Context) []validation.ConstraintViolation {
		rval := validation.UnwrapValue(ctx.Value().Node)
		if validation.IsEmpty(rval) {
			return nil
		}

		rtyp := validation.UnwrapType(rval.Type())

		violations := validation.ShouldBe(ctx, rtyp, reflect.Map)
		if len(violations) > 0 {
			return violations
		}

		// We don't want to be looping twice every time, so a map is made.
		allowed := make(map[interface{}]struct{}, len(keys))
		for _, k := range keys {
			allowed[k] = struct{}{}
		}

		var unexpected []string

		iter := rval.MapRange()
		for iter.Next() {
			key := iter.Key()
			if _, ok := allowed[key.Interface()]; !ok {
				unexpected = append(unexpected, fmt.Sprint(key.Interface()))
			}
		}

		if len(unexpected) > 0 {
			return []validation.ConstraintViolation{
				ctx.Violation("key must be one of the allowed keys", map[string]interface{}{
					"unexpected": unexpected,
				}),
			}
		}

		return nil
	}
}
