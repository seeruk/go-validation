package constraints

import (
	"regexp"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexp(t *testing.T) {
	pattern := regexp.MustCompile("^Hello, ")

	t.Run("should return no violations if the value does match the given regexp", func(t *testing.T) {
		violations := Regexp(pattern)(validation.NewContext("Hello, World!"))
		assert.Len(t, violations, 0)
		violations = Regexp(pattern)(validation.NewContext("Hello, Go!"))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the value doesn't match the given regexp", func(t *testing.T) {
		violations := Regexp(pattern)(validation.NewContext("test"))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := Regexp(pattern)(validation.NewContext(""))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the pattern with a violation", func(t *testing.T) {
		violations := Regexp(pattern)(validation.NewContext("test"))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"regexp": "^Hello, ",
		}, violations[0].Details)
	})

	t.Run("should not panic if given a nil pointer", func(t *testing.T) {
		assert.NotPanics(t, func() { Regexp(pattern)(validation.NewContext((*string)(nil))) })
	})

	t.Run("should return violations if given a value of the wrong type, and the value is not empty", func(t *testing.T) {
		ctx := validation.NewContext(123)
		assert.Len(t, Regexp(pattern)(ctx), 1)
	})
}
