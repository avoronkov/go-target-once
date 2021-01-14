package targets

import (
	"fmt"
	"time"
)

type Common struct {
	Id   string
	Deps map[string]Target
}

func (c Common) TargetId() string {
	if c.Id == "" {
		panic(fmt.Errorf("Custom.Id is not defined"))
	}
	return c.Id
}

func (c Common) Dependencies() map[string]Target {
	return c.Deps
}

func (c Common) DepsModified(since time.Time) bool {
	for _, t := range c.Deps {
		if t.IsModified(since) {
			return true
		}
	}
	return false
}
