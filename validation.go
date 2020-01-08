package validation

import "reflect"

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
	Message string `json:"message,omitempty"`
}

// Context ...
type Context struct {
	Value reflect.Value
	// TODO: Path.
}

type Fields map[string]Constraint

func (f Fields) Violations(ctx Context) []ConstraintViolation {
	return nil
}

type Constraints []Constraint

func (cc Constraints) Violations(ctx Context) []ConstraintViolation {
	var violations []ConstraintViolation
	for _, c := range cc {
		violations = append(violations, c.Violations(ctx)...)
	}
	// TODO: Length check and return nil? We should be consistent with what is returned if there are
	// no constraint violations.
	return violations
}

type Elements []Constraint

func (e Elements) Violations(ctx Context) []ConstraintViolation {
	return nil
}

type Keys []Constraint

func (k Keys) Violations(ctx Context) []ConstraintViolation {
	return nil
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
