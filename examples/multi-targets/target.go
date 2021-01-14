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
	deps   map[string]targets.Target
}

var _ targets.Target = (*Tgt)(nil)
var _ targets.WithDependencies = (*Tgt)(nil)

func NewTgt(values ...int) *Tgt {
	deps := make(map[string]targets.Target)
	for _, v := range values {
		deps[fmt.Sprintf("%v", v)] = NewSubTgt(v)
	}
	return &Tgt{
		values: values,
		deps:   deps,
	}
}

func (g *Tgt) Dependencies() map[string]targets.Target {
	return g.deps
}

func (g *Tgt) Build(bc targets.BuildContext) (content interface{}, t time.Time, err error) {
	log.Printf("[%v] Build()", g.TargetId())
	s := new(strings.Builder)
	for k := range g.deps {
		res, err := bc.GetDependency(k)
		fmt.Fprintf(s, "%v * %v => %v (%v)\n", k, k, res, err)
	}
	return s.String(), time.Now(), nil
}

func (g *Tgt) TargetId() string {
	return fmt.Sprintf("target-%v", g.values)
}
