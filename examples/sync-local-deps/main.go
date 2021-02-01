package main

import (
	"fmt"

	"github.com/avoronkov/go-target-once/lib/builder"
)

func main() {
	b := new(B)
	c := new(C)

	results := builder.Builds(b, c)

	fmt.Printf("B -> %v\n", results[0].Content)
	fmt.Printf("time: %v\n", results[0].Time)
	fmt.Printf("err: %v\n", results[0].Err)

	fmt.Printf("C -> %v\n", results[1].Content)
	fmt.Printf("time: %v\n", results[1].Time)
	fmt.Printf("err: %v\n", results[1].Err)
}
