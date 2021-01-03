package warehouse

import (
	"dont-repeat-twice/lib/id"
	"time"
)

type Warehouse interface {
	Get(targetId string, args []id.Interface) (content interface{}, t time.Time, ok bool)
	Put(targetId string, args []id.Interface, content interface{}, t time.Time)
}
