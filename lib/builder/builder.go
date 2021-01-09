package builder

import (
	"dont-repeat-twice/lib/targets"
	"dont-repeat-twice/lib/warehouse"
	"log"
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

	log.Printf("[debug] (%v) Rebuild content", t.TargetId())
	bc := &BuildContext{
		B: b,
		T: t,
	}
	cont, tm, err := t.Build(bc)
	if err != nil {
		return cont, err
	}

	b.w.Put(t.TargetId(), cont, tm)
	return cont, nil
}

func (b *Builder) isReady(t targets.Target) (bool, interface{}) {
	// Search for saved value
	cont, tm, ok := b.w.Get(t.TargetId())
	if !ok {
		return false, nil
	}

	if t.IsModified(tm) {
		return false, nil
	}

	for _, dep := range t.Dependencies() {
		if dep.IsModified(tm) {
			return false, nil
		}
	}
	return true, cont
}
