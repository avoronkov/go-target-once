package options

import (
	"testing"
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
)

type CacheableStruct struct {
	Cacheable
}

var _ (targets.IsCacheable) = CacheableStruct{}
var _ (targets.IsCacheable) = (*CacheableStruct)(nil)

func TestCacheable(t *testing.T) {
	c := CacheableStruct{}
	if !c.IsCacheable() {
		t.Errorf("CacheableStruct is not really cacheable")
	}
}

type IDStruct struct {
	ID
}

var _ (targets.Target) = (*IDStruct)(nil)

func (i *IDStruct) Build(bc targets.BuildContext) targets.Result {
	return targets.OK("OK!")
}

func newIDStruct() (t *IDStruct) {
	defer func() { t.SetTargetID("IDStruct") }()
	return new(IDStruct)
}

func TestID(t *testing.T) {
	i := newIDStruct()
	if act, exp := i.TargetID(), "IDStruct"; act != exp {
		t.Errorf("Incorrect TargetID(): got %v, want %v", act, exp)
	}
}

type MutableStruct struct {
	Mutable
}

var _ targets.Modifiable = MutableStruct{}
var _ targets.Modifiable = (*MutableStruct)(nil)

func TestMutable(t *testing.T) {
	m := MutableStruct{}
	if !m.IsModified(time.Now()) {
		t.Errorf("MutableStruct is not really mutable")
	}
}
