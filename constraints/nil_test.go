package constraints

import (
	"testing"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	t.Run("should return no violations if the value is nil", func(t *testing.T) {
		violations := Nil(validation.NewContext(([]string)(nil)))
		assert.Len(t, violations, 0)
		violations = Nil(validation.NewContext((*time.Time)(nil)))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the value is not nil", func(t *testing.T) {
		t.Run("empty slice", func(t *testing.T) {
			violations := Nil(validation.NewContext([]string{}))
			assert.Len(t, violations, 1)
		})

		t.Run("non-empty slice", func(t *testing.T) {
			violations := Nil(validation.NewContext([]string{"not", "nil"}))
			assert.Len(t, violations, 1)
		})

		t.Run("non-nil pointer", func(t *testing.T) {
			val := new(4)
			violations := Nil(validation.NewContext(val))
			assert.Len(t, violations, 1)
		})

		t.Run("non-nil double pointer", func(t *testing.T) {
			val := new(new(4))
			violations := Nil(validation.NewContext(val))
			assert.Len(t, violations, 1)
		})
	})
}
