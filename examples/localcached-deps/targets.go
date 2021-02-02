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

func (a *A) TargetID() string {
	return "target-A"
}

func (a *A) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"dep-b": new(B),
	}
}

func (a *A) Build(bc targets.BuildContext) targets.Result {
	log.Printf("A.Build()")
	time.Sleep(1 * time.Second)
	log.Printf("A: GetDependency(dep-b)")
	b := bc.GetDependency("dep-b")
	if b.Err != nil {
		return targets.Failed(b.Err)
	}

	time.Sleep(1 * time.Second)

	log.Printf("A: build C")
	targetC := new(C)
	c := bc.Build(targetC)
	if c.Err != nil {
		return targets.Failed(c.Err)
	}

	return targets.OK(fmt.Sprintf("A{B: %v, C: %v}", b.Content, c.Content))
}

type B struct{}

var _ targets.Target = (*B)(nil)

func (a *B) TargetID() string {
	return "target-B"
}

func (a *B) Build(bc targets.BuildContext) targets.Result {
	log.Printf("B.Build()")
	targetC := new(C)
	c := bc.Build(targetC)
	if c.Err != nil {
		return targets.Failed(c.Err)
	}
	return targets.OK(fmt.Sprintf("B{C: %v}", c.Content))
}

type C struct{}

var _ targets.Target = (*C)(nil)

func (a *C) TargetID() string {
	return "target-C"
}

func (a *C) Build(bc targets.BuildContext) targets.Result {
	log.Printf("C.Build()")
	return targets.OK("{C}")
}
