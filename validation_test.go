package validation_test

import (
	"reflect"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/seeruk/go-validation/constraints"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("should return violations if the value is invalid", func(t *testing.T) {
		violations := validation.Validate(0, constraints.Required)
		assert.Len(t, violations, 1)
	})

	t.Run("should return no violations if the value is valid", func(t *testing.T) {
		violations := validation.Validate(1, constraints.Required)
		assert.Len(t, violations, 0)
	})
}

func TestValidateContext(t *testing.T) {
	t.Run("should return violations if the value is invalid", func(t *testing.T) {
		violations := validation.ValidateContext(validation.NewContext(0), constraints.Required)
		assert.Len(t, violations, 1)
	})

	t.Run("should return no violations if the value is valid", func(t *testing.T) {
		violations := validation.ValidateContext(validation.NewContext(1), constraints.Required)
		assert.Len(t, violations, 0)
	})
}

func TestConstraintFunc_Violations(t *testing.T) {
	t.Run("should run the constraint function", func(t *testing.T) {
		constraint := validation.ConstraintFunc(func(ctx validation.Context) []validation.ConstraintViolation {
			return []validation.ConstraintViolation{ctx.Violation("test", nil)}
		})

		violations := constraint.Violations(validation.NewContext(123))

		require.Len(t, violations, 1)
		assert.Equal(t, "test", violations[0].Message)
	})
}

func TestPathKind_MarshalJSON(t *testing.T) {
	t.Run("should return the PathKind as a JSON string", func(t *testing.T) {
		pk := validation.PathKindValue

		bs, err := pk.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"value"`, string(bs))

		pk = validation.PathKindKey

		bs, err = pk.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"key"`, string(bs))
	})
}

func TestPathKind_String(t *testing.T) {
	t.Run("should return a string representation of the PathKind", func(t *testing.T) {
		pk := validation.PathKindValue
		assert.Equal(t, "value", pk.String())
		pk = validation.PathKindKey
		assert.Equal(t, "key", pk.String())
		pk = -1
		assert.Equal(t, "unknown", pk.String())
	})
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

func TestContext_Violation(t *testing.T) {
	t.Run("should return a new violation with the given message and details", func(t *testing.T) {
		message := "test violation"
		details := map[string]interface{}{
			"some": "value",
		}

		ctx := validation.NewContext("")
		violation := ctx.Violation(message, details)

		assert.Equal(t, message, violation.Message)
		assert.Equal(t, details, violation.Details)
	})

	t.Run("should build up the path on the violation", func(t *testing.T) {
		value := map[string]map[string]string{
			"Layer1": {
				"Layer2": "Test",
			},
		}

		ctx := validation.NewContext(value)
		ctx = ctx.WithValue("Layer1", reflect.ValueOf(value["Layer1"]))
		ctx = ctx.WithValue("Layer2", reflect.ValueOf(value["Layer1"]["Layer2"]))

		violation := ctx.Violation("", nil)

		assert.Equal(t, ".Layer1.Layer2", violation.Path)
	})

	t.Run("should set the PathKind from the Context on the returned violation", func(t *testing.T) {
		ctx := validation.NewContext("")
		violation := ctx.Violation("", nil)

		assert.Equal(t, ctx.PathKind, violation.PathKind)

		ctx = ctx.WithPathKind(validation.PathKindKey)
		violation = ctx.Violation("", nil)

		assert.Equal(t, ctx.PathKind, violation.PathKind)
	})
}

func TestContext_WithPathKind(t *testing.T) {

}

func TestContext_WithValue(t *testing.T) {

}

func TestFieldName(t *testing.T) {

}

func TestMustBe(t *testing.T) {

}

func TestUnwrapType(t *testing.T) {

}

func TestUnwrapValue(t *testing.T) {

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