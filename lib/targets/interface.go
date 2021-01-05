package targets

import (
	"dont-repeat-twice/lib/id"
	"time"
)

type Target interface {
	TargetId() string
	IsModified(since time.Time) bool
	Dependencies() []Target

	Build(bc BuildContext, args ...id.Interface) (interface{}, time.Time, error)
}
