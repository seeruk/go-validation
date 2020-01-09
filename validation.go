package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// defaultNameStructTag is the default struct tag used to override the name of struct fields in the
// path that's output.
const defaultNameStructTag = "validation"

// Validate executes the given constraint(s) against the given value, returning any violations of
// those constraints.
func Validate(value interface{}, constraints ...Constraint) []ConstraintViolation {
	return ValidateContext(NewContext(value), constraints...)
}

// ValidateContext is exactly like Validate, except it doesn't create a Context for you. This allows
// for more granular configuration provided by the Context type (and means we can avoid creating a
// Validator struct type to do this).
func ValidateContext(ctx Context, constraints ...Constraint) []ConstraintViolation {
	return Constraints(constraints).Violations(ctx)
}

// Constraint represents a type that will validate a value and/or adjust the validation scope for
// further validation (e.g. validating a field on a struct, or an element in a slice).
type Constraint interface {
	Violations(ctx Context) []ConstraintViolation
}

// ConstraintFunc provides a convenient way of defining a Constraint as a function instead of a
// struct, keeping code more compact.
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

// PathKind enumerates the different possible kinds of values that we validate. This is used to
// remove the ambiguity in what is being validated in cases where the path could either refer to a
// key, or a value (e.g. you might want to validate a map's key, but the path to the key would be
// the same as the path to the value under that key).
type PathKind int

// MarshalJSON returns a JSON encoded version of the string representation of this PathKind.
func (k PathKind) MarshalJSON() ([]byte, error) {
	// TODO: Be less lazy, more performant.
	return json.Marshal(k.String())
}

// String returns the string representation of this PathKind.
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

// ConstraintViolation contains information to highlight a value failing to fulfil the requirements
// of a Constraint. It contains information to find the value that is failing, and how to resolve
// the violation.
type ConstraintViolation struct {
	Path     string                 `json:"path"`
	PathKind PathKind               `json:"path_kind"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// Context contains useful information for a Constraint, including the value(s) being validated.
type Context struct {
	PathKind  PathKind
	StructTag string
	Values    []Value
}

// NewContext returns a new Context, with a Value created for the given interface{} value.
func NewContext(value interface{}) Context {
	ctx := Context{StructTag: defaultNameStructTag}
	return ctx.WithValue(Value{
		Node: reflect.ValueOf(value),
	})
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

// WithPathKind returns a shallow copy of this Context with the given Pathkind assigned, not
// modifying the original Context.
func (c Context) WithPathKind(pathKind PathKind) Context {
	c.PathKind = pathKind
	return c
}

// WithValue returns a shallow copy of this Context with the given value assigned, not modifying the
// original Context.
func (c Context) WithValue(value Value) Context {
	c.Values = append(c.Values, value)
	return c
}

// Value represents a value to be validated, and it's "name" (i.e. something we can use to build up
// a path to the value).
type Value struct {
	Name string
	Node reflect.Value
}

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

// Lazy is a Constraint that allows a function that returns Constraints to be evaluated at
// validation-time. This enables things like defining constraints for recursive structures.
type Lazy func() Constraint

// Violations ...
func (f Lazy) Violations(ctx Context) []ConstraintViolation {
	return f().Violations(ctx)
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

// FieldName returns the output name for a field of the given field name. The provided Context's
// value must be a struct, or this function will panic. FieldName will unwrap the value on the given
// Context, like many Constraints do.
func FieldName(ctx Context, fieldName string) string {
	rval := UnwrapValue(ctx.Value().Node)
	rtyp := UnwrapType(rval.Type())
	MustBe(rtyp, reflect.Struct)

	field, ok := rtyp.FieldByName(fieldName)
	if !ok {
		// TODO: More info, like type? Let's see what the stack trace looks like first.
		panic(fmt.Sprintf("validation: field '%s' does not exist", fieldName))
	}

	tag := field.Tag.Get(ctx.StructTag)

	name := fieldName
	if tag != "" {
		// Split should never return an empty slice as long as the separator is not empty.
		split := strings.Split(tag, ",")
		name = split[0]
	}

	return name
}

// MustBe will panic if the given reflect.Value is not one of the given reflect.Kind kinds.
func MustBe(typ reflect.Type, kinds ...reflect.Kind) {
	if len(kinds) == 0 {
		return
	}

	for _, kind := range kinds {
		if typ.Kind() == kind {
			return
		}
	}

	var kindNames []string
	for _, kind := range kinds {
		kindNames = append(kindNames, kind.String())
	}

	panic(fmt.Sprintf("validation: value must be one of: %s", strings.Join(kindNames, ", ")))
}

// UnwrapType takes the given reflect.Type, and if it's a pointer gets the pointer element's type.
// This process is recursive, so we always end up with a non-pointer type at the end of the process.
func UnwrapType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		return UnwrapType(typ.Elem())
	}

	return typ
}

// UnwrapValue takes the given reflect.Value, and if it's a pointer gets the pointer element. This
// process is recursive, so we always end up with a non-pointer type at the end of the process.
func UnwrapValue(val reflect.Value) reflect.Value {
	if IsNillable(val) && val.IsNil() {
		return val
	}

	// Unwrap any pointer values, we'll only get here if it's not nil.
	if val.Kind() == reflect.Ptr {
		return UnwrapValue(val.Elem())
	}

	return val
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
