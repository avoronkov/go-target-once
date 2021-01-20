package targets

import "time"

type keepAlive struct {
	T        Target
	duration time.Duration
}

var _ Target = (*keepAlive)(nil)
var _ Modifiable = (*keepAlive)(nil)
var _ WithDependencies = (*keepAlive)(nil)
var _ Cachable = (*keepAlive)(nil)

func KeepAlive(t Target, d time.Duration) Target {
	return &keepAlive{
		T:        t,
		duration: d,
	}
}

func (k *keepAlive) TargetId() string {
	return "keep-alive-" + k.T.TargetId()
}

func (k *keepAlive) IsModified(since time.Time) bool {
	return false
}

func (k *keepAlive) Dependencies() map[string]Target {
	if wd, ok := k.T.(WithDependencies); ok {
		return wd.Dependencies()
	}
	return map[string]Target{}
}

func (k *keepAlive) Build(bc BuildContext) (interface{}, time.Time, error) {
	return k.T.Build(bc)
}

func (k *keepAlive) ValidFor() time.Duration {
	return k.duration
}

func (k *keepAlive) Cachable() bool {
	return true
}
