// Code generated by go2go; DO NOT EDIT.


//line syncholder.go2:2
package syncholder

//line syncholder.go2:2
import (
//line syncholder.go2:2
 "github.com/avoronkov/go-target-once/lib/targets"
//line syncholder.go2:2
 "sync"
//line syncholder.go2:2
)

//line syncholder.go2:58
type ResultSyncHolder = instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult

func NewResultSyncHolder() *ResultSyncHolder {
	return instantiate୦୦NewSyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult()
}

//line syncholder.go2:62
type instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult struct {
//line syncholder.go2:12
 data  *targets.Result
				ready bool
				mutex sync.Mutex

				observers int
				notify    chan *targets.Result
}

//line syncholder.go2:26
func (o *instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult,) Put(data *targets.Result) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.data = data
	o.ready = true

	if o.observers >= 0 {
		for i := 0; i < o.observers; i++ {
			o.notify <- data
		}
		close(o.notify)
		o.observers = -1
	}
}

func (o *instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult,) Get() *targets.Result {
	o.mutex.Lock()

	if o.ready {
		o.mutex.Unlock()
		return o.data
	}

//line syncholder.go2:51
 o.observers++
	o.mutex.Unlock()

	return <-o.notify
}
//line syncholder.go2:20
func instantiate୦୦NewSyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult() *instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult {
	return &instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult{
		notify: make(chan *targets.Result),
	}
}

//line syncholder.go2:24
type Importable୦ int
//line syncholder.go2:24
type _ targets.BuildContext
//line syncholder.go2:24
type _ sync.Cond
