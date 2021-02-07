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

//line syncholder.go2:46
type ResultSyncHolder = instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult

func NewResultSyncHolder() *ResultSyncHolder {
	return instantiate୦୦NewSyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult()
}

//line syncholder.go2:50
type instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult struct {
//line syncholder.go2:12
 data  *targets.Result
				ready bool

				cond *sync.Cond
}

//line syncholder.go2:24
func (o *instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult,) Put(data *targets.Result) {
	o.cond.L.Lock()
	defer o.cond.L.Unlock()

	o.data = data
	o.ready = true

	o.cond.Broadcast()
}

func (o *instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult,) Get() *targets.Result {
	o.cond.L.Lock()
	defer o.cond.L.Unlock()

	for !o.ready {
		o.cond.Wait()
	}

	return o.data
}
//line syncholder.go2:18
func instantiate୦୦NewSyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult() *instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult {
	return &instantiate୦୦SyncHolder୦୮1github୮acom୮davoronkov୮dgo୮ctarget୮conce୮dlib୮dtargets୮aResult{
		cond: sync.NewCond(new(sync.Mutex)),
	}
}

//line syncholder.go2:22
type Importable୦ int
//line syncholder.go2:22
type _ targets.BuildContext
//line syncholder.go2:22
type _ sync.Cond
