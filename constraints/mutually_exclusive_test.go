package constraints

import (
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestMutuallyExclusive(t *testing.T) {
	type testSubject struct {
		Field1 string
		Field2 int
		Field3 []string       `validation:"field3"`
		Field4 map[string]int `validation:"field4"`
	}

	constraint := MutuallyExclusive("Field1", "Field2", "Field3", "Field4")

	t.Run("should return no violations for a valid value", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello"}
		ts2 := testSubject{Field2: 1234567}
		ts3 := testSubject{Field3: []string{"test"}}
		ts4 := testSubject{Field4: map[string]int{"test": 123}}

		assert.Empty(t, constraint(validation.NewContext(ts1)))
		assert.Empty(t, constraint(validation.NewContext(ts2)))
		assert.Empty(t, constraint(validation.NewContext(ts3)))
		assert.Empty(t, constraint(validation.NewContext(ts4)))
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := MutuallyExclusive("Field1", "Field2")(validation.NewContext((*testSubject)(nil)))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if multiple mx fields are set", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello", Field2: 1234567}
		ts2 := testSubject{Field3: []string{"test"}, Field4: map[string]int{"test": 123}}

		assert.NotEmpty(t, constraint(validation.NewContext(ts1)))
		assert.NotEmpty(t, constraint(validation.NewContext(ts2)))
	})

	t.Run("should return the fields that were not empty in the violation details", func(t *testing.T) {
		ts := testSubject{Field1: "hello", Field2: 1234567}

		violations := constraint(validation.NewContext(ts))

		assert.Equal(t, map[string]interface{}{
			"fields": []string{"Field1", "Field2"},
		}, violations[0].Details)
	})

	t.Run("should return the field aliases if set in the violation details", func(t *testing.T) {
		ts := testSubject{Field3: []string{"test"}, Field4: map[string]int{"test": 123}}

		violations := constraint(validation.NewContext(ts))

		assert.Equal(t, map[string]interface{}{
			"fields": []string{"field3", "field4"},
		}, violations[0].Details)
	})

	t.Run("should return no violations if the value is nil", func(t *testing.T) {
		var ts *testSubject
		assert.Empty(t, constraint(validation.NewContext(ts)))
	})
}
