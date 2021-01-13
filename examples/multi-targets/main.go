package main

import (
	"dont-repeat-twice/lib/builder"
	"fmt"
	"log"
)

func main() {
	t := NewTgt(1, 2, 3, 5, 8, 13)

	res, err := builder.Build(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)

	res, err = builder.Build(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)
}
