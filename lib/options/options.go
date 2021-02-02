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

func (i *ID) InitTargetID() {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Cannot fetch stack frame."))
	}
	i.id = fmt.Sprintf("id:%v:%v", file, line)
}

func (i *ID) TargetID() string {
	if i.id == "" {
		panic(errors.New("InitTargetID is not called"))
	}
	return i.id
}

// Mutable
type Mutable struct{}

var _ targets.Modifiable = Mutable{}

func (Mutable) IsModified(since time.Time) bool {
	return true
}

// Cacheable
type Cacheable struct{}

var _ targets.Cacheable = Cacheable{}

func (Cacheable) Cacheable() bool {
	return true
}
