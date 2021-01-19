package targets

import "time"

type keepAlive struct {
	T        Target
	duration time.Duration
}

var _ Target = (*keepAlive)(nil)
var _ Modifiable = (*File)(nil)

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

func (k *keepAlive) Build(bc BuildContext) (interface{}, time.Time, error) {
	return k.T.Build(bc)
}

func (k *keepAlive) ValidFor() time.Duration {
	return k.duration
}
