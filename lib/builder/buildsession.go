package builder

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
	"github.com/avoronkov/go-target-once/lib/targets"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

type BuildSession struct {
	// Target ID -> *ObservableResult
	targetResults sync.Map

	// Target ID -> bool
	modifiedTargets sync.Map

	globalCache warehouse.Warehouse

	mutex sync.Mutex

	// Target ID -> struct{}
	builtTargets sync.Map
}

func NewBuildSession(globalCache warehouse.Warehouse) *BuildSession {
	return &BuildSession{
		globalCache: globalCache,
	}
}

func (bc *BuildSession) Build(t targets.Target) targets.Result {
	tgts := map[string]*targetMeta{}

	bc.fillTargetDeps(t, &tgts)

	tid := t.TargetID()

	logger.Debugf("[target=%v] targets (%v) = %+V", tid, len(tgts), tgts)

	// sync access to bc.targetResults
	bc.mutex.Lock()

	// Check cached targets
	// TODO : handle KeepAlive targets
	// TODO : handle locally cached targets
T:
	for id, meta := range tgts {
		if meta.refCounter <= 0 {
			// Target already removed
			continue T
		}

		// check local cache
		if _, ok := bc.targetResults.Load(id); ok {
			// no need to build the target
			bc.removeTargetWithDeps(id, &tgts)
			continue T
		}

		if ct, ok := meta.t.(targets.Cacheable); ok && ct.Cacheable() {
			// check global cache
			cont, tm, ok := bc.globalCache.Get(id)
			if !ok {
				continue T
			}
			if _, keepAlive := meta.t.(targets.KeepingAlive); !keepAlive {
				// check all subtargets are not modified
				if bc.targetOrDepsModified(id, tm, tgts, &bc.modifiedTargets) {
					continue T
				}
			}

			// If OK then put content into targetResults
			obs := NewObservable()
			obs.Put(&targets.Result{
				Content: cont,
				Time:    tm,
			})
			bc.targetResults.Store(id, obs)

			// and remove target and subtargets from tgts
			bc.removeTargetWithDeps(id, &tgts)
		}
	}

	// cleanup tgts
	var removedIds []string
	for t, meta := range tgts {
		if meta.refCounter <= 0 {
			removedIds = append(removedIds, t)
		}
	}
	for _, t := range removedIds {
		delete(tgts, t)
	}

	// Create observable results
	for t := range tgts {
		bc.targetResults.Store(t, NewObservable())
	}

	bc.mutex.Unlock()

	// Build
	var wg sync.WaitGroup
	wg.Add(len(tgts))

	for id, tgt := range tgts {
		go func(id string, meta targetMeta) {
			logger.Debugf("Building target '%v'...", id)
			sc := NewSessionContext(id, bc, meta.depsNames)
			result := meta.t.Build(sc)
			o, ok := bc.targetResults.Load(id)
			if !ok {
				panic(fmt.Errorf("Target with ID `%v` not found in targetResults", id))
			}
			o.(*ObservableResult).Put(&result)
			logger.Debugf("Building target '%v': done.", id)
			wg.Done()
		}(id, *tgt)
	}

	wg.Wait()

	// cache Cacheable targets
	for id, tgt := range tgts {
		if ct, ok := tgt.t.(targets.Cacheable); ok && ct.Cacheable() {
			o, ok := bc.targetResults.Load(id)
			if !ok {
				panic(fmt.Errorf("Target with ID `%v` not found in targetResults", id))
			}
			res := o.(*ObservableResult).Get()
			if res.Err != nil {
				continue
			}
			bc.globalCache.Put(id, res.Content, res.Time)
		}
	}

	bc.updateBuiltTargets(tgts)

	// return result
	o, ok := bc.targetResults.Load(t.TargetID())
	if !ok {
		panic(fmt.Errorf("Target with ID `%v` not found in targetResults", t.TargetID()))
	}
	br := o.(*ObservableResult).Get()

	return *br
}

func (bc *BuildSession) Builds(ts ...targets.Target) []targets.Result {
	mt := newMultiTarget(ts)
	result := bc.Build(mt)
	if result.Err != nil {
		panic(fmt.Errorf("Internal error: multiTarget Build() failed: %v", result.Err))
	}
	results := result.Content.([]targets.Result)
	return results
}

func (bc *BuildSession) fillTargetDeps(t targets.Target, tgts *map[string]*targetMeta) {
	tid := t.TargetID()

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
			did := d.TargetID()
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
	/*
		if meta.refCounter == 0 {
			delete(*tgts, id)
		}
	*/
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

func (bc *BuildSession) updateBuiltTargets(tgts map[string]*targetMeta) {
	for id := range tgts {
		bc.builtTargets.Store(id, struct{}{})
	}
}

func (bc *BuildSession) BuiltTargets() (ts []string) {
	bc.builtTargets.Range(func(key, value interface{}) bool {
		ts = append(ts, key.(string))
		return true
	})
	sort.Strings(ts)
	return ts
}
