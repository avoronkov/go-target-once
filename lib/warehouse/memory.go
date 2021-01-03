package warehouse

import (
	"dont-repeat-twice/lib/id"
	"time"
)

var _ Warehouse = (*MemoryWarehouse)(nil)

type MemoryWarehouse struct {
	data map[string]map[string]memRecord
}

func NewMemoryWarehouse() *MemoryWarehouse {
	return &MemoryWarehouse{
		data: make(map[string]map[string]memRecord),
	}
}

func (w *MemoryWarehouse) Get(targetId string, args []id.Interface) (content interface{}, t time.Time, ok bool) {
	msm, exists := w.data[targetId]
	if !exists {
		return
	}
	argsId := id.Of(args...)
	mr, exists := msm[argsId]
	if !exists {
		return
	}
	return mr.Content, mr.Time, true
}

func (w *MemoryWarehouse) Put(targetId string, args []id.Interface, content interface{}, t time.Time) {
	argsId := id.Of(args...)
	mr := memRecord{
		Content: content,
		Time:    t,
	}

	_, exists := w.data[targetId]
	if !exists {
		w.data[targetId] = map[string]memRecord{
			argsId: mr,
		}
	} else {
		w.data[targetId][argsId] = mr
	}
}

type memRecord struct {
	Content interface{}
	Time    time.Time
}
