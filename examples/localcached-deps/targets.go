package main

import (
	"fmt"
	"log"
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type A struct {
}

var _ targets.Target = (*A)(nil)

func (a *A) TargetId() string {
	return "target-A"
}

func (a *A) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"dep-b": new(B),
	}
}

func (a *A) Build(bc targets.BuildContext) (interface{}, time.Time, error) {
	log.Printf("A.Build()")
	time.Sleep(1 * time.Second)
	log.Printf("A: GetDependency(dep-b)")
	b, err := bc.GetDependency("dep-b")
	if err != nil {
		return nil, time.Time{}, nil
	}

	time.Sleep(1 * time.Second)

	log.Printf("A: build C")
	targetC := new(C)
	c, _, err := bc.Build(targetC)
	if err != nil {
		return nil, time.Time{}, nil
	}

	return fmt.Sprintf("A{B: %v, C: %v}", b, c), time.Now(), nil
}

type B struct {
}

var _ targets.Target = (*B)(nil)

func (a *B) TargetId() string {
	return "target-B"
}

func (a *B) Build(bc targets.BuildContext) (interface{}, time.Time, error) {
	log.Printf("B.Build()")
	targetC := new(C)
	c, _, err := bc.Build(targetC)
	if err != nil {
		return nil, time.Time{}, nil
	}
	return fmt.Sprintf("B{C: %v}", c), time.Now(), nil
}

type C struct {
}

var _ targets.Target = (*C)(nil)

func (a *C) TargetId() string {
	return "target-C"
}

func (a *C) Build(bc targets.BuildContext) (interface{}, time.Time, error) {
	log.Printf("C.Build()")
	return "{C}", time.Now(), nil
}
