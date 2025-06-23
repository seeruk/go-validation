package constraints

import (
	"fmt"
	"reflect"

	"github.com/seeruk/go-validation"
)

// OneOfKeys ...
func OneOfKeys[T any](keys ...T) validation.ConstraintFunc {
	if len(keys) < 1 {
		panic("constraints: OneOfKeys must be given at least 1 allowed value")
	}

	return ValueFunc(func(ctx validation.Context, rval reflect.Value) []validation.ConstraintViolation {
		// We don't want to be looping twice every time, so a map is made.
		allowed := make(map[any]struct{}, len(keys))
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
				ctx.Violation("key must be one of the allowed keys", map[string]any{
					"unexpected": unexpected,
				}),
			}
		}

		return nil
	}, reflect.Map)
}
