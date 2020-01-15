package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLength(t *testing.T) {
	t.Run("should return no violations if the input's length is exactly the expected length", func(t *testing.T) {
		violations := Length(1)(validation.NewContext([]string{"test"}))
		assert.Len(t, violations, 0)
		violations = Length(3)(validation.NewContext([]string{"test", "test", "test"}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the input's length is not exactly the expected length", func(t *testing.T) {
		violations := Length(1)(validation.NewContext([]string{"hello", "world"}))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := Length(1)(validation.NewContext([]string{}))
		assert.Empty(t, violations)
	})

	t.Run("should return details about the expected length with a violation", func(t *testing.T) {
		violations := Length(1)(validation.NewContext([]string{"hello", "world"}))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"actual":   2,
			"expected": 1,
		}, violations[0].Details)
	})

	t.Run("should not panic if given values of any type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Length(1)(validation.NewContext([1]int{1}))
			Length(1)(validation.NewContext(make(chan struct{})))
			Length(1)(validation.NewContext(map[string]interface{}{}))
			Length(1)(validation.NewContext([]string{}))
			Length(1)(validation.NewContext(""))
		})
	})

	t.Run("should not panic if given a nil pointer to a type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Length(1)(validation.NewContext((*chan struct{})(nil)))
			Length(1)(validation.NewContext((*map[string]string)(nil)))
			Length(1)(validation.NewContext((*[]string)(nil)))
			Length(1)(validation.NewContext((*string)(nil)))
		})
	})

	t.Run("should panic if given a value of the wrong type, even if it's empty", func(t *testing.T) {
		assert.Panics(t, func() { Length(1)(validation.NewContext(123)) })
		assert.Panics(t, func() { Length(1)(validation.NewContext(0)) })
	})
}
