package targets

import (
	"time"
)

type Target interface {
	TargetId() string
	IsModified(since time.Time) bool

	Build(bc BuildContext) (interface{}, time.Time, error)
}

type WithDependencies interface {
	Dependencies() map[string]Target
}
