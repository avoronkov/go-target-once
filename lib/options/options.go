package options

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
)

// ID
type ID struct {
	id string
}

func (i *ID) InitTargetId() {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Cannot fetch stack frame."))
	}
	i.id = fmt.Sprintf("id:%v:%v", file, line)
}

func (i *ID) TargetId() string {
	if i.id == "" {
		panic(errors.New("InitTargetId is not called"))
	}
	return i.id
}

// Mutable
type Mutable struct{}

var _ targets.Modifiable = Mutable{}

func (Mutable) IsModified(since time.Time) bool {
	return true
}

// Cachable
type Cachable struct{}

var _ targets.Cachable = Cachable{}

func (Cachable) Cachable() bool {
	return true
}
