package builder

import (
	"dont-repeat-twice/lib/targets"
)

type BuildContext struct {
	B *Builder
	T targets.Target
}

var _ targets.BuildContext = (*BuildContext)(nil)

func (bc *BuildContext) GetDependency(dep string) (content interface{}, err error) {
	if targetWithDeps, ok := bc.T.(targets.WithDependencies); ok {
		d, exists := targetWithDeps.Dependencies()[dep]
		if !exists {
			return nil, NewDependencyNotFound(dep)
		}
		return bc.B.Build(d)
	}
	return nil, NewDependencyNotFound(dep)
}
