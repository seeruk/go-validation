package validation_test

import (
	"reflect"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/seeruk/go-validation/constraints"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {

}

func TestValidateContext(t *testing.T) {

}

func TestConstraintFunc(t *testing.T) {

}

func TestPathKind_MarshalJSON(t *testing.T) {

}

func TestPathKind_String(t *testing.T) {

}

func TestNewContext(t *testing.T) {
	t.Run("should return a context with a reflection of the given value assigned", func(t *testing.T) {
		val := "hello"
		rval := reflect.ValueOf(val)

		ctx := validation.NewContext(val)

		// The reflect.Value is different, but the underlying value is the same.
		assert.Equal(t, rval.Interface(), ctx.Value().Node.Interface())
	})

	t.Run("should return a context with the default struct tag set", func(t *testing.T) {
		ctx := validation.NewContext("hello")
		assert.Equal(t, validation.DefaultNameStructTag, ctx.StructTag)
	})
}

func TestContext_Value(t *testing.T) {
	t.Run("should return the most recently set value", func(t *testing.T) {
		value1 := "hello"
		value2 := "world"
		value3 := "test"

		ctx := validation.NewContext(value1)
		assert.Equal(t, value1, ctx.Value().Node.Interface())
		ctx = ctx.WithValue("value2", reflect.ValueOf(value2))
		assert.Equal(t, value2, ctx.Value().Node.Interface())
		ctx = ctx.WithValue("value3", reflect.ValueOf(value3))
		assert.Equal(t, value3, ctx.Value().Node.Interface())
	})

	t.Run("should panic if there are no values", func(t *testing.T) {
		// This shouldn't be possible if you use NewContext.
		ctx := validation.Context{}
		assert.Panics(t, func() {
			ctx.Value()
		})
	})
}

// Benchmarking ...
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

	ctx := validation.NewContext(tt)
	//ctx.StructTag = ""

	b.ReportAllocs()
	b.ResetTimer()

	var violations []validation.ConstraintViolation
	for i := 0; i < b.N; i++ {
		violations = validation.ValidateContext(ctx, cc)
	}

	b.Log(violations)
	if len(violations) == 0 {
		b.Error("expected constraint violations")
	}
}
