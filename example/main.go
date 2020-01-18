package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/go-validation"
)

func main() {
	example := Example{}
	example.Bool = true
	example.Text = "Hello, GitHub!"
	example.TextMap = map[string]string{"hello longer key": "world"}
	example.Int = 999
	example.Int2 = &example.Int
	example.Ints = []int{1}
	example.Float = math.Pi
	example.Nested = &NestedExample{Text: "Hello, GitHub!"}
	example.Adults = 2
	example.Children = 4
	example.Times = []time.Time{
		time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	example.Nested2 = &Example{}

	violations := validation.Validate(example, exampleConstraints(example))

	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.Encode(violations)

	fmt.Println(buf.String())

	protoViolations := validation.ConstraintViolationsToProto(violations)

	spew.Dump(protoViolations)
}
