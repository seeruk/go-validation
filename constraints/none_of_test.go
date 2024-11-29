package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoneOf(t *testing.T) {
	t.Run("should return no violations if the value is not one of the disallowed values", func(t *testing.T) {
		violations := NoneOf("test", "foo", "bar")(validation.NewContext("bla"))
		assert.Len(t, violations, 0)
		violations = NoneOf("foo", "test", "bar")(validation.NewContext("bla"))
		assert.Len(t, violations, 0)
		violations = NoneOf("foo", "bar", "test")(validation.NewContext("bla"))
		assert.Len(t, violations, 0)
		violations = NoneOf("foo", "bar", "baz", "test")(validation.NewContext("bla"))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the value is one of the disallowed values", func(t *testing.T) {
		violations := NoneOf("foo", "bar", "baz")(validation.NewContext("foo"))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := NoneOf("foo", "bar")(validation.NewContext(""))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the allowed values with a violation", func(t *testing.T) {
		violations := NoneOf("foo", "bar")(validation.NewContext("foo"))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"disallowed": []any{
				"foo",
				"bar",
			},
		}, violations[0].Details)
	})

	t.Run("should not panic if given a nil pointer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			NoneOf("hello", "world")(validation.NewContext((*chan struct{})(nil)))
			NoneOf("hello", "world")(validation.NewContext((*map[string]string)(nil)))
			NoneOf("hello", "world")(validation.NewContext((*[]string)(nil)))
			NoneOf("hello", "world")(validation.NewContext((*string)(nil)))
		})
	})

	t.Run("should panic if given less than two allowed values", func(t *testing.T) {
		assert.Panics(t, func() {
			NoneOf("test")(validation.NewContext("test"))
		})
	})

	t.Run("should work with wrapped/pointer values", func(t *testing.T) {
		val := "foo"
		violations := NoneOf("foo", "bar")(validation.NewContext(&val))
		assert.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"disallowed": []any{
				"foo",
				"bar",
			},
		}, violations[0].Details)
	})
}
