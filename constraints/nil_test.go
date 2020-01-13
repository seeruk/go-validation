package constraints

import (
	"testing"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	t.Run("should return no violations for a valid value", func(t *testing.T) {
		violations := Nil(validation.NewContext(([]string)(nil)))
		assert.Len(t, violations, 0)
		violations = Nil(validation.NewContext((*time.Time)(nil)))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the value is not nil", func(t *testing.T) {
		violations := Nil(validation.NewContext([]string{}))
		assert.Len(t, violations, 1)
	})
}
