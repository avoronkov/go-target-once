package builder

import (
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

var Default *Builder = New(warehouse.NewMemoryWarehouse())

func Build(t targets.Target) (interface{}, time.Time, error) {
	return Default.Build(t)
}

func Builds(ts ...targets.Target) ([]interface{}, []time.Time, []error) {
	return Default.Builds(ts...)
}
