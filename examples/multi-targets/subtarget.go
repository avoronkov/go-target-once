package main

import (
	"fmt"
	"log"
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type SubTgt struct {
	value int
}

var _ targets.Target = (*SubTgt)(nil)

func NewSubTgt(value int) *SubTgt {
	return &SubTgt{value: value}
}

func (s *SubTgt) Build(bc targets.BuildContext) (content interface{}, t time.Time, err error) {
	log.Printf("[%v] Build()", s.TargetID())
	time.Sleep(1 * time.Second)
	return s.value * s.value, time.Now(), nil
}

func (s *SubTgt) Dependencies() []targets.Target {
	return nil
}

func (s *SubTgt) IsModified(since time.Time) (m bool) {
	defer log.Printf("[%v] IsModified(): %v", s.TargetID(), m)
	return s.value > 5
}

func (s *SubTgt) TargetID() string {
	return fmt.Sprintf("sub-target-%d", s.value)
}

func (s *SubTgt) Cacheable() bool {
	return true
}
