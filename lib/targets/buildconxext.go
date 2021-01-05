package targets

import (
	"dont-repeat-twice/lib/id"
)

type BuildContext interface {
	GetDependency(dep int, args ...id.Interface) interface{}
}
