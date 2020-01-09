package validation_test

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/seeruk/go-validation/constraints"
)

type testType struct {
	Text   string `validation:"text"`
	Number int    `validation:"number"`
}

func (tt testType) Constraints() validation.Constraint {
	return validation.Constraints{
		constraints.MutuallyExclusive("Text", "Number"),
		validation.Fields{
			"Text":   constraints.Required,
			"Number": constraints.Required,
		},
	}
}

func BenchmarkValidate(b *testing.B) {
	tt := testType{
		Text:   "Hello, World",
		Number: 12345678901234,
	}

	cc := tt.Constraints()

	b.ReportAllocs()
	b.ResetTimer()

	var violations []validation.ConstraintViolation
	for i := 0; i < b.N; i++ {
		violations = validation.Validate(tt, cc)
	}

	b.Log(violations)
	if len(violations) == 0 {
		b.Error("expected constraint violations")
	}
}
