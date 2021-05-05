package constraints

import (
	"testing"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	t.Run("should return a violation if the value is not empty", func(t *testing.T) {
		foo := "foo"
		bar := make(chan struct{}, 1)
		bar <- struct{}{}

		assert.NotEmpty(t, Empty(validation.NewContext(true)))
		assert.NotEmpty(t, Empty(validation.NewContext(123)))
		assert.NotEmpty(t, Empty(validation.NewContext(123.456)))
		assert.NotEmpty(t, Empty(validation.NewContext("test")))
		assert.NotEmpty(t, Empty(validation.NewContext([]string{"test"})))
		assert.NotEmpty(t, Empty(validation.NewContext([1]int{1})))
		assert.NotEmpty(t, Empty(validation.NewContext(map[int]int{1: 2})))
		assert.NotEmpty(t, Empty(validation.NewContext(&foo)))
		assert.NotEmpty(t, Empty(validation.NewContext(time.Now())))
		assert.NotEmpty(t, Empty(validation.NewContext(bar)))
	})

	t.Run("should not return a violation if the value is empty", func(t *testing.T) {
		var foo *string
		bar := make(chan struct{}, 1)

		assert.Empty(t, Empty(validation.NewContext(false)))
		assert.Empty(t, Empty(validation.NewContext(0)))
		assert.Empty(t, Empty(validation.NewContext(0.0)))
		assert.Empty(t, Empty(validation.NewContext("")))
		assert.Empty(t, Empty(validation.NewContext([]string{})))
		assert.Empty(t, Empty(validation.NewContext(map[int]int{})))
		assert.Empty(t, Empty(validation.NewContext(foo)))
		assert.Empty(t, Empty(validation.NewContext(time.Time{})))
		assert.Empty(t, Empty(validation.NewContext(bar)))
	})
}
