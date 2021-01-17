package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/avoronkov/go-target-once/lib/builder"
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
