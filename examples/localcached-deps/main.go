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

	res := bs.Build(a)

	fmt.Printf("A -> %v\n", res.Content)
	fmt.Printf("time: %v\n", res.Time)
	fmt.Printf("err: %v\n", res.Err)
}
