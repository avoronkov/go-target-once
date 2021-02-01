package builder

import (
	"github.com/avoronkov/go-target-once/lib/targets"
	"github.com/avoronkov/go-target-once/lib/warehouse"
)

var DefaultWarehouse warehouse.Warehouse = warehouse.NewMemoryWarehouse()

func Build(t targets.Target) targets.Result {
	bs := NewBuildSession(DefaultWarehouse)
	return bs.Build(t)
}

func Builds(ts ...targets.Target) []targets.Result {
	bs := NewBuildSession(DefaultWarehouse)
	return bs.Builds(ts...)
}
