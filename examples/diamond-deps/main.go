package main

import (
	"fmt"

	"github.com/avoronkov/go-target-once/lib/builder"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

func main() {
	w := warehouse.NewMemoryWarehouse()

	func() {
		a := new(A)
		b := builder.NewBuildSession(w)

		cont, tm, err := b.Build(a)

		fmt.Printf("A -> %v\n", cont)
		fmt.Printf("time: %v\n", tm)
		fmt.Printf("err: %v\n", err)
	}()

	func() {
		a := new(A)
		b := builder.NewBuildSession(w)

		cont, tm, err := b.Build(a)

		fmt.Printf("A -> %v\n", cont)
		fmt.Printf("time: %v\n", tm)
		fmt.Printf("err: %v\n", err)
	}()
}
