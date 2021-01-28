package builder

import (
	"sync"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
	"github.com/avoronkov/go-target-once/lib/targets"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

type BuildSession struct {
	targetResults map[string]*ObservableResult

	globalCache warehouse.Warehouse
}

func NewBuildSession(globalCache warehouse.Warehouse) *BuildSession {
	return &BuildSession{
		targetResults: make(map[string]*ObservableResult),
		globalCache:   globalCache,
	}
}

func (bc *BuildSession) Build(t targets.Target) (content interface{}, tm time.Time, err error) {
	tgts := map[string]*targetMeta{}

	bc.fillTargetDeps(t, &tgts)

	tid := t.TargetId()

	logger.Debugf("[target=%v] targets (%v) = %+V", tid, len(tgts), tgts)

	// Create observable results
	for t := range tgts {
		bc.targetResults[t] = NewObservable()
	}

	// Check cached targets
T:
	for id, meta := range tgts {
		if ct, ok := meta.t.(targets.Cachable); ok && ct.Cachable() {
			cont, tm, ok := bc.globalCache.Get(id)
			if !ok {
				continue T
			}
			// check all subtargets are not modified
			// TODO

			// If OK then put content into targetResults
			bc.targetResults[id].Put(&buildResult{
				C: cont,
				T: tm,
			})
			// and remove target and subtargets from tgts
			bc.removeTargetWithDeps(id, &tgts)
		}
	}

	// Build
	var wg sync.WaitGroup
	wg.Add(len(tgts))

	for id, tgt := range tgts {
		go func(id string, meta targetMeta) {
			logger.Debugf("Building target '%v'...", id)
			sc := NewSessionContext(id, bc, meta.depsNames)
			cont, tm, err := meta.t.Build(sc)
			br := &buildResult{
				C: cont,
				T: tm,
				E: err,
			}
			bc.targetResults[id].Put(br)
			logger.Debugf("Building target '%v': done.", id)
			wg.Done()
		}(id, *tgt)
	}

	wg.Wait()

	// cache cachable targets
	for id, tgt := range tgts {
		if ct, ok := tgt.t.(targets.Cachable); ok && ct.Cachable() {
			res := bc.targetResults[id].Get()
			if res.E != nil {
				continue
			}
			bc.globalCache.Put(id, res.C, res.T)
		}
	}

	// return result
	br := bc.targetResults[t.TargetId()].Get()

	return br.C, br.T, br.E
}

func (bc *BuildSession) fillTargetDeps(t targets.Target, tgts *map[string]*targetMeta) {
	tid := t.TargetId()

	if meta, ok := (*tgts)[tid]; ok {
		meta.refCounter++

		logger.Debugf("Target `%v` already defined. Skipping.", tid)
		return
	}

	meta := newTargetMeta()
	meta.t = t

	if withDeps, ok := t.(targets.WithDependencies); ok {
		namesMp := map[string]string{}
		for name, d := range withDeps.Dependencies() {
			did := d.TargetId()
			namesMp[name] = did

			bc.fillTargetDeps(d, tgts)
		}
		meta.depsNames = namesMp
	}
	(*tgts)[tid] = meta
}

func (bc *BuildSession) removeTargetWithDeps(id string, tgts *map[string]*targetMeta) {
	meta := (*tgts)[id]

	for _, depId := range meta.depsNames {
		bc.removeTargetWithDeps(depId, tgts)
	}

	meta.refCounter--
	if meta.refCounter == 0 {
		delete(*tgts, id)
	}
}

type buildResult struct {
	C interface{}
	T time.Time
	E error
}
