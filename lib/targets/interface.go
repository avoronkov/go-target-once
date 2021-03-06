package targets

import (
	"time"
)

type Target interface {
	TargetID() string

	Build(bc BuildContext) Result
}

type WithDependencies interface {
	Dependencies() map[string]Target
}

type Modifiable interface {
	IsModified(since time.Time) bool
}

type IsCacheable interface {
	IsCacheable() bool
}

type ValidFor interface {
	ValidFor() time.Duration
}

type KeepingAlive interface {
	KeepingAlive()
}
