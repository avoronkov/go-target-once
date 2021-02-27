package main

import (
	"fmt"
	"time"

	bl "github.com/avoronkov/go-target-once/lib/builder"
)

func main() {
	t1 := &Slow{"T1", 15 * time.Second}
	t2 := &Slow{"T2", 31 * time.Second}
	results := bl.Builds(t1, t2)

	for _, res := range results {
		fmt.Printf("Content: %v\n", res.Content)
		fmt.Printf("Error: %v\n", res.Err)
		fmt.Printf("Time: %v\n", res.Time)
	}
}
