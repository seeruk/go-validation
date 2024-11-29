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
		assert.Equal(t, map[string]any{
			"actual":  2,
			"maximum": 1,
		}, violations[0].Details)
	})

	t.Run("should return violations if given a value of the wrong type, if not empty", func(t *testing.T) {
		ctx := validation.NewContext(123)
		assert.Len(t, MaxLength(1)(ctx), 1)
	})
}
