package builder

import (
	"sync"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type ObservableResult struct {
	data  *targets.Result
	ready bool
	mutex sync.Mutex

	observers []chan *targets.Result
}

func NewObservable() *ObservableResult {
	return &ObservableResult{}
}

func (o *ObservableResult) Put(data *targets.Result) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.data = data
	o.ready = true

	for _, observer := range o.observers {
		observer <- data
		close(observer)
	}
	o.observers = []chan *targets.Result{}
}

func (o *ObservableResult) Get() *targets.Result {
	o.mutex.Lock()

	if o.ready {
		o.mutex.Unlock()
		return o.data
	}

	// wait for data is ready
	ch := make(chan *targets.Result)
	o.observers = append(o.observers, ch)
	o.mutex.Unlock()

	return <-ch
}
