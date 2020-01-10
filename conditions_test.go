package validation

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	t.Run("should return true if the value is invalid", func(t *testing.T) {
		assert.True(t, IsEmpty(reflect.ValueOf((fmt.Stringer)(nil))))
	})

	t.Run("should return true if the value is a zero value", func(t *testing.T) {
		assert.True(t, IsEmpty(reflect.ValueOf("")))
		assert.True(t, IsEmpty(reflect.ValueOf(0)))
		assert.True(t, IsEmpty(reflect.ValueOf(0.0)))
		assert.True(t, IsEmpty(reflect.ValueOf(false)))
		assert.True(t, IsEmpty(reflect.ValueOf([]string{})))
		assert.True(t, IsEmpty(reflect.ValueOf(map[string]string{})))
		assert.True(t, IsEmpty(reflect.ValueOf(struct{}{})))
		assert.True(t, IsEmpty(reflect.ValueOf((chan struct{})(nil))))
		assert.True(t, IsEmpty(reflect.ValueOf(time.Time{})))
	})

	t.Run("should return false is the value is not empty", func(t *testing.T) {
		ch := make(chan string, 1)
		ch <- "Hello, World!"

		assert.False(t, IsEmpty(reflect.ValueOf("hello")))
		assert.False(t, IsEmpty(reflect.ValueOf(123)))
		assert.False(t, IsEmpty(reflect.ValueOf(123.456)))
		assert.False(t, IsEmpty(reflect.ValueOf(true)))
		assert.False(t, IsEmpty(reflect.ValueOf([]string{"Hello", "World"})))
		assert.False(t, IsEmpty(reflect.ValueOf(map[string]string{"foo": "bar"})))
		assert.False(t, IsEmpty(reflect.ValueOf(ch)))
		assert.False(t, IsEmpty(reflect.ValueOf(time.Now())))

		<-ch
	})
}

func TestIsNillable(t *testing.T) {
	t.Run("should return true for nillable types", func(t *testing.T) {
		assert.True(t, IsNillable(reflect.ValueOf((chan struct{})(nil))))
		assert.True(t, IsNillable(reflect.ValueOf((func())(nil))))
		assert.True(t, IsNillable(reflect.ValueOf((map[string]interface{})(nil))))
		assert.True(t, IsNillable(reflect.ValueOf((*string)(nil))))
		assert.True(t, IsNillable(reflect.ValueOf(([]string)(nil))))
		assert.True(t, IsNillable(reflect.ValueOf((unsafe.Pointer)(nil))))
	})

	t.Run("should return false for non-nillable types", func(t *testing.T) {
		assert.False(t, IsNillable(reflect.ValueOf("hello")))
		assert.False(t, IsNillable(reflect.ValueOf(123)))
		assert.False(t, IsNillable(reflect.ValueOf(123.456)))
		assert.False(t, IsNillable(reflect.ValueOf(true)))
		assert.False(t, IsNillable(reflect.ValueOf(struct{}{})))
	})
}
