package constraints

import (
	"net/url"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtLeastNRequired(t *testing.T) {
	type testSubject struct {
		Field1 string
		Field2 int
		Field3 []string       `validation:"field3"`
		Field4 map[string]int `validation:"field4"`
	}

	constraint := AtLeastNRequired(2, "Field1", "Field2", "Field3", "Field4")

	t.Run("should return no violations for a valid value", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello", Field2: 1234567}
		ts2 := testSubject{Field3: []string{"test"}, Field4: map[string]int{"test": 123}}

		assert.Empty(t, constraint(validation.NewContext(ts1)))
		assert.Empty(t, constraint(validation.NewContext(ts2)))
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := AtLeastNRequired(1, "Field1", "Field2")(validation.NewContext((*testSubject)(nil)))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if minimum number of fields are not set", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello"}
		ts2 := testSubject{Field3: []string{"test"}}

		assert.NotEmpty(t, constraint(validation.NewContext(ts1)))
		assert.NotEmpty(t, constraint(validation.NewContext(ts2)))
	})

	t.Run("should return the fields that at least n of should be set in the violation details", func(t *testing.T) {
		ts := testSubject{Field1: "hello"}

		violations := constraint(validation.NewContext(ts))

		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"minimum": 2,
			"fields":  []string{"Field1", "Field2", "field3", "field4"},
		}, violations[0].Details)
	})

	t.Run("should return no violations if the value is nil", func(t *testing.T) {
		var ts *testSubject
		assert.Empty(t, constraint(validation.NewContext(ts)))
	})

	t.Run("should panic if given a value of the wrong type, even if it's empty", func(t *testing.T) {
		assert.Panics(t, func() { constraint(validation.NewContext("")) })
		assert.Panics(t, func() { constraint(validation.NewContext(0)) })
		assert.Panics(t, func() { constraint(validation.NewContext(url.Values{})) })
	})

	t.Run("should panic if the value if n is 0 or less", func(t *testing.T) {
		assert.Panics(t, func() { AtLeastNRequired(0, "test")(validation.NewContext(testSubject{})) })
		assert.Panics(t, func() { AtLeastNRequired(-10, "test")(validation.NewContext(testSubject{})) })
		assert.Panics(t, func() { AtLeastNRequired(-99999, "test")(validation.NewContext(testSubject{})) })
	})

	t.Run("should panic if number of fields passed to constraint doesn't exceed n", func(t *testing.T) {
		assert.Panics(t, func() { AtLeastNRequired(2, "test")(validation.NewContext(testSubject{})) })
	})
}
