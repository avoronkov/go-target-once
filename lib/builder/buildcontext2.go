package builder

import (
	"fmt"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
	"github.com/avoronkov/go-target-once/lib/targets"
)

type BuildContext2 struct {
}

func (bc *BuildContext2) Build(t targets.Target) (content interface{}, tm time.Time, err error) {
	// targetId -> dependencies ids
	targetDeps := map[string][]string{}
	// targetId -> dep name -> dep id
	targetDepsNames := map[string]map[string]string{}

	bc.fillTargetDeps(t, &targetDeps, &targetDepsNames)

	return nil, time.Time{}, fmt.Errorf("Not implemented yet")
}

func (bc *BuildContext2) fillTargetDeps(t targets.Target, td *map[string][]string, tdn *map[string]map[string]string) {
	tid := t.TargetId()
	if _, ok := (*td)[tid]; ok {
		logger.Debugf("Target `%v` already defined. Skipping.", tid)
		return
	}

	tdeps := []string{}
	if withDeps, ok := t.(targets.WithDependencies); ok {
		namesMp := map[string]string{}
		for name, d := range withDeps.Dependencies() {
			did := d.TargetId()
			namesMp[name] = did

			tdeps = append(tdeps, d.TargetId())
			bc.fillTargetDeps(d, td, tdn)
		}
		(*tdn)[tid] = namesMp
	}
	(*td)[tid] = tdeps
}
