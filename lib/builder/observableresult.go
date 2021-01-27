package builder

import "sync"

type ObservableResult struct {
	data  *buildResult
	ready bool
	mutex sync.Mutex

	observers []chan *buildResult
}

func NewObservable() *ObservableResult {
	return &ObservableResult{}
}

func (o *ObservableResult) Put(data *buildResult) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.data = data
	o.ready = true

	for _, observer := range o.observers {
		observer <- data
		close(observer)
	}
	o.observers = []chan *buildResult{}
}

func (o *ObservableResult) Get() *buildResult {
	o.mutex.Lock()

	if o.ready {
		o.mutex.Unlock()
		return o.data
	}

	// wait for data is ready
	ch := make(chan *buildResult)
	o.observers = append(o.observers, ch)
	o.mutex.Unlock()

	return <-ch
}
