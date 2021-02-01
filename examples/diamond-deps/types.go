package main

import (
	"fmt"
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type A struct {
}

func (a *A) TargetId() string {
	return "target-A"
}

func (a *A) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"b": new(B),
		"c": new(C),
	}
}

func (a *A) Build(bc targets.BuildContext) targets.Result {
	fmt.Printf("A.Build()...\n")
	b := bc.GetDependency("b")
	c := bc.GetDependency("c")
	content := fmt.Sprintf("A (%v, %v)", b.Content, c.Content)
	time.Sleep(1 * time.Second)
	fmt.Printf("A.Build(): done.\n")
	return targets.ResultOk(content)
}

type B struct {
}

func (b *B) TargetId() string {
	return "target-B"
}

func (b *B) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"d": new(D),
		"e": new(E),
	}
}

func (b *B) Build(bc targets.BuildContext) targets.Result {
	fmt.Printf("B.Build()...\n")
	d := bc.GetDependency("d")
	e := bc.GetDependency("e")
	content := fmt.Sprintf("B (%v, %v)", d.Content, e.Content)
	time.Sleep(1 * time.Second)
	fmt.Printf("B.Build(): done.\n")
	return targets.ResultOk(content)
}

type C struct {
}

func (c *C) TargetId() string {
	return "target-C"
}

func (c *C) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"e": new(E),
		"f": new(F),
	}
}

func (c *C) Cachable() bool {
	return true
}

func (c *C) Build(bc targets.BuildContext) targets.Result {
	fmt.Printf("C.Build()...\n")
	e := bc.GetDependency("e")
	f := bc.GetDependency("f")
	content := fmt.Sprintf("C (%v, %v)", e.Content, f.Content)
	time.Sleep(1 * time.Second)
	fmt.Printf("C.Build():u done.\n")
	return targets.ResultOk(content)
}

type D struct {
}

func (d *D) TargetId() string {
	return "target-D"
}

func (d *D) Build(bc targets.BuildContext) targets.Result {
	fmt.Printf("D.Build()\n")
	time.Sleep(1 * time.Second)
	return targets.ResultOk("{D}")
}

type E struct {
}

func (e *E) TargetId() string {
	return "target-E"
}

func (e *E) Build(bc targets.BuildContext) targets.Result {
	fmt.Printf("E.Build()\n")
	time.Sleep(1 * time.Second)
	return targets.ResultOk("{E}")
}

type F struct {
}

func (f *F) TargetId() string {
	return "target-F"
}

func (f *F) Build(bc targets.BuildContext) targets.Result {
	fmt.Printf("F.Build()\n")
	time.Sleep(1 * time.Second)
	return targets.ResultOk("{F}")
}
