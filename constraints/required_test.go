package constraints

import (
	"testing"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	t.Run("should not return a violation if the value is not empty", func(t *testing.T) {
		foo := "foo"
		bar := make(chan struct{}, 1)
		bar <- struct{}{}

		assert.Empty(t, Required(validation.NewContext(true)))
		assert.Empty(t, Required(validation.NewContext(123)))
		assert.Empty(t, Required(validation.NewContext(123.456)))
		assert.Empty(t, Required(validation.NewContext("test")))
		assert.Empty(t, Required(validation.NewContext([]string{"test"})))
		assert.Empty(t, Required(validation.NewContext([1]int{1})))
		assert.Empty(t, Required(validation.NewContext(map[int]int{1: 2})))
		assert.Empty(t, Required(validation.NewContext(&foo)))
		assert.Empty(t, Required(validation.NewContext(time.Now())))
		assert.Empty(t, Required(validation.NewContext(bar)))
	})

	t.Run("should return a violation if the value is empty", func(t *testing.T) {
		var foo *string
		bar := make(chan struct{}, 1)

		assert.NotEmpty(t, Required(validation.NewContext(false)))
		assert.NotEmpty(t, Required(validation.NewContext(0)))
		assert.NotEmpty(t, Required(validation.NewContext(0.0)))
		assert.NotEmpty(t, Required(validation.NewContext("")))
		assert.NotEmpty(t, Required(validation.NewContext([]string{})))
		assert.NotEmpty(t, Required(validation.NewContext(map[int]int{})))
		assert.NotEmpty(t, Required(validation.NewContext(foo)))
		assert.NotEmpty(t, Required(validation.NewContext(time.Time{})))
		assert.NotEmpty(t, Required(validation.NewContext(bar)))
	})
}
