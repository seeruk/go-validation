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

				ctx := ctx.WithValue(valueString(key), value)
				violations = append(violations, constraint.Violations(ctx)...)
			}
		case reflect.Array, reflect.Slice:
			for i := 0; i < rval.Len(); i++ {
				ctx := ctx.WithValue(fmt.Sprintf("[%d]", i), rval.Index(i))
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
		ctx := ctx.WithValue(FieldName(ctx, fieldName), rval.FieldByName(fieldName))
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
			ctx := ctx.WithValue(valueString(key), key).WithPathKind(PathKindKey)
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
	MustBe(UnwrapType(ctx.Value().Node.Type()), reflect.Struct)
	return f().Violations(ctx)
}

// LazyDynamic is a Constraint that's extremely similar to Lazy, and fulfils mostly the same
// purpose, but instead expects a function that has a single argument of the type being validated.
func LazyDynamic(constraintFn interface{}) Constraint {
	// TODO: Check constraintFn is not nil?
	return &lazyDynamic{
		constraintFn: constraintFn,
	}
}

// lazyDynamic is the implementation of the LazyDynamic constraint.
type lazyDynamic struct {
	constraintFn interface{}
}

// constraintType is kept on it's own here because it won't change, we don't need to fetch it every
// time a constraint is run that uses it.
var constraintType = reflect.TypeOf((*Constraint)(nil)).Elem()

// Violations ...
func (ld *lazyDynamic) Violations(ctx Context) []ConstraintViolation {
	rval := UnwrapValue(ctx.Value().Node)
	MustBe(UnwrapType(ctx.Value().Node.Type()), reflect.Struct)

	if IsNillable(rval) && rval.IsNil() {
		return nil
	}

	rfn := reflect.ValueOf(ld.constraintFn)
	rfnt := rfn.Type()

	if rfnt.NumIn() != 1 {
		panic("validation: LazyDynamic expects a function that accepts a single argument of the type being validated (or it's unwrapped value type)")
	}

	if rfnt.NumOut() != 1 || rfnt.Out(0) != constraintType {
		panic("validation: LazyDynamic expects a function that returns a Constraint")
	}

	isContextType := rfnt.In(0) == ctx.Value().Node.Type()
	isUnwrappedType := rfnt.In(0) == rval.Type()

	if !isContextType && !isUnwrappedType {
		panic("validation: LazyDynamic expects a function that accepts a single argument of the type being validated (or it's unwrapped value type)")
	}

	var constraint Constraint
	if isUnwrappedType {
		constraint = rfn.Call([]reflect.Value{rval})[0].Interface().(Constraint)
	} else {
		constraint = rfn.Call([]reflect.Value{ctx.Value().Node})[0].Interface().(Constraint)
	}

	return constraint.Violations(ctx)
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
		ctx := ctx.WithValue(
			valueString(reflect.ValueOf(mapKey)),
			rval.MapIndex(reflect.ValueOf(mapKey)),
		)

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

// valueString returns a string representation of the given value. It handles any type that may be
// nil by returning "nil", otherwise defers to fmt.Sprint.
func valueString(val reflect.Value) string {
	unwrapped := UnwrapValue(val)
	if IsNillable(unwrapped) && unwrapped.IsNil() {
		return "nil"
	}

	return fmt.Sprint(unwrapped)
}
