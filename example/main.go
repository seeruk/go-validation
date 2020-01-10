package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/seeruk/go-validation"
)

func main() {
	e := Example{}
	e.TextMap = map[string]string{
		"Hello": "World!",
	}

	violations := validation.Validate(e, exampleConstraints())

	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.Encode(violations)

	fmt.Println(buf.String())
}
