package constraints

import (
	"math"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	t.Run("should return no violations if the minimum value is met or exceeded", func(t *testing.T) {
		assert.Empty(t, Min(1)(validation.NewContext(1)))
		assert.Empty(t, Min(0)(validation.NewContext(123)))
		assert.Empty(t, Min(0)(validation.NewContext(uint(123))))
		assert.Empty(t, Min(0)(validation.NewContext(123.456)))
	})

	t.Run("should return a violation if the minimum value is not met", func(t *testing.T) {
		assert.NotEmpty(t, Min(math.MaxFloat64)(validation.NewContext(123)))
		assert.NotEmpty(t, Min(math.MaxFloat64)(validation.NewContext(uint(123))))
		assert.NotEmpty(t, Min(math.MaxFloat64)(validation.NewContext(123.456)))
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := Min(1)(validation.NewContext(0))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the minimum value with a violation", func(t *testing.T) {
		violations := Min(math.MaxFloat64)(validation.NewContext(123))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"minimum": math.MaxFloat64,
		}, violations[0].Details)
	})

	t.Run("should not panic if given values of any regular numeric type", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Min(0)(validation.NewContext(123))
			Min(0)(validation.NewContext(int8(123)))
			Min(0)(validation.NewContext(int16(123)))
			Min(0)(validation.NewContext(int32(123)))
			Min(0)(validation.NewContext(int64(123)))
			Min(0)(validation.NewContext(uint(123)))
			Min(0)(validation.NewContext(uint8(123)))
			Min(0)(validation.NewContext(uint16(123)))
			Min(0)(validation.NewContext(uint32(123)))
			Min(0)(validation.NewContext(uint64(123)))
			Min(0)(validation.NewContext(float32(123.456)))
			Min(0)(validation.NewContext(123.456))
		})
	})

	t.Run("should not panic if given a nil pointer to a numeric type", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Min(1)(validation.NewContext((*int)(nil)))
		})
	})

	t.Run("should return violations if given a value of the wrong type, and the value is not empty", func(t *testing.T) {
		ctx := validation.NewContext("hello world")
		assert.Len(t, Min(1)(ctx), 1)
	})
}
