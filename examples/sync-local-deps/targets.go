package main

import (
	"fmt"
	"log"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type A struct{}

var _ targets.Target = (*A)(nil)

func (a *A) TargetID() string {
	return "target-A"
}

func (a *A) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"dep-b": new(B),
		"dep-c": new(C),
	}
}

func (a *A) Build(bc targets.BuildContext) targets.Result {
	b := bc.GetDependency("dep-b")
	if b.Err != nil {
		return targets.Failed(b.Err)
	}

	c := bc.GetDependency("dep-c")
	if c.Err != nil {
		return targets.Failed(c.Err)
	}

	return targets.OK(fmt.Sprintf("A{b: %v, c: %v}", b.Content, c.Content))
}

type B struct{}

var _ targets.Target = (*B)(nil)

func (a *B) TargetID() string {
	return "target-B"
}

func (b *B) Build(bc targets.BuildContext) targets.Result {
	d := new(D)
	cd := bc.Build(d)
	if cd.Err != nil {
		return targets.Failed(cd.Err)
	}
	return targets.OK(fmt.Sprintf("B{d: %v}", cd.Content))
}

type C struct {
}

var _ targets.Target = (*C)(nil)

func (c *C) TargetID() string {
	return "target-C"
}

func (c *C) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"dep-b": new(B),
	}
}

func (c *C) Build(bc targets.BuildContext) targets.Result {
	d := new(D)
	cd := bc.Build(d)
	if cd.Err != nil {
		return targets.Failed(cd.Err)
	}
	return targets.OK(fmt.Sprintf("C{d: %v}", cd.Content))
}

type D struct {
}

var _ targets.Target = (*D)(nil)

func (a *D) TargetID() string {
	return "target-D"
}

func (d *D) Build(bc targets.BuildContext) targets.Result {
	log.Printf("D.Build()")
	return targets.OK("{D}")
}
