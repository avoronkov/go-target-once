package targets

import "time"

type keepAlive struct {
	T        Target
	duration time.Duration
}

var _ Target = (*keepAlive)(nil)
var _ Modifiable = (*keepAlive)(nil)
var _ WithDependencies = (*keepAlive)(nil)
var _ IsCacheable = (*keepAlive)(nil)

func KeepAlive(t Target, d time.Duration) Target {
	return &keepAlive{
		T:        t,
		duration: d,
	}
}

func (k *keepAlive) TargetID() string {
	return "keep-alive-" + k.T.TargetID()
}

func (k *keepAlive) IsModified(since time.Time) bool {
	return false
}

func (k *keepAlive) Dependencies() map[string]Target {
	return map[string]Target{
		"target": k.T,
	}
}

func (k *keepAlive) Build(bc BuildContext) Result {
	return bc.GetDependency("target")
}

func (k *keepAlive) ValidFor() time.Duration {
	return k.duration
}

func (k *keepAlive) IsCacheable() bool {
	return true
}

// Some hack to skip dependencies check
func (k *keepAlive) KeepingAlive() {}
