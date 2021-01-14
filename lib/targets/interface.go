package targets

import (
	"time"
)

type Target interface {
	TargetId() string

	Build(bc BuildContext) (interface{}, time.Time, error)
}

type WithDependencies interface {
	Dependencies() map[string]Target
}

type Modifiable interface {
	IsModified(since time.Time) bool
}
