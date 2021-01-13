package main

import (
	"dont-repeat-twice/lib/targets"
	"fmt"
	"log"
	"strings"
	"time"
)

type Tgt struct {
	values []int
}

var _ targets.Target = (*Tgt)(nil)

func NewTgt(values ...int) *Tgt {
	return &Tgt{values: values}
}

func (g *Tgt) Dependencies() (deps []targets.Target) {
	for _, v := range g.values {
		deps = append(deps, NewSubTgt(v))
	}
	return
}

func (g *Tgt) Build(bc targets.BuildContext) (content interface{}, t time.Time, err error) {
	log.Printf("[%v] Build()", g.TargetId())
	s := new(strings.Builder)
	for i, v := range g.values {
		res, err := bc.GetDependency(i)
		fmt.Fprintf(s, "%v * %v => %v (%v)\n", v, v, res, err)
	}
	return s.String(), time.Now(), nil
}

func (g *Tgt) IsModified(since time.Time) bool {
	log.Printf("[%v] IsModified()", g.TargetId())
	return true
}

func (g *Tgt) TargetId() string {
	return fmt.Sprintf("target-%v", g.values)
}
