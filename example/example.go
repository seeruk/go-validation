package main

import (
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

func exampleConstraints() validation.Constraint {
	return validation.Constraints{
		// Struct constraints ...
		constraints.MutuallyExclusive("Text", "Texts"),
		//constraints.MutuallyInclusive("Int", "Int2", "Ints"),
		//constraints.ExactlyNRequired(3, "Text", "Int", "Int2", "Ints"),

		validation.Fields{
			"Text": validation.Constraints{
				constraints.Required,
				constraints.OneOf("hello", "world", "example"),
			},
			"TextMap": validation.Constraints{
				constraints.Required,
				validation.Elements{
					constraints.Required,
				},
				validation.Keys{
					constraints.Required,
					//constraints.MinLength(10),
				},
			},
			"Nested": validation.Constraints{
				constraints.Required,
				nestedExampleConstraints(),
			},
		},
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
