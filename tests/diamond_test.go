package tests

import (
	"fmt"
	"sync/atomic"

	"github.com/avoronkov/go-target-once/lib/builder"
	op "github.com/avoronkov/go-target-once/lib/options"
	tg "github.com/avoronkov/go-target-once/lib/targets"
)

var (
	cnA, cnB, cnC int32
)

type Ta struct {
	op.ID
}

var _ tg.Target = (*Ta)(nil)
var _ tg.WithDependencies = (*Ta)(nil)

func newTa() (t *Ta) {
	defer func() { t.InitTargetID() }()
	return new(Ta)
}

func (t *Ta) Dependencies() map[string]tg.Target {
	return map[string]tg.Target{
		"dep-b": newTb(),
		"dep-c": newTc(),
	}
}

func (t *Ta) Build(bc tg.BuildContext) tg.Result {
	atomic.AddInt32(&cnA, 1)
	b := bc.GetDependency("dep-b")
	c := bc.GetDependency("dep-c")
	return tg.OK(fmt.Sprintf("A{b: %v, c: %v}", b.Content, c.Content))
}

type Tb struct {
	op.ID
}

var _ tg.Target = (*Tb)(nil)
var _ tg.WithDependencies = (*Tb)(nil)

func newTb() (t *Tb) {
	defer func() { t.InitTargetID() }()
	return new(Tb)
}

func (t *Tb) Dependencies() map[string]tg.Target {
	return map[string]tg.Target{
		"dep-c": newTc(),
	}
}

func (t *Tb) Build(bc tg.BuildContext) tg.Result {
	atomic.AddInt32(&cnB, 1)
	c := bc.GetDependency("dep-c")
	return tg.OK(fmt.Sprintf("B{c: %v}", c.Content))
}

type Tc struct {
	op.ID
}

var _ tg.Target = (*Tc)(nil)

func newTc() (t *Tc) {
	defer func() { t.InitTargetID() }()
	return new(Tc)
}

func (t *Tc) Build(bc tg.BuildContext) tg.Result {
	atomic.AddInt32(&cnC, 1)
	return tg.OK("{C}")
}

// Tests

func ExampleDiamondDeps() {
	t := newTa()
	res := builder.Build(t)

	fmt.Printf("T: %v (%v)\n", res.Content, res.Err)

	fmt.Printf("A built %v times.\n", cnA)
	fmt.Printf("B built %v times.\n", cnB)
	fmt.Printf("C built %v times.\n", cnC)
	// Output:
	// T: A{b: B{c: {C}}, c: {C}} (<nil>)
	// A built 1 times.
	// B built 1 times.
	// C built 1 times.
}
