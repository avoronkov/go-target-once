package builder

import (
	"dont-repeat-twice/lib/id"
	"dont-repeat-twice/lib/targets"
	"dont-repeat-twice/lib/warehouse"
)

var Default *Builder = New(warehouse.NewMemoryWarehouse())

func Build(t targets.Target, args ...id.Interface) interface{} {
	return Default.Build(t, args...)
}
