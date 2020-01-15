package constraints

import (
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeBefore(t *testing.T) {
	past := time.Date(1000, time.January, 1, 0, 0, 0, 0, time.UTC)
	present := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	future := time.Date(3000, time.January, 1, 0, 0, 0, 0, time.UTC)

	t.Run("should return no violations if the context's time if before the constraints 'before' time", func(t *testing.T) {
		violations := TimeBefore(future)(validation.NewContext(present))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the context's time is not before the constraints 'before' time", func(t *testing.T) {
		violations := TimeBefore(past)(validation.NewContext(present))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := TimeBefore(future)(validation.NewContext(time.Time{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the time the value should be after with a violation", func(t *testing.T) {
		violations := TimeBefore(past)(validation.NewContext(present))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"time": past.Format(time.RFC3339),
		}, violations[0].Details)
	})

	t.Run("should not panic if given a nil pointer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			TimeBefore(future)(validation.NewContext((*time.Time)(nil)))
		})
	})

	t.Run("should panic if given a value of the wrong type, even if it's empty", func(t *testing.T) {
		assert.Panics(t, func() { TimeBefore(future)(validation.NewContext("test")) })
		assert.Panics(t, func() { TimeBefore(future)(validation.NewContext(123)) })
		assert.Panics(t, func() { TimeBefore(future)(validation.NewContext(url.Values{})) })
		assert.Panics(t, func() { TimeBefore(future)(validation.NewContext(regexp.MustCompile("^test"))) })
	})
}
