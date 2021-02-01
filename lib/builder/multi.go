package builder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type multiTarget struct {
	id   string
	deps map[string]targets.Target
}

var _ targets.Target = (*multiTarget)(nil)
var _ targets.WithDependencies = (*multiTarget)(nil)

func newMultiTarget(ts []targets.Target) *multiTarget {
	id := &strings.Builder{}
	fmt.Fprintf(id, "multi:")
	deps := make(map[string]targets.Target)
	for i, t := range ts {
		if i > 0 {
			fmt.Fprintf(id, ",")
		}
		fmt.Fprintf(id, t.TargetId())
		deps[strconv.Itoa(i)] = t
	}
	return &multiTarget{
		id:   id.String(),
		deps: deps,
	}
}

func (t *multiTarget) TargetId() string {
	return t.id
}

func (t *multiTarget) Dependencies() map[string]targets.Target {
	return t.deps
}

// Content = []targets.Result
func (t *multiTarget) Build(bc targets.BuildContext) targets.Result {
	var results []targets.Result
	l := len(t.deps)
	for i := 0; i < l; i++ {
		d := bc.GetDependency(strconv.Itoa(i))
		results = append(results, d)
	}
	return targets.ResultOk(results)
}
