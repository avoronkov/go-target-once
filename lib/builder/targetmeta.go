package builder

import "github.com/avoronkov/go-target-once/lib/targets"

type targetMeta struct {
	t targets.Target

	// dep name -> dep id
	depsNames map[string]string

	refCounter int
}

func newTargetMeta() *targetMeta {
	return &targetMeta{
		depsNames:  map[string]string{},
		refCounter: 1,
	}
}
