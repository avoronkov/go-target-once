package targets

import (
	"dont-repeat-twice/lib/id"
	"time"
)

type BuildContext interface {
	GetDependency(dep int, args ...id.Interface) (content interface{}, t time.Time)
}
