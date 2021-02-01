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

		res := b.Build(a)

		fmt.Printf("A -> %v\n", res.Content)
		fmt.Printf("time: %v\n", res.Time)
		fmt.Printf("err: %v\n", res.Err)
	}()

	func() {
		a := new(A)
		b := builder.NewBuildSession(w)

		res := b.Build(a)

		fmt.Printf("A -> %v\n", res.Content)
		fmt.Printf("time: %v\n", res.Time)
		fmt.Printf("err: %v\n", res.Err)
	}()
}
