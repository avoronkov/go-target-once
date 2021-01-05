package builder

import (
	"dont-repeat-twice/lib/id"
	"dont-repeat-twice/lib/targets"
	"time"
)

type BuildContext struct {
	B *Builder
	T targets.Target
}

var _ targets.BuildContext = (*BuildContext)(nil)

func (bc *BuildContext) GetDependency(dep int, args ...id.Interface) (content interface{}, t time.Time) {
	d := bc.T.Dependencies()[dep]
	dbc := &BuildContext{
		B: bc.B,
		T: d,
	}
	return d.Build(dbc, args...)
}
