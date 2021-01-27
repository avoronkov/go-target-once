package builder

import (
	"sync"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
	"github.com/avoronkov/go-target-once/lib/targets"
)

type BuildSession struct {
	targetResults map[string]*ObservableResult
}

func NewBuildSession() *BuildSession {
	return &BuildSession{
		targetResults: make(map[string]*ObservableResult),
	}
}

func (bc *BuildSession) Build(t targets.Target) (content interface{}, tm time.Time, err error) {
	// targetId -> dependencies ids
	targetDeps := map[string][]string{}
	// targetId -> dep name -> dep id
	targetDepsNames := map[string]map[string]string{}

	tgts := map[string]targets.Target{}

	bc.fillTargetDeps(t, &targetDeps, &targetDepsNames, &tgts)

	tid := t.TargetId()

	logger.Debugf("[target=%v] targetDeps = %+v", tid, targetDeps)
	logger.Debugf("[target=%v] targetDeps = %+v", tid, targetDepsNames)
	logger.Debugf("[target=%v] targets (%v) = %+V", tid, len(tgts), tgts)

	// Create observable results
	for t := range targetDeps {
		bc.targetResults[t] = NewObservable()
	}

	// Build
	var wg sync.WaitGroup
	wg.Add(len(tgts))

	for id, tgt := range tgts {
		go func(id string, tg targets.Target) {
			logger.Debugf("Building target '%v'...", id)
			sc := NewSessionContext(id, bc, targetDepsNames[id])
			cont, tm, err := tg.Build(sc)
			br := &buildResult{
				C: cont,
				T: tm,
				E: err,
			}
			bc.targetResults[id].Put(br)
			logger.Debugf("Building target '%v': done.", id)
			wg.Done()
		}(id, tgt)
	}

	wg.Wait()

	br := bc.targetResults[t.TargetId()].Get()

	return br.C, br.T, br.E
}

func (bc *BuildSession) fillTargetDeps(t targets.Target, td *map[string][]string, tdn *map[string]map[string]string, tgts *map[string]targets.Target) {
	tid := t.TargetId()
	if _, ok := (*tgts)[tid]; ok {
		logger.Debugf("Target `%v` already defined. Skipping.", tid)
		return
	}

	(*tgts)[tid] = t

	tdeps := []string{}
	if withDeps, ok := t.(targets.WithDependencies); ok {
		namesMp := map[string]string{}
		for name, d := range withDeps.Dependencies() {
			did := d.TargetId()
			namesMp[name] = did

			tdeps = append(tdeps, d.TargetId())
			bc.fillTargetDeps(d, td, tdn, tgts)
		}
		(*tdn)[tid] = namesMp
	}
	(*td)[tid] = tdeps
}

type buildResult struct {
	C interface{}
	T time.Time
	E error
}
