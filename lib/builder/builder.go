package builder

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
	"github.com/avoronkov/go-target-once/lib/targets"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

type Builder struct {
	w warehouse.Warehouse
}

func New(w warehouse.Warehouse) *Builder {
	return &Builder{
		w: w,
	}
}

func (b *Builder) Builds(ts ...targets.Target) (contents []interface{}, errs []error) {
	l := len(ts)
	contents = make([]interface{}, l)
	errs = make([]error, l)
	var wg sync.WaitGroup
	wg.Add(l)
	for index, target := range ts {
		go func(i int, t targets.Target) {
			defer wg.Done()
			content, err := b.Build(t)
			contents[i] = content
			errs[i] = err
		}(index, target)
	}
	wg.Wait()
	return
}

func (b *Builder) Build(t targets.Target) (content interface{}, err error) {
	ready, cont := b.isReady(t)
	if ready {
		logger.Debugf("(%v) Return content from cache.", t.TargetId())
		return cont, nil
	}

	bc := NewBuildContext(b, t)

	logger.Debugf("(%v) Rebuild content", t.TargetId())
	cont, tm, err := t.Build(bc)
	if err != nil {
		return cont, err
	}

	go bc.Close()

	// Cache only cachable targets
	if c, ok := t.(targets.Cachable); ok && c.Cachable() {
		var opts []warehouse.Option
		if v, ok := t.(targets.ValidFor); ok {
			opts = append(opts, warehouse.OptValidFor(v.ValidFor()))
		}
		b.w.Put(t.TargetId(), cont, tm, opts...)
	}

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
