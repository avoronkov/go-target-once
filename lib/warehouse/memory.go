package warehouse

import (
	"fmt"
	"sync"
	"time"
)

var _ Warehouse = (*MemoryWarehouse)(nil)

type MemoryWarehouse struct {
	data  map[string]memRecord
	mutex sync.RWMutex
}

func NewMemoryWarehouse() *MemoryWarehouse {
	return &MemoryWarehouse{
		data: make(map[string]memRecord),
	}
}

func (w *MemoryWarehouse) Get(targetId string) (content interface{}, t time.Time, ok bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	mr, exists := w.data[targetId]
	if !exists {
		return
	}
	if !mr.Expires.IsZero() && time.Now().After(mr.Expires) {
		delete(w.data, targetId)
		return
	}
	return mr.Content, mr.Time, true
}

func (w *MemoryWarehouse) Put(targetId string, content interface{}, t time.Time, opts ...Option) {
	mr := memRecord{
		Content: content,
		Time:    t,
	}

	for _, opt := range opts {
		switch a := opt.(type) {
		case ValidForOption:
			mr.Expires = time.Now().Add(a.ValidFor())
		default:
			panic(fmt.Errorf("Unknown option: %v (%T)", opt, opt))
		}
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.data[targetId] = mr
}

type memRecord struct {
	Content interface{}
	Time    time.Time
	Expires time.Time
}
