package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinLength(t *testing.T) {
	t.Run("should return no violations if the minimum length is met or exceeded", func(t *testing.T) {
		violations := MinLength(1)(validation.NewContext([]string{"test"}))
		assert.Len(t, violations, 0)
		violations = MinLength(3)(validation.NewContext([]string{"test", "test", "test", "test"}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the minimum length is not met", func(t *testing.T) {
		violations := MinLength(2)(validation.NewContext([]string{"test"}))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := MinLength(1)(validation.NewContext([]string{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the minimum length with a violation", func(t *testing.T) {
		violations := MinLength(2)(validation.NewContext([]string{"test"}))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"actual":  1,
			"minimum": 2,
		}, violations[0].Details)
	})

	t.Run("should not panic if given values of any type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MinLength(1)(validation.NewContext([1]int{1}))
			MinLength(1)(validation.NewContext(make(chan struct{})))
			MinLength(1)(validation.NewContext(map[string]interface{}{}))
			MinLength(1)(validation.NewContext([]string{}))
			MinLength(1)(validation.NewContext(""))
		})
	})

	t.Run("should not panic if given a nil pointer to a type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MinLength(1)(validation.NewContext((*chan struct{})(nil)))
			MinLength(1)(validation.NewContext((*map[string]string)(nil)))
			MinLength(1)(validation.NewContext((*[]string)(nil)))
			MinLength(1)(validation.NewContext((*string)(nil)))
		})
	})

	t.Run("should panic if given a value of the wrong type, even if it's empty", func(t *testing.T) {
		assert.Panics(t, func() { MinLength(1)(validation.NewContext(123)) })
		assert.Panics(t, func() { MinLength(1)(validation.NewContext(0)) })
	})

	t.Run("should return violations if given a value of the wrong type, even if it's empty, if strict types is false", func(t *testing.T) {
		ctx := validation.NewContext(123)
		ctx.StrictTypes = false

		assert.Len(t, MinLength(1)(ctx), 1)
	})
}
