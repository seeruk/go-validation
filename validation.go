package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Validate ...
func Validate(value interface{}, constraints ...Constraint) []ConstraintViolation {
	ctx := Context{Value: reflect.ValueOf(value)}
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

// ConstraintViolation ...
type ConstraintViolation struct {
	Path     string                 `json:"path"`
	PathKind PathKind               `json:"path_kind"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
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

// Context ...
type Context struct {
	PathKind PathKind
	Value    reflect.Value
}

// Violation provides a convenient way to produce a ConstraintViolation using the information found
// on the Context. If a custom violation is needed, one can always be made using the information on
// the Context manually.
func (c *Context) Violation(message string, details map[string]interface{}) ConstraintViolation {
	return ConstraintViolation{
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
func (c Context) WithValue(value reflect.Value) Context {
	c.Value = value
	return c
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
	mustBe(ctx.Value, reflect.Array, reflect.Map, reflect.Slice)

	var violations []ConstraintViolation
	for _, constraint := range e {
		switch ctx.Value.Kind() {
		case reflect.Map:
			for _, key := range ctx.Value.MapKeys() {
				ctx := ctx.WithValue(ctx.Value.MapIndex(key).Elem())
				violations = append(violations, constraint.Violations(ctx)...)
			}
		case reflect.Array, reflect.Slice:
			for i := 0; i < ctx.Value.Len(); i++ {
				ctx := ctx.WithValue(ctx.Value.Index(i))
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
	mustBe(ctx.Value, reflect.Struct)

	var violations []ConstraintViolation
	for fieldName, constraint := range f {
		ctx := ctx.WithValue(ctx.Value.FieldByName(fieldName))
		violations = append(violations, constraint.Violations(ctx)...)
	}

	return violations
}

// Keys ...
type Keys []Constraint

// Violations ...
func (k Keys) Violations(ctx Context) []ConstraintViolation {
	mustBe(ctx.Value, reflect.Map)

	var violations []ConstraintViolation
	for _, constraint := range k {
		for _, key := range ctx.Value.MapKeys() {
			ctx := ctx.WithValue(key).WithPathKind(PathKindKey)
			violations = append(violations, constraint.Violations(ctx)...)
		}
	}

	return violations
}

// Map ...
type Map map[interface{}]Constraint

// Violations ...
func (m Map) Violations(ctx Context) []ConstraintViolation {
	mustBe(ctx.Value, reflect.Map)

	var violations []ConstraintViolation
	for mapKey, constraint := range m {
		ctx := ctx.WithValue(ctx.Value.MapIndex(reflect.ValueOf(mapKey)))
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
func mustBe(value reflect.Value, kinds ...reflect.Kind) {
	if len(kinds) == 0 {
		return
	}

	for _, kind := range kinds {
		if value.Kind() == kind {
			return
		}
	}

	var kindNames []string
	for _, kind := range kinds {
		kindNames = append(kindNames, kind.String())
	}

	panic(fmt.Sprintf("validation: value must be one of: %s", strings.Join(kindNames, ", ")))
}
