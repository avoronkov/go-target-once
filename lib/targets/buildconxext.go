package targets

import "time"

type BuildContext interface {
	GetDependency(dep string) (interface{}, error)
	Build(t Target) (content interface{}, tm time.Time, err error)
	Builds(ts ...Target) (contents []interface{}, times []time.Time, errs []error)
}
