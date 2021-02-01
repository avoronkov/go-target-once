package builder

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
	"github.com/avoronkov/go-target-once/lib/targets"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

type BuildSession struct {
	targetResults map[string]*ObservableResult

	// Target ID -> bool
	modifiedTargets sync.Map

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

	// Check cached targets
	// TODO : handle KeepAlive targets
	// TODO : handle locally cached targets
T:
	for id, meta := range tgts {
		// check local cache
		if _, ok := bc.targetResults[id]; ok {
			// no need to build the target
			bc.removeTargetWithDeps(id, &tgts)
			continue T
		}

		if ct, ok := meta.t.(targets.Cachable); ok && ct.Cachable() {
			// check global cache
			cont, tm, ok := bc.globalCache.Get(id)
			if !ok {
				continue T
			}
			// check all subtargets are not modified
			if bc.targetOrDepsModified(id, tm, tgts, &bc.modifiedTargets) {
				continue T
			}

			// If OK then put content into targetResults
			bc.targetResults[id] = NewObservable()
			bc.targetResults[id].Put(&buildResult{
				C: cont,
				T: tm,
			})
			// and remove target and subtargets from tgts
			bc.removeTargetWithDeps(id, &tgts)
		}
	}

	// Create observable results
	for t := range tgts {
		bc.targetResults[t] = NewObservable()
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

func (bc *BuildSession) targetOrDepsModified(id string, since time.Time, tgts map[string]*targetMeta, modified *sync.Map) bool {
	var wg sync.WaitGroup
	var tModified func(string)

	var mod int32

	tModified = func(tid string) {
		defer wg.Done()

		meta := tgts[tid]
		wg.Add(len(meta.depsNames))
		for _, depId := range meta.depsNames {
			go tModified(depId)
		}

		if md, ok := modified.Load(tid); ok {
			if md.(bool) {
				atomic.StoreInt32(&mod, 1)
			}
			return
		}

		if mf, ok := meta.t.(targets.Modifiable); ok {
			md := mf.IsModified(since)
			modified.Store(tid, md)
			if md {
				atomic.StoreInt32(&mod, 1)
			}
		}
	}

	wg.Add(1)
	go tModified(id)
	wg.Wait()
	return mod > 0
}

type buildResult struct {
	C interface{}
	T time.Time
	E error
}
