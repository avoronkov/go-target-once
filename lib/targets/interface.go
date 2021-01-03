package targets

import (
	"dont-repeat-twice/lib/id"
	"time"
)

type Target interface {
	TargetId() string
	IsModified(since time.Time) bool
	Dependencies() []Target

	Build(args ...id.Interface) (content interface{}, t time.Time)
}
