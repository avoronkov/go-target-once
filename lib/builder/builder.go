package builder

import (
	"dont-repeat-twice/lib/targets"
	"dont-repeat-twice/lib/warehouse"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Builder struct {
	w warehouse.Warehouse
}

func New(w warehouse.Warehouse) *Builder {
	return &Builder{
		w: w,
	}
}

func (b *Builder) Build(t targets.Target) (content interface{}, err error) {
	ready, cont := b.isReady(t)
	if ready {
		log.Printf("[debug] (%v) Return content from cache.", t.TargetId())
		return cont, nil
	}

	bc := NewBuildContext(b, t)

	log.Printf("[debug] (%v) Rebuild content", t.TargetId())
	cont, tm, err := t.Build(bc)
	if err != nil {
		return cont, err
	}

	go bc.Close()

	b.w.Put(t.TargetId(), cont, tm)
	return cont, nil
}

func (b *Builder) IsModified(t targets.Target, since time.Time) bool {
	var modified int32
	var wg sync.WaitGroup

	if mod, ok := t.(targets.Modifiable); ok {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if mod.IsModified(since) {
				atomic.StoreInt32(&modified, 1)
			}
		}()
	}

	if wd, ok := t.(targets.WithDependencies); ok {
		deps := wd.Dependencies()
		wg.Add(len(deps))
		for _, dep := range wd.Dependencies() {
			go func(d targets.Target) {
				defer wg.Done()
				if b.IsModified(d, since) {
					atomic.StoreInt32(&modified, 1)
				}
			}(dep)
		}
	}
	return modified > 0
}

func (b *Builder) isReady(t targets.Target) (bool, interface{}) {
	// Search for saved value
	cont, tm, ok := b.w.Get(t.TargetId())
	if !ok {
		return false, nil
	}

	if mod, ok := t.(targets.Modifiable); ok {
		if mod.IsModified(tm) {
			return false, nil
		}
	}

	if b.IsModified(t, tm) {
		return false, nil
	}

	if targetWithDeps, ok := t.(targets.WithDependencies); ok {
		for _, dep := range targetWithDeps.Dependencies() {
			// There is an error here
			if mod, ok := dep.(targets.Modifiable); ok {
				if mod.IsModified(tm) {
					return false, nil
				}
			}
		}
	}
	return true, cont
}
