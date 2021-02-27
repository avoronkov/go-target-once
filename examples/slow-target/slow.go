package main

import (
	"time"

	tg "github.com/avoronkov/go-target-once/lib/targets"
)

type Slow struct {
	Name  string
	Sleep time.Duration
}

var _ tg.Target = (*Slow)(nil)

func (s *Slow) TargetID() string {
	return "slow-" + s.Name
}

func (s *Slow) Build(bc tg.BuildContext) tg.Result {
	time.Sleep(s.Sleep)
	return tg.OK(s.Name + ": finished.")
}
