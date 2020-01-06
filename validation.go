package validation

type Constraint interface {
	Violations() []ConstraintViolation
}

type ConstraintFunc func() []ConstraintViolation

type ConstraintViolation struct {
}

type Struct []Constraint

func (s Struct) Violations() []ConstraintViolation {
	return nil
}

type Fields map[string]Constraint

func (f Fields) Violations() []ConstraintViolation {
	return nil
}

type Constraints []Constraint

func (c Constraints) Violations() []ConstraintViolation {
	return nil
}

type Elements []Constraint

func (e Elements) Violations() []ConstraintViolation {
	return nil
}

type Keys []Constraint

func (k Keys) Violations() []ConstraintViolation {
	return nil
}

func When(predicate bool, constraint ...Constraint) Constraint {
	return nil
}
