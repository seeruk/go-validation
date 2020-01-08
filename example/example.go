package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/go-validation"
	"github.com/seeruk/go-validation/constraints"
)

func main() {
	e := Example{}
	e.Text = "Hello, World"
	e.Number = 123

	spew.Dump(validation.Validate(0, constraints.Required))
	spew.Dump(validation.Validate(e, exampleConstraints()))
}

// Example ...
type Example struct {
	Text   string                 `json:"text"`
	Number int                    `json:"number"`
	Object map[string]interface{} `json:"object"`
}

// exampleConstraints ...
func exampleConstraints() validation.Constraint {
	// NOTE: The Example value doesn't need to be passed in, it can just be used to build more
	// dynamic constraints. If you don't need the value, you could actually run this function once
	// which would probably be more efficient - not sure how much by.

	// This approach leaves open quite a flexible approach to building up validation. You wouldn't
	// have to return Fields for example, and you could pass any value to be validated, any any set
	// of constraints really.

	// This approach is extremely similar to Phil's. Implementing it could be quite tricky.
	// Realistically, if we avoid traversing paths, it's probably not so bad. We still need to build
	// up paths to the current thing being validated though. Constraints would need to accept some
	// kind of context that includes the current thing being validated. A constraint that calls
	// child constraints would choose what context it provides to it's children (e.g. fields should
	// get a struct, and pass each field to each of the child constraints.

	structConstraints := validation.Constraints{
		constraints.MutuallyExclusive("Text", "Number"),
	}

	fieldConstraints := validation.Fields{
		"Text": validation.Constraints{
			constraints.Required,
		},
		"Object": validation.Constraints{
			constraints.Required,
			validation.Elements{
				constraints.Required,
			},
			validation.Keys{
				constraints.Required,
			},
		},
	}

	return validation.Constraints{
		structConstraints,
		fieldConstraints,
	}
}
