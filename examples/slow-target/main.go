package main

import (
	"fmt"

	bl "github.com/avoronkov/go-target-once/lib/builder"
)

func main() {
	t := new(Slow)
	res := bl.Build(t)
	fmt.Printf("Content: %v\n", res.Content)
	fmt.Printf("Error: %v\n", res.Err)
	fmt.Printf("Time: %v\n", res.Time)
}
