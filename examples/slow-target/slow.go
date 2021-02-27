package main

import (
	"time"

	tg "github.com/avoronkov/go-target-once/lib/targets"
)

type Slow struct {
}

var _ tg.Target = (*Slow)(nil)

func (*Slow) TargetID() string {
	return "slow"
}

func (*Slow) Build(bc tg.BuildContext) tg.Result {
	time.Sleep(31 * time.Second)
	return tg.OK("Finished.")
}
