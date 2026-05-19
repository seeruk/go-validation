package constraints

import (
	"errors"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoError(t *testing.T) {
	t.Run("should return no violations if the function returns no error", func(t *testing.T) {
		violations := NoError(func(value string) error {
			return nil
		}, "value is invalid")(validation.NewContext("test"))

		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the function returns an error", func(t *testing.T) {
		violations := NoError(func(value string) error {
			return errors.New("test error")
		}, "value is invalid")(validation.NewContext("test"))

		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := NoError(func(value string) error {
			return errors.New("test error")
		}, "value is invalid")(validation.NewContext(""))

		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the error with a violation", func(t *testing.T) {
		violations := NoError(func(value string) error {
			return errors.New("test error")
		}, "value is invalid")(validation.NewContext("test"))

		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"error": "test error",
		}, violations[0].Details)
	})

	t.Run("should return violations if given a value of the wrong type, and the value is not empty", func(t *testing.T) {
		violations := NoError(func(value string) error {
			return nil
		}, "value is invalid")(validation.NewContext(123))

		require.Len(t, violations, 1)
		assert.Equal(t, "value does not match expected type", violations[0].Message)
		assert.Equal(t, map[string]any{
			"expected": "string",
			"actual":   "int",
		}, violations[0].Details)
	})

	t.Run("should work with wrapped/pointer values", func(t *testing.T) {
		value := "test"
		violations := NoError(func(value string) error {
			return nil
		}, "value is invalid")(validation.NewContext(&value))

		assert.Len(t, violations, 0)
	})
}
