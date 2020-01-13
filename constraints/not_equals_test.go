package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotEquals(t *testing.T) {
	t.Run("should return no violations for a valid value", func(t *testing.T) {
		violations := NotEquals(1)(validation.NewContext(0))
		assert.Len(t, violations, 0)
		violations = NotEquals("hello")(validation.NewContext("goodbye"))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the values are equal", func(t *testing.T) {
		violations := NotEquals(1)(validation.NewContext(1))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := NotEquals(1)(validation.NewContext(0))
		assert.Len(t, violations, 0)
		violations = NotEquals(1)(validation.NewContext(""))
		assert.Len(t, violations, 0)
		violations = NotEquals(1)(validation.NewContext([]string{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the expected value with a violation", func(t *testing.T) {
		violations := NotEquals("test")(validation.NewContext("test"))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"expected": "test",
		}, violations[0].Details)
	})

	t.Run("should not panic if given a nil pointer to a type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			NotEquals(1)(validation.NewContext((*chan struct{})(nil)))
			NotEquals(1)(validation.NewContext((*map[string]string)(nil)))
			NotEquals(1)(validation.NewContext((*[]string)(nil)))
			NotEquals(1)(validation.NewContext((*string)(nil)))
		})
	})
}
