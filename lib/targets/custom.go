package targets

import (
	"fmt"
	"time"

	"dont-repeat-twice/lib/id"
)

type Custom struct {
	Id string

	Deps []Target

	DoBuild func(this *Custom, args ...id.Interface) (interface{}, time.Time)
}

var _ Target = (*Custom)(nil)

func (c *Custom) TargetId() string {
	if c.Id == "" {
		panic(fmt.Errorf("Custom.Id is not defined"))
	}
	return c.Id
}

func (c *Custom) Build(args ...id.Interface) (content interface{}, t time.Time) {
	return c.DoBuild(c, args...)
}

func (c *Custom) IsModified(since time.Time) bool {
	return false
}

func (c *Custom) Dependencies() []Target {
	return c.Deps
}
