package constraints

import (
	"net/url"
	"testing"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMutuallyInclusive(t *testing.T) {
	type testSubject struct {
		Field1 string
		Field2 int
		Field3 []string       `validation:"field3"`
		Field4 map[string]int `validation:"field4"`
	}

	constraint := MutuallyInclusive("Field1", "Field2", "Field3", "Field4")

	t.Run("should return no violations if all mutually inclusive fields are set", func(t *testing.T) {
		ts1 := testSubject{
			Field1: "hello",
			Field2: 123,
			Field3: []string{"test"},
			Field4: map[string]int{
				"test": 123,
			},
		}

		assert.Empty(t, constraint(validation.NewContext(ts1)))
	})

	t.Run("should return a violation if all mutually inclusive fields are not set", func(t *testing.T) {
		ts1 := testSubject{Field1: "hello", Field2: 1234567}
		ts2 := testSubject{Field3: []string{"test"}, Field4: map[string]int{"test": 123}}

		assert.NotEmpty(t, constraint(validation.NewContext(ts1)))
		assert.NotEmpty(t, constraint(validation.NewContext(ts2)))
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := MutuallyInclusive("Field1", "Field2")(validation.NewContext((*testSubject)(nil)))
		assert.Empty(t, violations)
	})

	t.Run("should return the fields that are mutually inclusive in the violation details", func(t *testing.T) {
		ts := testSubject{Field1: "hello", Field2: 1234567}

		violations := constraint(validation.NewContext(ts))

		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"fields": []string{"Field1", "Field2", "field3", "field4"},
		}, violations[0].Details)
	})

	t.Run("should return the field aliases if set in the violation details", func(t *testing.T) {
		ts := testSubject{Field1: "hello", Field2: 1234567}

		violations := constraint(validation.NewContext(ts))

		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"fields": []string{"Field1", "Field2", "field3", "field4"},
		}, violations[0].Details)
	})

	t.Run("should return violations if given a value of the wrong type, and the value is not empty", func(t *testing.T) {
		ctx1 := validation.NewContext("hi")
		ctx2 := validation.NewContext(123)
		ctx3 := validation.NewContext(url.Values{"not": []string{"empty"}})

		assert.Len(t, constraint(ctx1), 1)
		assert.Len(t, constraint(ctx2), 1)
		assert.Len(t, constraint(ctx3), 1)
	})
}
