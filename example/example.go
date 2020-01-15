package main

import (
	"math"
	"regexp"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/seeruk/go-validation/constraints"
)

// patternGreeting is a regular expression to test that a string starts with "Hello".
var patternGreeting = regexp.MustCompile("^Hello")

// timeYosemite is a time that represents when Yosemite National Park was founded.
var timeYosemite = time.Date(1890, time.October, 1, 0, 0, 0, 0, time.UTC)

// Example ...
type Example struct {
	Bool      bool                       `json:"bool,omitempty"`
	Chan      <-chan string              `json:"chan" validation:"chan"`
	Text      string                     `json:"text"`
	Texts     []string                   `json:"texts" validation:"texts"`
	TextMap   map[string]string          `json:"text_map"`
	Adults    int                        `json:"adults"`
	Children  int                        `json:"children" validation:"children"`
	Int       int                        `json:"int"`
	Int2      *int                       `json:"int2" validation:"int2"`
	Ints      []int                      `json:"ints"`
	Float     float64                    `json:"float" validation:"float"`
	Time      time.Time                  `json:"time" validation:"time"`
	Times     []time.Time                `json:"times"`
	Nested    *NestedExample             `json:"nested" validation:"nested"`
	Nesteds   []*NestedExample           `json:"nesteds"`
	NestedMap map[NestedExample]struct{} `json:"nested_map" validation:"nested_map"`
}

func exampleConstraints(e Example) validation.Constraint {
	return validation.Constraints{
		// Struct constraints ...
		constraints.MutuallyExclusive("Text", "Texts"),
		constraints.MutuallyInclusive("Int", "Int2", "Ints"),
		constraints.AtLeastNRequired(3, "Text", "Int", "Int2", "Ints"),

		validation.Fields{
			"Bool": validation.Constraints{
				constraints.NotEquals(false),
				constraints.Equals(true),
			},
			"Chan": constraints.MaxLength(12),
			"Text": validation.Constraints{
				constraints.Required,
				constraints.Regexp(patternGreeting),
				constraints.MaxLength(14),
				constraints.Length(14),
				constraints.OneOf("Hello, World!", "Hello, SeerUK!", "Hello, GitHub!"),
			},
			"TextMap": validation.Constraints{
				constraints.Required,
				validation.Elements{
					constraints.Required,
				},
				validation.Keys{
					constraints.MinLength(10),
				},
			},
			"Int": constraints.Required,
			"Int2": validation.Constraints{
				constraints.Required,
				constraints.NotNil,
				constraints.Min(0),
			},
			"Ints": validation.Constraints{
				constraints.Required,
				constraints.MaxLength(3),
				validation.Elements{
					constraints.Required,
					constraints.Min(0),
				},
			},
			"Float": constraints.Equals(math.Pi),
			"Time":  constraints.TimeBefore(timeYosemite),
			"Times": validation.Constraints{
				constraints.MinLength(1),
				validation.Elements{
					constraints.TimeBefore(timeYosemite),
				},
			},
			"Adults": validation.Constraints{
				constraints.Min(1),
				constraints.Max(9),
			},
			"Children": validation.Constraints{
				constraints.Min(0),
				constraints.Equals(e.Adults + 2),
				constraints.Max(math.Max(float64(8-(e.Adults-1)), 0)),
			},
			"Nested": validation.Constraints{
				constraints.Required,
				nestedExampleConstraints(),
			},
			"Nesteds": validation.Elements{
				nestedExampleConstraints(),
			},
			"NestedMap": validation.Keys{
				nestedExampleConstraints(),
			},
		},

		validation.When(
			len(e.Text) > 32,
			validation.Constraints{
				constraints.Required,
				constraints.MinLength(64),
			},
		),
	}
}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}

func nestedExampleConstraints() validation.Constraint {
	return validation.Fields{
		"Text": constraints.Required,
	}
}
