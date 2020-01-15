package constraints

import (
	"net/url"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMutuallyExclusive(t *testing.T) {
	type testSubject struct {
		Field1 string
		Field2 int
		Field3 []string       `validation:"field3"`
		Field4 map[string]int `validation:"field4"`
	}

	constraint := MutuallyExclusive("Field1", "Field2", "Field3", "Field4")

	t.Run("should return no violations if multiple mx fields are not set", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello"}
		ts2 := testSubject{Field2: 1234567}
		ts3 := testSubject{Field3: []string{"test"}}
		ts4 := testSubject{Field4: map[string]int{"test": 123}}

		assert.Empty(t, constraint(validation.NewContext(ts1)))
		assert.Empty(t, constraint(validation.NewContext(ts2)))
		assert.Empty(t, constraint(validation.NewContext(ts3)))
		assert.Empty(t, constraint(validation.NewContext(ts4)))
	})

	t.Run("should return a violation if multiple mx fields are set", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello", Field2: 1234567}
		ts2 := testSubject{Field3: []string{"test"}, Field4: map[string]int{"test": 123}}

		assert.NotEmpty(t, constraint(validation.NewContext(ts1)))
		assert.NotEmpty(t, constraint(validation.NewContext(ts2)))
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := MutuallyExclusive("Field1", "Field2")(validation.NewContext((*testSubject)(nil)))
		assert.Empty(t, violations)
	})

	t.Run("should return the fields that were not empty in the violation details", func(t *testing.T) {
		ts := testSubject{Field1: "hello", Field2: 1234567}

		violations := constraint(validation.NewContext(ts))

		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"fields": []string{"Field1", "Field2"},
		}, violations[0].Details)
	})

	t.Run("should return the field aliases if set in the violation details", func(t *testing.T) {
		ts := testSubject{Field3: []string{"test"}, Field4: map[string]int{"test": 123}}

		violations := constraint(validation.NewContext(ts))

		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"fields": []string{"field3", "field4"},
		}, violations[0].Details)
	})

	t.Run("should panic if given a value of the wrong type, even if it's empty", func(t *testing.T) {
		assert.Panics(t, func() { constraint(validation.NewContext("")) })
		assert.Panics(t, func() { constraint(validation.NewContext(0)) })
		assert.Panics(t, func() { constraint(validation.NewContext(url.Values{})) })
	})
}
