package warehouse

import (
	"time"
)

type Warehouse interface {
	Get(targetId string) (content interface{}, t time.Time, ok bool)
	Put(targetId string, content interface{}, t time.Time)
}
