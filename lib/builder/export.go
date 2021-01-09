package builder

import (
	"dont-repeat-twice/lib/targets"
	"dont-repeat-twice/lib/warehouse"
)

var Default *Builder = New(warehouse.NewMemoryWarehouse())

func Build(t targets.Target) (interface{}, error) {
	return Default.Build(t)
}
