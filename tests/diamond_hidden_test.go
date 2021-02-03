package tests

import (
	"fmt"
	"sync/atomic"

	"github.com/avoronkov/go-target-once/lib/builder"
	op "github.com/avoronkov/go-target-once/lib/options"
	tg "github.com/avoronkov/go-target-once/lib/targets"
)

var (
	cxa, cxb, cxc int32
)

type xa struct {
	op.ID
}

var _ tg.Target = (*xa)(nil)

func newxa() (t *xa) {
	defer func() { t.InitTargetID() }()
	return new(xa)
}

func (t *xa) Build(bc tg.BuildContext) tg.Result {
	atomic.AddInt32(&cxa, 1)
	ct := newxc()
	c := bc.Build(ct)
	return tg.OK(fmt.Sprintf("A{c: %v}", c.Content))
}

type xb struct {
	op.ID
}

var _ tg.Target = (*xb)(nil)

func newxb() (t *xb) {
	defer func() { t.InitTargetID() }()
	return new(xb)
}

func (t *xb) Build(bc tg.BuildContext) tg.Result {
	atomic.AddInt32(&cxb, 1)
	ct := newxc()
	c := bc.Build(ct)
	return tg.OK(fmt.Sprintf("B{c: %v}", c.Content))
}

type xc struct {
	op.ID
}

var _ tg.Target = (*xc)(nil)

func newxc() (t *xc) {
	defer func() { t.InitTargetID() }()
	return new(xc)
}

func (t *xc) Build(bc tg.BuildContext) tg.Result {
	atomic.AddInt32(&cxc, 1)
	return tg.OK("{C}")
}

func ExampleHiddedDiamondDependency() {
	a := newxa()
	b := newxb()
	results := builder.Builds(a, b)
	for _, res := range results {
		fmt.Printf("Result: %v (%v)\n", res.Content, res.Err)
	}
	fmt.Printf("A built %v times.\n", cxa)
	fmt.Printf("B built %v times.\n", cxb)
	fmt.Printf("C built %v times.\n", cxc)
	// Output:
	// Result: A{c: {C}} (<nil>)
	// Result: B{c: {C}} (<nil>)
	// A built 1 times.
	// B built 1 times.
	// C built 1 times.
}
