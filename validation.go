package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/seeruk/go-validation/protobuf"
	"github.com/seeruk/go-validation/validationpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DefaultNameStructTag is the default struct tag used to override the name of struct fields in the
// path that's output.
const DefaultNameStructTag = "validation"

// Validate executes the given constraint(s) against the given value, returning any violations of
// those constraints.
func Validate(value interface{}, constraints ...Constraint) []ConstraintViolation {
	return ValidateContext(NewContext(value), constraints...)
}

// CreateValidateFn allows the caller to create a customised validate function, using the given
// option(s), allowing them to avoid manually creating a context and using the simpler API while
// maintaining the ability to customise the validation context.
// TODO: Any more options to add?
func CreateValidateFunc(structTag string) func(value interface{}, constraints ...Constraint) []ConstraintViolation {
	return func(value interface{}, constraints ...Constraint) []ConstraintViolation {
		ctx := NewContext(value)
		ctx.StructTag = structTag
		return ValidateContext(ctx, constraints...)
	}
}

// ValidateContext is exactly like Validate, except it doesn't create a Context for you. This allows
// for more granular configuration provided by the Context type (and means we can avoid creating a
// Validator struct type to do this).
func ValidateContext(ctx Context, constraints ...Constraint) []ConstraintViolation {
	if !ctx.Value().Node.IsValid() {
		panic("validation: expected a valid type to be given (i.e. valid to Go's reflect library)")
	}

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
	PathKindValue PathKind = iota
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
	ctx := Context{StructTag: DefaultNameStructTag}
	return ctx.WithValue("", reflect.ValueOf(value))
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
		if val.Name != "" && i < len(c.Values[1:])-1 {
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

// WithPathKind returns a shallow copy of this Context with the given PathKind assigned, not
// modifying the original Context.
func (c Context) WithPathKind(pathKind PathKind) Context {
	c.PathKind = pathKind
	return c
}

// WithValue returns a shallow copy of this Context with the given value assigned, not modifying the
// original Context.
func (c Context) WithValue(name string, val reflect.Value) Context {
	value := Value{
		Name: name,
		Node: val,
	}

	// TODO: This would be far more efficient with a linked list probably?
	c.Values = append(c.Values, value)
	return c
}

// Value represents a value to be validated, and it's "name" (i.e. something we can use to build up
// a path to the value).
type Value struct {
	Name string
	Node reflect.Value
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

	name := fieldName

	if ctx.StructTag != "" {
		tag := field.Tag.Get(ctx.StructTag)
		switch tag {
		case "":
			// Split should never return an empty slice as long as the separator is not empty.
			split := strings.Split(tag, ",")
			name = split[0]
		case "-":
			name = ""
		}
	}

	return name
}

// MustBe will panic if the given reflect.Value is not one of the given reflect.Kind kinds.
func MustBe(typ reflect.Type, kinds ...reflect.Kind) {
	if len(kinds) == 0 {
		panic("validation: at least one kind must be given to MustBe")
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

	panic(fmt.Sprintf(
		"validation: value must be one of: %s, got %s",
		strings.Join(kindNames, ", "),
		typ.Kind(),
	))
}

// ShouldBe is the non-panicking alternative to MustBe. Instead of panicking it returns a slice of
// ConstraintViolation which can be directly returned from a Constraint.
func ShouldBe(ctx Context, typ reflect.Type, kinds ...reflect.Kind) []ConstraintViolation {
	if len(kinds) == 0 {
		panic("validation: at least one kind must be given to ShouldBe")
	}

	for _, kind := range kinds {
		if typ.Kind() == kind {
			return nil
		}
	}

	var kindNames []string
	for _, kind := range kinds {
		kindNames = append(kindNames, kind.String())
	}

	return []ConstraintViolation{
		ctx.Violation("value should be of one of the allowed kinds", map[string]interface{}{
			"allowed_kinds": kindNames,
		}),
	}
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

	switch val.Kind() {
	case reflect.Interface, reflect.Ptr:
		// Unwrap any pointer values, we'll only get here if it's not nil.
		// Also handle interface values, which we might get from dynamic structures.
		return UnwrapValue(val.Elem())
	}

	return val
}

// ConstraintViolationsToProto converts a slice of ConstraintViolations into a slice of the ProtoBuf
// representation of those ConstraintViolations, making them ready to use for example in a gRPC
// service in a similar way to how ConstraintViolation can already be used in JSON web services.
//
// They're returned as a slice of proto.Message so that they can be passed straight into something
// like a gRPC status without transformation.
func ConstraintViolationsToProto(violations []ConstraintViolation) []proto.Message {
	protoViolations := make([]proto.Message, 0, len(violations))

	for _, violation := range violations {
		protoViolations = append(protoViolations, &validationpb.ConstraintViolation{
			Path: violation.Path,
			// Currently these enum values are both just numbers, and both start at the same number,
			// and the values are in the same order.
			PathKind: validationpb.PathKind(violation.PathKind),
			Message:  violation.Message,
			Details:  protobuf.MapToStruct(violation.Details),
		})
	}

	return protoViolations
}

// ViolationsToStatus returns the given set of constraint violations as a gRPC status.
func ViolationsToStatus(violations []ConstraintViolation) *status.Status {
	sts, err := status.New(codes.InvalidArgument, "validation failed").
		WithDetails(ConstraintViolationsToProto(violations)...)

	if err != nil {
		return status.New(codes.Internal, "failed to generate status for validation failures")
	}

	return sts
}
