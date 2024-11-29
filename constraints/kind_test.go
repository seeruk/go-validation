package constraints

import (
	"reflect"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKind(t *testing.T) {
	t.Run("should return no violations if the value is of one of the allowed kinds", func(t *testing.T) {
		assert.Empty(t, Kind(reflect.String)(validation.NewContext("hello world")))
		assert.Empty(t, Kind(reflect.Int)(validation.NewContext(123)))
	})

	t.Run("should return no violations if the value is empty, even if it's the wrong type", func(t *testing.T) {
		assert.Empty(t, Kind(reflect.String)(validation.NewContext(0)))
	})

	t.Run("should return a violation if the value is not one of the allowed kinds", func(t *testing.T) {
		assert.Len(t, Kind(reflect.Struct)(validation.NewContext("hello world")), 1)
		assert.Len(t, Kind(reflect.Struct)(validation.NewContext(123)), 1)
	})

	t.Run("should return details about the allowed kinds", func(t *testing.T) {
		violations := Kind(reflect.Struct)(validation.NewContext(123))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"allowed_kinds": []string{
				"struct",
			},
		}, violations[0].Details)
	})
}
