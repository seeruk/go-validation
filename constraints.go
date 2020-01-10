package validation

import (
	"fmt"
	"reflect"
)

// Constraints is simply a collection of many constraints. All of the constraints will be run, and
// their results will be aggregated and returned.
type Constraints []Constraint

// Violations ...
func (cc Constraints) Violations(ctx Context) []ConstraintViolation {
	var violations []ConstraintViolation
	for _, c := range cc {
		violations = append(violations, c.Violations(ctx)...)
	}

	return violations
}

// Elements is a Constraint used to validate every value (element) in an array, a map, or a slice.
type Elements []Constraint

// Violations ...
func (e Elements) Violations(ctx Context) []ConstraintViolation {
	rval := UnwrapValue(ctx.Value().Node)
	MustBe(UnwrapType(rval.Type()), reflect.Array, reflect.Map, reflect.Slice)

	var violations []ConstraintViolation
	if rval.IsZero() {
		return violations
	}

	if rval.Len() == 0 {
		return violations
	}

	for _, constraint := range e {
		switch rval.Kind() {
		case reflect.Map:
			iter := rval.MapRange()
			for iter.Next() {
				key := iter.Key()
				value := iter.Value()

				ctx := ctx.WithValue(Value{
					Name: valueString(key),
					Node: value,
				})
				violations = append(violations, constraint.Violations(ctx)...)
			}
		case reflect.Array, reflect.Slice:
			for i := 0; i < rval.Len(); i++ {
				ctx := ctx.WithValue(Value{
					Name: fmt.Sprintf("[%d]", i),
					Node: rval.Index(i),
				})
				violations = append(violations, constraint.Violations(ctx)...)
			}
		}
	}

	return violations
}

// Fields is a Constraint used to validate the values of specific fields on a struct.
type Fields map[string]Constraint

// Violations ...
func (f Fields) Violations(ctx Context) []ConstraintViolation {
	rval := UnwrapValue(ctx.Value().Node)
	rtyp := UnwrapType(rval.Type())
	MustBe(rtyp, reflect.Struct)

	var violations []ConstraintViolation
	if IsNillable(rval) && rval.IsNil() {
		return violations
	}

	for fieldName, constraint := range f {
		ctx := ctx.WithValue(Value{
			Name: FieldName(ctx, fieldName),
			Node: rval.FieldByName(fieldName),
		})

		violations = append(violations, constraint.Violations(ctx)...)
	}

	return violations
}

// Keys is a Constraint used to validate the keys of a map.
type Keys []Constraint

// Violations ...
func (k Keys) Violations(ctx Context) []ConstraintViolation {
	rval := UnwrapValue(ctx.Value().Node)
	MustBe(UnwrapType(rval.Type()), reflect.Map)

	var violations []ConstraintViolation
	if rval.IsNil() {
		return violations
	}

	if rval.Len() == 0 {
		return violations
	}

	for _, constraint := range k {
		for _, key := range rval.MapKeys() {
			ctx := ctx.WithValue(Value{
				Name: valueString(key),
				Node: key,
			}).WithPathKind(PathKindKey)

			violations = append(violations, constraint.Violations(ctx)...)
		}
	}

	return violations
}

// Lazy is a Constraint that allows a function that returns Constraints to be evaluated at
// validation-time. This enables things like defining constraints for recursive structures.
type Lazy func() Constraint

// Violations ...
func (f Lazy) Violations(ctx Context) []ConstraintViolation {
	return f().Violations(ctx)
}

// Map is a Constraint used to validate a map. This Constraint validates the values in the map, by
// specific keys. If you want to use the same validation on all keys of a map, use Elements instead.
// If you want to validate the keys of the map, use Keys instead.
type Map map[interface{}]Constraint

// Violations ...
func (m Map) Violations(ctx Context) []ConstraintViolation {
	rval := UnwrapValue(ctx.Value().Node)
	MustBe(UnwrapType(rval.Type()), reflect.Map)

	var violations []ConstraintViolation
	if rval.IsNil() { // Maps, or pointers to maps can be nil.
		return violations
	}

	for mapKey, constraint := range m {
		ctx := ctx.WithValue(Value{
			Name: valueString(reflect.ValueOf(mapKey)),
			Node: rval.MapIndex(reflect.ValueOf(mapKey)),
		})

		violations = append(violations, constraint.Violations(ctx)...)
	}

	return violations
}

// When conditionally runs some constraints. The predicate is set up at the time of creating the
// constraints. If you want a more dynamic approach, you should use WhenFn instead. You can also
// build up constraints programmatically and use the value being validated to build the constraints.
func When(predicate bool, constraints ...Constraint) ConstraintFunc {
	return func(ctx Context) []ConstraintViolation {
		var violations []ConstraintViolation
		if predicate {
			for _, c := range constraints {
				violations = append(violations, c.Violations(ctx)...)
			}
		}

		return violations
	}
}

// WhenFn lazily conditionally runs some constraints. The predicate function is called during the
// validation process. If you need an even more dynamic approach, you can also build up constraints
// programmatically and use the value being validated to build the constraints.
func WhenFn(predicateFn func() bool, constraints ...Constraint) ConstraintFunc {
	return func(ctx Context) []ConstraintViolation {
		var violations []ConstraintViolation
		if predicateFn() {
			for _, c := range constraints {
				violations = append(violations, c.Violations(ctx)...)
			}
		}

		return violations
	}
}
