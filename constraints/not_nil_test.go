package constraints

import (
	"testing"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestNotNil(t *testing.T) {
	t.Run("should return no violations for a valid value", func(t *testing.T) {
		violations := NotNil(validation.NewContext([]string{}))
		assert.Len(t, violations, 0)
		violations = NotNil(validation.NewContext(&time.Time{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the value is nil", func(t *testing.T) {
		violations := NotNil(validation.NewContext((*string)(nil)))
		assert.Len(t, violations, 1)
	})
}
