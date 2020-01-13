package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEquals(t *testing.T) {
	t.Run("should return no violations for a valid value", func(t *testing.T) {
		violations := Equals(1)(validation.NewContext(1))
		assert.Len(t, violations, 0)
		violations = Equals("hello")(validation.NewContext("hello"))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the exact length is not met", func(t *testing.T) {
		violations := Equals(1)(validation.NewContext(2))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := Equals(1)(validation.NewContext(0))
		assert.Len(t, violations, 0)
		violations = Equals(1)(validation.NewContext(""))
		assert.Len(t, violations, 0)
		violations = Equals(1)(validation.NewContext([]string{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the expected value with a violation", func(t *testing.T) {
		violations := Equals("test")(validation.NewContext("not test"))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"expected": "test",
		}, violations[0].Details)
	})

	t.Run("should not panic if given a nil pointer to a type 'len' can be called on", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Equals(1)(validation.NewContext((*chan struct{})(nil)))
			Equals(1)(validation.NewContext((*map[string]string)(nil)))
			Equals(1)(validation.NewContext((*[]string)(nil)))
			Equals(1)(validation.NewContext((*string)(nil)))
		})
	})
}
