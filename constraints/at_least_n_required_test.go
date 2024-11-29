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

	t.Run("should return no violations if minimum number of fields is met", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello", Field2: 1234567}
		ts2 := testSubject{Field3: []string{"test"}, Field4: map[string]int{"test": 123}}

		assert.Empty(t, constraint(validation.NewContext(ts1)))
		assert.Empty(t, constraint(validation.NewContext(ts2)))
	})

	t.Run("should return a violation if minimum number of fields is not met", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello"}
		ts2 := testSubject{Field3: []string{"test"}}

		assert.NotEmpty(t, constraint(validation.NewContext(ts1)))
		assert.NotEmpty(t, constraint(validation.NewContext(ts2)))
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := AtLeastNRequired(1, "Field1", "Field2")(validation.NewContext((*testSubject)(nil)))
		assert.Empty(t, violations)
	})

	t.Run("should return the fields that at least n of should be set in the violation details", func(t *testing.T) {
		ts := testSubject{Field1: "hello"}

		violations := constraint(validation.NewContext(ts))

		require.Len(t, violations, 1)
		assert.Equal(t, map[string]any{
			"actual":  1,
			"minimum": 2,
			"fields":  []string{"Field1", "Field2", "field3", "field4"},
		}, violations[0].Details)
	})

	t.Run("should return violations if given a value of the wrong type, and the value is not empty", func(t *testing.T) {
		ctx1 := validation.NewContext("hello")
		ctx2 := validation.NewContext(123)
		ctx3 := validation.NewContext(url.Values{"not": []string{"empty"}})

		assert.Len(t, constraint(ctx1), 1)
		assert.Len(t, constraint(ctx2), 1)
		assert.Len(t, constraint(ctx3), 1)
	})

	t.Run("should panic if the value of n is 0 or less", func(t *testing.T) {
		assert.Panics(t, func() { AtLeastNRequired(0, "test")(validation.NewContext(testSubject{})) })
		assert.Panics(t, func() { AtLeastNRequired(-10, "test")(validation.NewContext(testSubject{})) })
		assert.Panics(t, func() { AtLeastNRequired(-99999, "test")(validation.NewContext(testSubject{})) })
	})

	t.Run("should panic if number of fields passed to constraint doesn't exceed n", func(t *testing.T) {
		assert.Panics(t, func() { AtLeastNRequired(2, "test")(validation.NewContext(testSubject{})) })
	})
}
