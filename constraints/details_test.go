package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestDetails(t *testing.T) {
	t.Run("should return no violations if the constraint argument has no violations", func(t *testing.T) {
		violations := Details(Equals("test"), "should not happen")(validation.NewContext("test"))
		assert.Len(t, violations, 0, "should not return any violations")
	})

	t.Run("should return a violation with the given message if the constraint argument has violations", func(t *testing.T) {
		violations := Details(Equals("test"), "should happen")(validation.NewContext("not test"))
		assert.Len(t, violations, 1, "should return a single violation")
		assert.Equal(t, "should happen", violations[0].Message, "should have the given message")
		assert.Len(t, violations[0].Details, 0, "should not have any details")
	})

	t.Run("should interpret variadic arguments as details kv map", func(t *testing.T) {
		violations := Details(Equals("test"), "should happen", "key", "value", "hello", 1)(validation.NewContext("not test"))
		assert.Len(t, violations, 1, "should return a single violation")
		assert.Equal(t, "should happen", violations[0].Message, "should have the given message")
		assert.Equal(t, map[string]any{
			"key":   "value",
			"hello": 1,
		}, violations[0].Details, "should have the given details")
	})
}
