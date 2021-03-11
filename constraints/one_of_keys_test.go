package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneOfKeys(t *testing.T) {
	t.Run("should return no violations if the key is one of the allowed keys", func(t *testing.T) {
		violations := OneOfKeys("foo", 12, true)(validation.NewContext(map[interface{}]interface{}{
			"foo": "bar",
			12:    34,
			true:  false,
		}))

		assert.Len(t, violations, 0)
	})

	t.Run("should return a violations if a key isn't one of the allowed keys", func(t *testing.T) {
		violations := OneOfKeys("foo", 12, true)(validation.NewContext(map[interface{}]interface{}{
			false: true,
		}))

		assert.Len(t, violations, 1)

		violations = OneOfKeys("foo", 12, true)(validation.NewContext(map[interface{}]interface{}{
			"bar": "foo",
			34:    12,
			false: true,
		}))

		require.Len(t, violations, 1)
		assert.Len(t, violations[0].Details["unexpected"].([]string), 3)
	})

	t.Run("should not return any violations if the map is empty", func(t *testing.T) {
		violations := OneOfKeys("foo", 12, true)(validation.NewContext(map[interface{}]interface{}{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should not panic if given a nil map", func(t *testing.T) {
		assert.NotPanics(t, func() {
			OneOfKeys("hello", "world")(validation.NewContext((map[interface{}]interface{})(nil)))
		})
	})

	t.Run("should panic if given no allowed keys", func(t *testing.T) {
		assert.Panics(t, func() {
			OneOfKeys()(validation.NewContext(map[interface{}]interface{}{}))
		})
	})
}
