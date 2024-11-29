package constraints

import (
	"math"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMax(t *testing.T) {
	t.Run("should return no violations if the maximum value is not exceeded", func(t *testing.T) {
		assert.Empty(t, Max(math.MaxFloat64)(validation.NewContext(123)))
		assert.Empty(t, Max(math.MaxFloat64)(validation.NewContext(uint(123))))
		assert.Empty(t, Max(math.MaxFloat64)(validation.NewContext(123.456)))
		assert.Empty(t, Max(math.MaxFloat64)(validation.NewContext(math.MaxFloat64)))
	})

	t.Run("should return a violation if the maximum value is exceeded", func(t *testing.T) {
		assert.NotEmpty(t, Max(1)(validation.NewContext(123)))
		assert.NotEmpty(t, Max(1)(validation.NewContext(uint(123))))
		assert.NotEmpty(t, Max(1)(validation.NewContext(123.456)))
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := Max(-10)(validation.NewContext(0))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the maximum value with a violation", func(t *testing.T) {
		violations := Max(1)(validation.NewContext(123))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"actual":  float64(123),
			"maximum": float64(1),
		}, violations[0].Details)
	})

	t.Run("should not panic if given values of any regular numeric type", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Max(math.MaxFloat64)(validation.NewContext(123))
			Max(math.MaxFloat64)(validation.NewContext(int8(123)))
			Max(math.MaxFloat64)(validation.NewContext(int16(123)))
			Max(math.MaxFloat64)(validation.NewContext(int32(123)))
			Max(math.MaxFloat64)(validation.NewContext(int64(123)))
			Max(math.MaxFloat64)(validation.NewContext(uint(123)))
			Max(math.MaxFloat64)(validation.NewContext(uint8(123)))
			Max(math.MaxFloat64)(validation.NewContext(uint16(123)))
			Max(math.MaxFloat64)(validation.NewContext(uint32(123)))
			Max(math.MaxFloat64)(validation.NewContext(uint64(123)))
			Max(math.MaxFloat64)(validation.NewContext(float32(123.456)))
			Max(math.MaxFloat64)(validation.NewContext(123.456))
		})
	})

	t.Run("should not panic if given a nil pointer to a numeric type", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Max(1)(validation.NewContext((*int)(nil)))
		})
	})

	t.Run("should return violations if given a value of the wrong type, and the value is not empty", func(t *testing.T) {
		ctx := validation.NewContext("hello world")
		assert.Len(t, Max(1)(ctx), 1)
	})
}
