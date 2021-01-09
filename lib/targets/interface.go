package targets

import (
	"time"
)

type Target interface {
	TargetId() string
	IsModified(since time.Time) bool
	Dependencies() []Target

	Build(bc BuildContext) (interface{}, time.Time, error)
}
