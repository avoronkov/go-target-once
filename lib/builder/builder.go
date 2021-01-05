package builder

import (
	"dont-repeat-twice/lib/id"
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

func (b *Builder) Build(t targets.Target, args ...id.Interface) (content interface{}) {
	ready, cont := b.isReady(t, args...)
	if ready {
		log.Printf("[debug] (%v) Return content from cache.", t.TargetId())
		return cont
	}

	log.Printf("[debug] (%v) Rebuild content", t.TargetId())
	bc := &BuildContext{
		B: b,
		T: t,
	}
	cont, tm := t.Build(bc, args...)
	b.w.Put(t.TargetId(), args, cont, tm)
	return cont
}

func (b *Builder) isReady(t targets.Target, args ...id.Interface) (bool, interface{}) {
	// Search for saved value
	cont, tm, ok := b.w.Get(t.TargetId(), args)
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
