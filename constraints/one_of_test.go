package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneOf(t *testing.T) {
	t.Run("should return no violations if the value is one of the allowed values", func(t *testing.T) {
		violations := OneOf("test", "foo", "bar")(validation.NewContext("test"))
		assert.Len(t, violations, 0)
		violations = OneOf("foo", "test", "bar")(validation.NewContext("test"))
		assert.Len(t, violations, 0)
		violations = OneOf("foo", "bar", "test")(validation.NewContext("test"))
		assert.Len(t, violations, 0)
		violations = OneOf("foo", "bar", "baz", "test")(validation.NewContext("test"))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the value isn't one of the allowed values", func(t *testing.T) {
		violations := OneOf("foo", "bar", "baz")(validation.NewContext("test"))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := OneOf("foo", "bar")(validation.NewContext(""))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the allowed values with a violation", func(t *testing.T) {
		violations := OneOf("foo", "bar")(validation.NewContext([]string{"test"}))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"allowed": []any{
				"foo",
				"bar",
			},
		}, violations[0].Details)
	})

	t.Run("should not panic if given a nil pointer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			OneOf("hello", "world")(validation.NewContext((*chan struct{})(nil)))
			OneOf("hello", "world")(validation.NewContext((*map[string]string)(nil)))
			OneOf("hello", "world")(validation.NewContext((*[]string)(nil)))
			OneOf("hello", "world")(validation.NewContext((*string)(nil)))
		})
	})

	t.Run("should panic if given less than two allowed values", func(t *testing.T) {
		assert.Panics(t, func() {
			OneOf("test")(validation.NewContext("test"))
		})
	})

	t.Run("should work with wrapped/pointer values", func(t *testing.T) {
		val := "foo"
		violations := OneOf("foo", "bar")(validation.NewContext(&val))
		require.Len(t, violations, 0)
	})
}
