package builder

import (
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
	"github.com/avoronkov/go-target-once/lib/targets"
)

type BuildContext struct {
	B        *Builder
	T        targets.Target
	contents map[string]chan contentError
}

var _ targets.BuildContext = (*BuildContext)(nil)

func NewBuildContext(b *Builder, t targets.Target) *BuildContext {
	bc := &BuildContext{
		B:        b,
		T:        t,
		contents: make(map[string]chan contentError),
	}
	bc.buildDependencies()
	return bc
}

func (bc *BuildContext) GetDependency(dep string) (content interface{}, err error) {
	ch, ok := bc.contents[dep]
	if !ok {
		return nil, NewDependencyNotFound(dep)
	}
	ce := <-ch
	return ce.content, ce.err
}

func (bc *BuildContext) Close() error {
	logger.Debugf("BuildContext: closing...")
	for _, ch := range bc.contents {
		<-ch
	}
	logger.Debugf("BuildContext: closed.")
	return nil
}

func (bc *BuildContext) buildDependencies() {
	tdeps, ok := bc.T.(targets.WithDependencies)
	if !ok {
		return
	}

	deps := tdeps.Dependencies()
	for name := range deps {
		bc.contents[name] = make(chan contentError)
	}

	// Starting building dependencies
	go func() {
		for name, dep := range deps {
			go func(n string, d targets.Target) {
				defer close(bc.contents[n])
				cont, tm, err := bc.B.Build(d)
				bc.contents[n] <- contentError{
					content: cont,
					tm:      tm,
					err:     err,
				}
			}(name, dep)
		}
	}()
}

func (bc *BuildContext) Build(t targets.Target) (content interface{}, tm time.Time, err error) {
	return bc.B.Build(t)
}

func (bc *BuildContext) Builds(ts ...targets.Target) (contents []interface{}, times []time.Time, errs []error) {
	return bc.B.Builds(ts...)
}
