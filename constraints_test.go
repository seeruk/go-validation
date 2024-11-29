package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstraints(t *testing.T) {
	t.Run("should run all constraints", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		Validate((*TestSubject)(nil), Constraints{
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		})

		assert.Equal(t, 4, testConstraint.Calls)
	})

	t.Run("should return all constraint violations", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		violations := Validate((*TestSubject)(nil), Constraints{
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		})

		assert.Len(t, violations, 4)
	})

	t.Run("should sort the returned violations by path", func(t *testing.T) {
		subject := TestSubject{
			Text:   "abc",
			Number: 123,
		}

		for i := 0; i < 10000; i++ {
			violations := Validate(subject, Fields{
				"Text":   &TestConstraint{},
				"Number": &TestConstraint{},
			})

			require.Len(t, violations, 2)
			require.Equal(t, ".number", violations[0].Path)
			require.Equal(t, ".text", violations[1].Path)
		}
	})
}

func TestElements(t *testing.T) {
	t.Run("should run all constraints", func(t *testing.T) {
		t.Run("against an array of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := [6]int{1, 2, 3, 4, 5, 6}

			Validate(values, Elements{
				testConstraint,
			})

			assert.Equal(t, 6, testConstraint.Calls)
		})

		t.Run("against a map of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := map[int]int{1: 2, 3: 4, 5: 6, 7: 8, 9: 10, 11: 12}

			Validate(values, Elements{
				testConstraint,
			})

			assert.Equal(t, 6, testConstraint.Calls)
		})

		t.Run("against a slice of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := []int{1, 2, 3, 4, 5, 6}

			Validate(values, Elements{
				testConstraint,
			})

			assert.Equal(t, 6, testConstraint.Calls)
		})
	})

	t.Run("should return all constraint violations", func(t *testing.T) {
		t.Run("against an array of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := [6]int{1, 2, 3, 4, 5, 6}

			violations := Validate(values, Elements{
				testConstraint,
			})

			assert.Len(t, violations, 6)
		})

		t.Run("against a map of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := map[int]int{1: 2, 3: 4, 5: 6, 7: 8, 9: 10, 11: 12}

			violations := Validate(values, Elements{
				testConstraint,
			})

			assert.Len(t, violations, 6)
		})

		t.Run("against a slice of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := []int{1, 2, 3, 4, 5, 6}

			violations := Validate(values, Elements{
				testConstraint,
			})

			assert.Len(t, violations, 6)
		})
	})

	t.Run("should return no violations if the given value is nil", func(t *testing.T) {
		// NOTE: An array cannot be nil, and must have the length specified by it's type.

		t.Run("against a map", func(t *testing.T) {
			var m map[string]any

			testConstraint := &TestConstraint{}
			violations := Validate(m, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)
		})

		t.Run("against a slice", func(t *testing.T) {
			var s []string

			testConstraint := &TestConstraint{}
			violations := Validate(s, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)
		})
	})

	t.Run("should return no violations if the given value's length is 0", func(t *testing.T) {
		// NOTE: An array cannot be nil, and must have the length specified by it's type.

		t.Run("against a map", func(t *testing.T) {
			m := make(map[string]any, 0)

			testConstraint := &TestConstraint{}
			violations := Validate(m, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)
		})

		t.Run("against a slice", func(t *testing.T) {
			s := make([]string, 0)

			testConstraint := &TestConstraint{}
			violations := Validate(s, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)
		})
	})

	t.Run("should update the context's value node to the elements of the given value", func(t *testing.T) {
		t.Run("against an array", func(t *testing.T) {
			a := [2]string{"Hello", "World"}

			var values []string

			Validate(a, Elements{
				ConstraintFunc(func(ctx Context) []ConstraintViolation {
					values = append(values, ctx.Value().Node.Interface().(string))
					return nil
				}),
			})

			require.Len(t, values, len(a))
			for i := range a {
				assert.Equal(t, a[i], values[i])
			}
		})

		t.Run("against a map", func(t *testing.T) {
			m := map[string]string{"Hello": "World"}

			var value string
			Validate(m, Elements{
				ConstraintFunc(func(ctx Context) []ConstraintViolation {
					value = ctx.Value().Node.Interface().(string)
					return nil
				}),
			})

			assert.Equal(t, m["Hello"], value)
		})

		t.Run("against a slice", func(t *testing.T) {
			s := []string{"Hello", "World"}

			var values []string

			Validate(s, Elements{
				ConstraintFunc(func(ctx Context) []ConstraintViolation {
					values = append(values, ctx.Value().Node.Interface().(string))
					return nil
				}),
			})

			require.Len(t, values, len(s))
			for i := range s {
				assert.Equal(t, s[i], values[i])
			}
		})
	})

	t.Run("should update the path", func(t *testing.T) {
		t.Run("against an array", func(t *testing.T) {
			a := [2]string{"Hello", "World"}

			testConstraint := &TestConstraint{}
			violations := Validate(a, Elements{
				testConstraint,
			})

			require.Len(t, violations, 2)
			assert.Equal(t, ".[0]", violations[0].Path)
			assert.Equal(t, ".[1]", violations[1].Path)
		})

		t.Run("against a map", func(t *testing.T) {
			m := map[string]string{
				"Hello": "World",
			}

			testConstraint := &TestConstraint{}
			violations := Validate(m, Elements{
				testConstraint,
			})

			require.Len(t, violations, 1)
			assert.Equal(t, ".Hello", violations[0].Path)
		})

		t.Run("against a map with a nil key", func(t *testing.T) {
			m := map[*string]string{
				nil: "test",
			}

			testConstraint := &TestConstraint{}
			violations := Validate(m, Elements{
				testConstraint,
			})

			require.Len(t, violations, 1)
			assert.Equal(t, ".nil", violations[0].Path)
		})

		t.Run("against a slice", func(t *testing.T) {
			a := []string{"Hello", "World"}

			testConstraint := &TestConstraint{}
			violations := Validate(a, Elements{
				testConstraint,
			})

			require.Len(t, violations, 2)
			assert.Equal(t, ".[0]", violations[0].Path)
			assert.Equal(t, ".[1]", violations[1].Path)
		})
	})

	t.Run("should return violations if the given type is not allowed, and the value is not empty", func(t *testing.T) {
		vctx := NewContext("hello world")

		testConstraint := &TestConstraint{}

		violations := ValidateContext(vctx, Elements{
			testConstraint,
		})

		assert.Len(t, violations, 1)
	})
}

func TestFields(t *testing.T) {
	type fieldsTester struct {
		Foo string
		Bar int
		Baz []string
	}

	t.Run("should run all constraints", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		Validate(fieldsTester{}, Fields{
			"Foo": testConstraint,
			"Bar": testConstraint,
			"Baz": testConstraint,
		})

		assert.Equal(t, 3, testConstraint.Calls)
	})

	t.Run("should return all constraint violations", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate(fieldsTester{}, Fields{
			"Foo": testConstraint,
			"Bar": testConstraint,
			"Baz": testConstraint,
		})

		assert.Len(t, violations, 3)
	})

	t.Run("should return no violations if the given value is nil", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate((*fieldsTester)(nil), Fields{
			"Foo": testConstraint,
			"Bar": testConstraint,
			"Baz": testConstraint,
		})

		assert.Len(t, violations, 0)
	})

	t.Run("should update the context's value node to the fields of the given value", func(t *testing.T) {
		fieldsTester := fieldsTester{
			Foo: "this is a test",
			Bar: 123,
			Baz: []string{"Hello, Go!"},
		}

		var foo string
		var bar int
		var baz []string

		violations := Validate(fieldsTester, Fields{
			"Foo": ConstraintFunc(func(ctx Context) []ConstraintViolation {
				foo = ctx.Value().Node.Interface().(string)
				return nil
			}),
			"Bar": ConstraintFunc(func(ctx Context) []ConstraintViolation {
				bar = ctx.Value().Node.Interface().(int)
				return nil
			}),
			"Baz": ConstraintFunc(func(ctx Context) []ConstraintViolation {
				baz = ctx.Value().Node.Interface().([]string)
				return nil
			}),
		})

		require.Len(t, violations, 0)
		assert.Equal(t, fieldsTester.Foo, foo)
		assert.Equal(t, fieldsTester.Bar, bar)
		assert.Equal(t, fieldsTester.Baz, baz)
	})

	t.Run("should update the path", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		fieldsTester := fieldsTester{
			Foo: "will fail",
		}

		violations := Validate(fieldsTester, Fields{
			"Foo": testConstraint,
		})

		require.Len(t, violations, 1)
		assert.Equal(t, ".Foo", violations[0].Path)
	})

	t.Run("should return violations if the given type is not allowed, and the value is not empty", func(t *testing.T) {
		vctx := NewContext("hello world")

		testConstraint := &TestConstraint{}

		violations := ValidateContext(vctx, Fields{
			"Foo": testConstraint,
		})

		assert.Len(t, violations, 1)
	})
}

func TestKeys(t *testing.T) {
	t.Run("should run all constraints", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		mapTester := map[string]any{
			"Foo": "Hello",
			"Bar": 123,
			"Baz": []string{"Hello", "World"},
		}

		Validate(mapTester, Keys{
			testConstraint,
		})

		assert.Equal(t, 3, testConstraint.Calls)
	})

	t.Run("should return all constraint violations", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		mapTester := map[string]any{
			"Foo": "Hello",
			"Bar": 123,
			"Baz": []string{"Hello", "World"},
		}

		violations := Validate(mapTester, Keys{
			testConstraint,
		})

		assert.Len(t, violations, 3)
	})

	t.Run("should return no violations if the given value is nil", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate((map[string]any)(nil), Keys{
			testConstraint,
			testConstraint,
			testConstraint,
		})

		assert.Len(t, violations, 0)

		// Pointers to map are also possible.
		violations = Validate((*map[string]any)(nil), Keys{
			testConstraint,
			testConstraint,
			testConstraint,
		})

		assert.Len(t, violations, 0)
	})

	t.Run("should return no violations if the given map is empty (but not nil)", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate(map[string]any{}, Keys{
			testConstraint,
			testConstraint,
			testConstraint,
		})

		assert.Len(t, violations, 0)
	})

	t.Run("should update the context's value node to the fields of the given value", func(t *testing.T) {
		mapTester := map[string]any{
			"Foo": "Hello",
			"Bar": 123,
			"Baz": []string{"Hello", "World"},
		}

		var keys []string

		violations := Validate(mapTester, Keys{
			ConstraintFunc(func(ctx Context) []ConstraintViolation {
				keys = append(keys, ctx.Value().Node.Interface().(string))
				return nil
			}),
		})

		require.Len(t, violations, 0)
		require.Len(t, keys, 3)
		assert.Contains(t, keys, "Foo")
		assert.Contains(t, keys, "Bar")
		assert.Contains(t, keys, "Baz")
	})

	t.Run("should update the path", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		mapTester := map[string]any{
			"Foo": "Hello",
		}

		violations := Validate(mapTester, Keys{
			testConstraint,
		})

		require.Len(t, violations, 1)
		assert.Equal(t, ".Foo", violations[0].Path)
	})

	t.Run("should update the path, even with a nil key", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		mapTester := map[*string]any{
			nil: "Hello",
		}

		violations := Validate(mapTester, Keys{
			testConstraint,
		})

		require.Len(t, violations, 1)
		assert.Equal(t, ".nil", violations[0].Path)
	})

	t.Run("should return violations if the given type is not allowed, and the value is not empty", func(t *testing.T) {
		vctx := NewContext("hello world")

		testConstraint := &TestConstraint{}

		violations := ValidateContext(vctx, Keys{
			testConstraint,
		})

		assert.Len(t, violations, 1)
	})
}

func TestLazy(t *testing.T) {
	t.Run("should not execute the constraint(s) upon construction", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		_ = Lazy(func() Constraint { return testConstraint })

		assert.Equal(t, 0, testConstraint.Calls)
	})

	t.Run("should execute the constraint(s) during validation", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate((*TestSubject)(nil), Lazy(func() Constraint { return testConstraint }))

		assert.Equal(t, 1, testConstraint.Calls)
		assert.Len(t, violations, 1)
	})

	t.Run("should support non-struct values", func(t *testing.T) {
		testConstraint := &TestConstraint{
			NoViolation: true,
		}

		violations := Validate("hello world", Lazy(func() Constraint { return testConstraint }))

		assert.Equal(t, 1, testConstraint.Calls)
		assert.Len(t, violations, 0)
	})
}

func TestLazyDynamic(t *testing.T) {
	t.Run("should not execute the constraint(s) upon construction", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		_ = LazyDynamic(func(i int) Constraint { return testConstraint })

		assert.Equal(t, 0, testConstraint.Calls)
	})

	t.Run("should execute the constraint(s) during validation", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate(&TestSubject{}, LazyDynamic(func(ts TestSubject) Constraint { return testConstraint }))

		assert.Equal(t, 1, testConstraint.Calls)
		assert.Len(t, violations, 1)
	})

	t.Run("should not execute the constraint(s) if the value is nil", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate((*TestSubject)(nil), LazyDynamic(func(ts TestSubject) Constraint { return testConstraint }))

		assert.Equal(t, 0, testConstraint.Calls)
		assert.Empty(t, violations)
	})

	t.Run("should panic if the given lazy function does not have exactly 1 argument", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		assert.Panics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func() Constraint { return testConstraint }))
		})

		assert.Panics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(a1, a2 int) Constraint { return testConstraint }))
		})

		assert.Panics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(a1, a2, a3 int, a4 string) Constraint { return testConstraint }))
		})
	})

	t.Run("should panic if the given lazy function doesn't return a Constraint", func(t *testing.T) {
		assert.Panics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(ts TestSubject) {}))
		})

		assert.Panics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(ts TestSubject) int { return 1 }))
		})
	})

	t.Run("should panic if the given lazy function doesn't accept either the original value, or the unwrapped value's type", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		assert.Panics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(t time.Time) Constraint { return testConstraint }))
		})

		assert.Panics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(t **TestSubject) Constraint { return testConstraint }))
		})
	})

	t.Run("should not panic if the given lazy function accepts the original value's type", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		assert.NotPanics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(t *TestSubject) Constraint { return testConstraint }))
		})
	})

	t.Run("should not panic if the given lazy function accepts the unwrapped value's type", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		assert.NotPanics(t, func() {
			Validate(&TestSubject{}, LazyDynamic(func(t TestSubject) Constraint { return testConstraint }))
		})
	})

	t.Run("should return violations if the given type is not allowed, and the value is not empty", func(t *testing.T) {
		vctx := NewContext("hello world")

		testConstraint := &TestConstraint{}

		violations := ValidateContext(vctx, LazyDynamic(func(t TestSubject) Constraint { return testConstraint }))

		assert.Len(t, violations, 1)
	})
}

func TestMap(t *testing.T) {
	t.Run("should run all constraints", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		mapTester := map[string]any{
			"Foo": "Hello",
			"Bar": 123,
			"Baz": []string{"Hello", "World"},
		}

		Validate(mapTester, Map{
			"Foo": testConstraint,
			"Bar": testConstraint,
			"Baz": testConstraint,
		})

		assert.Equal(t, 3, testConstraint.Calls)
	})

	t.Run("should return all constraint violations", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		mapTester := map[string]any{
			"Foo": "Hello",
			"Bar": 123,
			"Baz": []string{"Hello", "World"},
		}

		violations := Validate(mapTester, Map{
			"Foo": testConstraint,
			"Bar": testConstraint,
			"Baz": testConstraint,
		})

		assert.Len(t, violations, 3)
	})

	t.Run("should return no violations if the given value is nil", func(t *testing.T) {
		testConstraint := &TestConstraint{}

		violations := Validate((map[string]any)(nil), Map{
			"Foo": testConstraint,
			"Bar": testConstraint,
			"Baz": testConstraint,
		})

		assert.Len(t, violations, 0)

		// Pointers to map are also possible.
		violations = Validate((*map[string]any)(nil), Map{
			"Foo": testConstraint,
			"Bar": testConstraint,
			"Baz": testConstraint,
		})

		assert.Len(t, violations, 0)
	})

	t.Run("should update the context's value node to the fields of the given value", func(t *testing.T) {
		mapTester := map[string]interface{}{
			"Foo": "Hello",
			"Bar": 123,
			"Baz": []string{"Hello", "World"},
		}

		var foo string
		var bar int
		var baz []string

		violations := Validate(mapTester, Map{
			"Foo": ConstraintFunc(func(ctx Context) []ConstraintViolation {
				foo = ctx.Value().Node.Interface().(string)
				return nil
			}),
			"Bar": ConstraintFunc(func(ctx Context) []ConstraintViolation {
				bar = ctx.Value().Node.Interface().(int)
				return nil
			}),
			"Baz": ConstraintFunc(func(ctx Context) []ConstraintViolation {
				baz = ctx.Value().Node.Interface().([]string)
				return nil
			}),
		})

		require.Len(t, violations, 0)
		assert.Equal(t, mapTester["Foo"], foo)
		assert.Equal(t, mapTester["Bar"], bar)
		assert.Equal(t, mapTester["Baz"], baz)
	})

	t.Run("should update the path", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		mapTester := map[string]interface{}{
			"Foo": "Hello",
		}

		violations := Validate(mapTester, Map{
			"Foo": testConstraint,
		})

		require.Len(t, violations, 1)
		assert.Equal(t, ".Foo", violations[0].Path)
	})

	t.Run("should return violations if the given type is not allowed, and the value is not empty", func(t *testing.T) {
		vctx := NewContext("hello world")

		testConstraint := &TestConstraint{}

		violations := ValidateContext(vctx, Map{
			"Foo": testConstraint,
		})

		assert.Len(t, violations, 1)
	})

	t.Run("should run constraints on empty fields", func(t *testing.T) {
		constraint := Map{
			"foo": &TestConstraint{},
		}

		violations := Validate(map[string]interface{}{}, constraint)

		assert.Len(t, violations, 1)
	})
}

func TestWhen(t *testing.T) {
	t.Run("should run all constraints when the predicate is true", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		Validate((*TestSubject)(nil), When(true,
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Equal(t, 4, testConstraint.Calls)
	})

	t.Run("should return all constraint violations when the predicate is true", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		violations := Validate((*TestSubject)(nil), When(true,
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Len(t, violations, 4)
	})

	t.Run("should not run any constraints when the predicate is false", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		Validate((*TestSubject)(nil), When(false,
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Equal(t, 0, testConstraint.Calls)
	})

	t.Run("should not return any constraint violations when the predicate is false", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		violations := Validate((*TestSubject)(nil), When(false,
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Len(t, violations, 0)
	})
}

func TestWhenFn(t *testing.T) {
	t.Run("should run all constraints when the predicate function returns true", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		Validate((*TestSubject)(nil), WhenFn(func(ctx Context) bool { return true },
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Equal(t, 4, testConstraint.Calls)
	})

	t.Run("should return all constraint violations when the predicate function returns true", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		violations := Validate((*TestSubject)(nil), WhenFn(func(ctx Context) bool { return true },
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Len(t, violations, 4)
	})

	t.Run("should not run any constraints when the predicate function returns false", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		Validate((*TestSubject)(nil), WhenFn(func(ctx Context) bool { return false },
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Equal(t, 0, testConstraint.Calls)
	})

	t.Run("should not return any constraint violations when the predicate function returns false", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		violations := Validate((*TestSubject)(nil), WhenFn(func(ctx Context) bool { return false },
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		))

		assert.Len(t, violations, 0)
	})
}

type TestSubject struct {
	Text   string `validation:"text"`
	Number int    `validation:"number"`
}

type TestConstraint struct {
	Calls       int
	NoViolation bool
}

func (c *TestConstraint) Violations(ctx Context) []ConstraintViolation {
	c.Calls++

	var violations []ConstraintViolation
	if !c.NoViolation {
		violations = append(violations, ctx.Violation("test violations", nil))
	}

	return violations
}
