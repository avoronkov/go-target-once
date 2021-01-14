package warehouse

import (
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
	return mr.Content, mr.Time, true
}

func (w *MemoryWarehouse) Put(targetId string, content interface{}, t time.Time) {
	mr := memRecord{
		Content: content,
		Time:    t,
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.data[targetId] = mr
}

type memRecord struct {
	Content interface{}
	Time    time.Time
}
