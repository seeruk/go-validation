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
	t.Run("should return a copy of the original context with the given PathKind", func(t *testing.T) {
		oldCtx := validation.NewContext("hello")
		oldCtx.PathKind = validation.PathKindValue

		newCtx := oldCtx.WithPathKind(validation.PathKindKey)

		assert.NotEqual(t, oldCtx, newCtx)
		assert.Equal(t, validation.PathKindKey, newCtx.PathKind)
	})
}

func TestContext_WithValue(t *testing.T) {
	t.Run("should return a copy of the original context with the given value", func(t *testing.T) {
		oldCtx := validation.NewContext("hello")
		newCtx := oldCtx.WithValue("subject", reflect.ValueOf("world"))

		assert.NotEmpty(t, oldCtx, newCtx)
		assert.Equal(t, newCtx.Value().Node.Interface(), "world")
	})
}

func TestFieldName(t *testing.T) {
	type testSubject struct {
		Test1 string
		Test2 string `validation:"test2"`
		Test3 string `json:"test3,omitempty"`
	}

	t.Run("should return the given field's output name", func(t *testing.T) {
		ts := testSubject{}
		ctx := validation.NewContext(ts)

		name := validation.FieldName(ctx, "Test1")

		assert.Equal(t, "Test1", name)
	})

	t.Run("should use the 'validation' struct tag if set (default tag)", func(t *testing.T) {
		ts := testSubject{}
		ctx := validation.NewContext(ts)

		name := validation.FieldName(ctx, "Test2")

		assert.Equal(t, "test2", name)
	})

	t.Run("should other struct tags if configured, with support for CSVs in tags", func(t *testing.T) {
		ts := testSubject{}
		ctx := validation.NewContext(ts)
		ctx.StructTag = "json"

		name := validation.FieldName(ctx, "Test3")

		assert.Equal(t, "test3", name)
	})

	t.Run("should panic if run on a type other than struct (or pointers to structs)", func(t *testing.T) {
		ctx := validation.NewContext("test")
		assert.Panics(t, func() {
			validation.FieldName(ctx, "test")
		})
	})

	t.Run("should panic if the given field doesn't exist on the struct", func(t *testing.T) {
		ctx := validation.NewContext(testSubject{})
		assert.Panics(t, func() {
			validation.FieldName(ctx, "ThisFieldDoesNotExist")
		})
	})
}

func TestMustBe(t *testing.T) {
	t.Run("should not panic if the given type is of one of the given kinds", func(t *testing.T) {
		assert.NotPanics(t, func() {
			validation.MustBe(reflect.TypeOf("test"), reflect.Array, reflect.String)
		})
	})

	t.Run("should panic if the given type is not one of the given kinds", func(t *testing.T) {
		assert.Panics(t, func() {
			validation.MustBe(reflect.TypeOf("test"), reflect.Array, reflect.Map)
		})
	})

	t.Run("should panic if no kinds are given", func(t *testing.T) {
		assert.Panics(t, func() {
			validation.MustBe(reflect.TypeOf("hello"))
		})
	})
}

func TestUnwrapType(t *testing.T) {
	t.Run("should find the root, non-pointer type of the given type", func(t *testing.T) {
		var wrapped ****string

		wrappedType := reflect.TypeOf(wrapped)
		require.Equal(t, reflect.Ptr, wrappedType.Kind())

		unwrappedType := validation.UnwrapType(wrappedType)
		assert.Equal(t, reflect.String, unwrappedType.Kind())
	})

	t.Run("should returned types that aren't pointers", func(t *testing.T) {
		var val string

		typ := reflect.TypeOf(val)
		assert.Equal(t, typ, validation.UnwrapType(typ))
	})
}

func TestUnwrapValue(t *testing.T) {
	t.Run("should return the underlying non-pointer value for the given pointer value", func(t *testing.T) {
		val := "hello"
		layer1 := &val
		layer2 := &layer1
		layer3 := &layer2
		layer4 := &layer3

		unwrapped := validation.UnwrapValue(reflect.ValueOf(layer4))

		assert.Equal(t, val, unwrapped.Interface())
	})

	t.Run("should return nil values as they are", func(t *testing.T) {
		var foo ***string

		assert.Equal(t, foo, validation.UnwrapValue(reflect.ValueOf(foo)).Interface())
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
