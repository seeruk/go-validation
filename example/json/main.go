package main

import (
	"encoding/json"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/go-validation"
	"github.com/seeruk/go-validation/constraints"
)

var raw = []byte(`{
	"name": ["hi", "bar", "baz"]
}`)

func main() {
	var input map[string]interface{}

	err := json.Unmarshal(raw, &input)
	if err != nil {
		panic(err)
	}

	spew.Dump(input)

	cc := validation.Map{
		"name": validation.Constraints{
			constraints.Required,
			constraints.Kind(reflect.String),
			constraints.MinLength(2),
			constraints.MaxLength(128),
		},
	}

	violations := validation.Validate(input, cc)
	spew.Dump(violations)
}
