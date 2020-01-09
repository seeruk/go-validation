package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/seeruk/go-validation"
	"github.com/seeruk/go-validation/constraints"
)

var foo = "foo"

func main() {
	e := &Example{}
	e.Text = "Hello, World"
	e.Number = 123

	e.Object = map[*string]interface{}{
		&foo: 123,
	}

	list := &[]*string{
		&foo,
		nil,
	}

	list2 := &list

	e.List = &list2
	e.List = nil

	e2 := &Example{}

	e.Foo = &e2

	var buf bytes.Buffer

	ctx := validation.NewContext(e)
	ctx.StructTag = "json"

	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.Encode(validation.ValidateContext(ctx, exampleConstraints()))

	fmt.Println(buf.String())
}

// Example ...
type Example struct {
	Text   string                  `json:"text"`
	Number int                     `json:"number"`
	List   ***[]*string            `json:"list,omitempty"`
	Object map[*string]interface{} `json:"object"`
	Foo    **Example               `json:"foo"`
}

// exampleConstraints ...
func exampleConstraints() validation.Constraint {
	structConstraints := validation.Constraints{
		constraints.Required,
		constraints.MutuallyExclusive("Text", "Number"),
	}

	fieldConstraints := validation.Fields{
		"Text": validation.Constraints{
			constraints.Required,
		},
		"List": validation.Constraints{
			constraints.Required,
			validation.Elements{
				constraints.Required,
			},
		},
		"Object": validation.Constraints{
			constraints.Required,
			validation.Elements{
				constraints.Required,
			},
			validation.Keys{
				constraints.Required,
			},
			validation.Map{
				&foo: constraints.Required,
			},
		},
		//"Foo": validation.Lazy(exampleConstraints),
	}

	return validation.Constraints{
		structConstraints,
		fieldConstraints,
	}
}
