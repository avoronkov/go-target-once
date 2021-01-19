package warehouse

import (
	"time"
)

type Warehouse interface {
	Get(targetId string) (content interface{}, t time.Time, ok bool)
	Put(targetId string, content interface{}, t time.Time, opts ...Option)
}

type Option interface {
}

type ValidForOption interface {
	ValidFor() time.Duration
}

type validFor struct {
	d time.Duration
}

var _ ValidForOption = validFor{}

func OptValidFor(d time.Duration) Option {
	return validFor{d: d}
}

func (v validFor) ValidFor() time.Duration {
	return v.d
}
