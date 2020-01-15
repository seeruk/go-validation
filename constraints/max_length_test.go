package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaxLength(t *testing.T) {
	t.Run("should return no violations is the max length is not exceeded", func(t *testing.T) {
		violations := MaxLength(1)(validation.NewContext([]string{"test"}))
		assert.Len(t, violations, 0)
		violations = MaxLength(3)(validation.NewContext([]string{"test", "test", "test"}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the max length is exceeded", func(t *testing.T) {
		violations := MaxLength(1)(validation.NewContext([]string{"foo", "bar"}))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := MaxLength(0)(validation.NewContext([]string{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the maximum length with a violation", func(t *testing.T) {
		violations := MaxLength(1)(validation.NewContext([]string{"foo", "bar"}))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"actual":  2,
			"maximum": 1,
		}, violations[0].Details)
	})

	t.Run("should not panic if given values of any type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MaxLength(1)(validation.NewContext([1]int{1}))
			MaxLength(1)(validation.NewContext(make(chan struct{})))
			MaxLength(1)(validation.NewContext(map[string]interface{}{}))
			MaxLength(1)(validation.NewContext([]string{}))
			MaxLength(1)(validation.NewContext(""))
		})
	})

	t.Run("should not panic if given a nil pointer to a type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MaxLength(1)(validation.NewContext((*chan struct{})(nil)))
			MaxLength(1)(validation.NewContext((*map[string]string)(nil)))
			MaxLength(1)(validation.NewContext((*[]string)(nil)))
			MaxLength(1)(validation.NewContext((*string)(nil)))
		})
	})

	t.Run("should panic if given a value of the wrong type, even if it's empty", func(t *testing.T) {
		assert.Panics(t, func() { MaxLength(1)(validation.NewContext(123)) })
		assert.Panics(t, func() { MaxLength(1)(validation.NewContext(0)) })
	})
}
