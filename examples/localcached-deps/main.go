package main

import (
	"fmt"

	"github.com/avoronkov/go-target-once/lib/builder"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

func main() {
	a := new(A)

	w := warehouse.NewMemoryWarehouse()
	bs := builder.NewBuildSession(w)

	cont, tm, err := bs.Build(a)

	fmt.Printf("A -> %v\n", cont)
	fmt.Printf("time: %v\n", tm)
	fmt.Printf("err: %v\n", err)
}
