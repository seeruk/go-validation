package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Validate ...
func Validate(value interface{}, constraints ...Constraint) []ConstraintViolation {
	ctx := (Context{}).WithValue(Value{
		Node: reflect.ValueOf(value),
	})

	return Constraints(constraints).Violations(ctx)
}

// Constraint ...
type Constraint interface {
	Violations(ctx Context) []ConstraintViolation
}

// ConstraintFunc ...
type ConstraintFunc func(ctx Context) []ConstraintViolation

// Validate ...
func (c ConstraintFunc) Violations(ctx Context) []ConstraintViolation {
	return c(ctx)
}

// All possible PathKind values.
const (
	PathKindValue = iota
	PathKindKey
)

// PathKind ...
type PathKind int

// MarshalJSON ...
func (k PathKind) MarshalJSON() ([]byte, error) {
	// TODO: Be less lazy, more performant.
	return json.Marshal(k.String())
}

// String ...
func (k PathKind) String() string {
	switch k {
	case PathKindValue:
		return "value"
	case PathKindKey:
		return "key"
	default:
		return "unknown"
	}
}

// ConstraintViolation ...
type ConstraintViolation struct {
	Path     string                 `json:"path"`
	PathKind PathKind               `json:"path_kind"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// Context ...
type Context struct {
	PathKind PathKind
	Values   []Value
}

// Value gets the current value (the last value in Values).
func (c *Context) Value() Value {
	if len(c.Values) == 0 {
		panic("validation: Value called on Context with empty Values")
	}

	return c.Values[len(c.Values)-1]
}

// Violation provides a convenient way to produce a ConstraintViolation using the information found
// on the Context. If a custom violation is needed, one can always be made using the information on
// the Context manually.
func (c *Context) Violation(message string, details map[string]interface{}) ConstraintViolation {
	pathBuilder := strings.Builder{}
	pathBuilder.WriteString(".")

	// The first is skipped because if validation.Validate was used it won't have a name, we'll just
	// refer to it as ".".
	for i, val := range c.Values[1:] {
		pathBuilder.WriteString(val.Name)
		if i < len(c.Values[1:])-1 {
			pathBuilder.WriteString(".")
		}
	}

	return ConstraintViolation{
		Path:     pathBuilder.String(),
		PathKind: c.PathKind,
		Message:  message,
		Details:  details,
	}
}

// WithPathKind ...
func (c Context) WithPathKind(pathKind PathKind) Context {
	c.PathKind = pathKind
	return c
}

// WithValue returns a shallow copy of this Context with the given value assigned, not modifying the
// original context.
func (c Context) WithValue(value Value) Context {
	c.Values = append(c.Values, value)
	return c
}

// Value ...
type Value struct {
	Name string
	Node reflect.Value
}

// Constraints ...
type Constraints []Constraint

// Violations ...
func (cc Constraints) Violations(ctx Context) []ConstraintViolation {
	var violations []ConstraintViolation
	for _, c := range cc {
		violations = append(violations, c.Violations(ctx)...)
	}

	return violations
}

// Elements ...
type Elements []Constraint

// Violations ...
func (e Elements) Violations(ctx Context) []ConstraintViolation {
	mustBe(ctx.Value(), reflect.Array, reflect.Map, reflect.Slice)

	rval := ctx.Value().Node

	var violations []ConstraintViolation
	for _, constraint := range e {
		switch rval.Kind() {
		case reflect.Map:
			for _, key := range rval.MapKeys() {
				ctx := ctx.WithValue(Value{
					Name: valueString(key),
					Node: rval.MapIndex(key).Elem(),
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

// Fields ...
type Fields map[string]Constraint

// Violations ...
func (f Fields) Violations(ctx Context) []ConstraintViolation {
	mustBe(ctx.Value(), reflect.Struct)

	var violations []ConstraintViolation
	for fieldName, constraint := range f {
		ctx := ctx.WithValue(Value{
			Name: fieldName,
			Node: ctx.Value().Node.FieldByName(fieldName),
		})

		violations = append(violations, constraint.Violations(ctx)...)
	}

	return violations
}

// Keys ...
type Keys []Constraint

// Violations ...
func (k Keys) Violations(ctx Context) []ConstraintViolation {
	mustBe(ctx.Value(), reflect.Map)

	rval := ctx.Value().Node

	var violations []ConstraintViolation
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

// Map ...
type Map map[interface{}]Constraint

// Violations ...
func (m Map) Violations(ctx Context) []ConstraintViolation {
	mustBe(ctx.Value(), reflect.Map)

	rval := ctx.Value().Node

	var violations []ConstraintViolation
	for mapKey, constraint := range m {
		ctx := ctx.WithValue(Value{
			Name: valueString(reflect.ValueOf(mapKey)),
			Node: rval.MapIndex(reflect.ValueOf(mapKey)),
		})

		violations = append(violations, constraint.Violations(ctx)...)
	}

	return violations
}

// When ...
func When(predicate bool, constraints ...Constraint) Constraint {
	return ConstraintFunc(func(ctx Context) []ConstraintViolation {
		var violations []ConstraintViolation
		if predicate {
			for _, c := range constraints {
				violations = append(violations, c.Violations(ctx)...)
			}
		}

		return violations
	})
}

// mustBe ...
func mustBe(value Value, kinds ...reflect.Kind) {
	if len(kinds) == 0 {
		return
	}

	for _, kind := range kinds {
		if value.Node.Kind() == kind {
			return
		}
	}

	var kindNames []string
	for _, kind := range kinds {
		kindNames = append(kindNames, kind.String())
	}

	panic(fmt.Sprintf("validation: value must be one of: %s", strings.Join(kindNames, ", ")))
}

func unwrap(val interface{}) interface{} {
	return unwrapValue(reflect.ValueOf(val)).Interface()
}

func unwrapValue(val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		if val.IsNil() {
			return val
		}

		return unwrapValue(val.Elem())
	}

	return val
}

func valueString(val reflect.Value) string {
	unwrapped := unwrapValue(val)
	str := fmt.Sprint(unwrapped)
	switch unwrapped.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		if unwrapped.IsNil() {
			return "nil"
		}
	}

	return str
}
