package main

import (
	"dont-repeat-twice/lib/builder"
	"fmt"
	"log"
	"runtime"
	"time"
)

func main() {
	t := NewTgt(1, 2, 3, 5, 8, 13, 21)

	for {
		res, err := builder.Build(t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(res)
		fmt.Println("goroutines: ", runtime.NumGoroutine())
		time.Sleep(500 * time.Millisecond)
	}
}
