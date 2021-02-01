package builder

import (
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type SessionContext struct {
	// dep name -> dep id
	depNames map[string]string

	targetId     string
	buildSession *BuildSession
}

var _ targets.BuildContext = (*SessionContext)(nil)

func NewSessionContext(targetId string, bs *BuildSession, depNames map[string]string) *SessionContext {
	return &SessionContext{
		targetId:     targetId,
		buildSession: bs,
		depNames:     depNames,
	}
}

func (sc *SessionContext) GetDependency(name string) (interface{}, error) {
	depId, ok := sc.depNames[name]
	if !ok {
		return nil, NewDependencyNotFound(name)
	}
	o, ok := sc.buildSession.targetResults.Load(depId)
	if !ok {
		return nil, NewDependencyNotFound(name)
	}
	br := o.(*ObservableResult).Get()
	return br.C, br.E
}

func (sc *SessionContext) Build(t targets.Target) (interface{}, time.Time, error) {
	return sc.buildSession.Build(t)
}

func (sc *SessionContext) Builds(t ...targets.Target) ([]interface{}, []time.Time, []error) {
	// TODO
	return nil, nil, nil
}
