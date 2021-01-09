package warehouse

import (
	"time"
)

var _ Warehouse = (*MemoryWarehouse)(nil)

type MemoryWarehouse struct {
	data map[string]memRecord
}

func NewMemoryWarehouse() *MemoryWarehouse {
	return &MemoryWarehouse{
		data: make(map[string]memRecord),
	}
}

func (w *MemoryWarehouse) Get(targetId string) (content interface{}, t time.Time, ok bool) {
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

	w.data[targetId] = mr
}

type memRecord struct {
	Content interface{}
	Time    time.Time
}
