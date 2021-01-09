package builder

import (
	"dont-repeat-twice/lib/id"
	"dont-repeat-twice/lib/targets"
)

type BuildContext struct {
	B *Builder
	T targets.Target
}

var _ targets.BuildContext = (*BuildContext)(nil)

func (bc *BuildContext) GetDependency(dep int, args ...id.Interface) (content interface{}, err error) {
	d := bc.T.Dependencies()[dep]
	return bc.B.Build(d)
}
