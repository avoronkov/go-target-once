package main

import (
	"fmt"

	"github.com/avoronkov/go-target-once/lib/builder"
)

func main() {
	a := new(A)
	b := builder.NewBuildSession()

	cont, tm, err := b.Build(a)

	fmt.Printf("A -> %v\n", cont)
	fmt.Printf("time: %v\n", tm)
	fmt.Printf("err: %v\n", err)
}
