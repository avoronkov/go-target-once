package tests

import (
	"fmt"
	"time"

	"github.com/avoronkov/go-target-once/lib/builder"
	"github.com/avoronkov/go-target-once/lib/targets"
)

// TestTarget

type TestTarget struct{}

var _ targets.Target = (*TestTarget)(nil)
var _ targets.WithDependencies = (*TestTarget)(nil)
var _ targets.Modifiable = (*TestTarget)(nil)

func (t *TestTarget) TargetID() string {
	return "test-target"
}

func (t *TestTarget) IsModified(since time.Time) bool {
	return true
}

func (t *TestTarget) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"dep": &TestDependency{},
	}
}

func (t *TestTarget) Build(bc targets.BuildContext) targets.Result {
	dep := bc.GetDependency("dep")
	if dep.Err != nil {
		return targets.Failed(dep.Err)
	}
	fmt.Println("Target: Build()")
	return targets.OK(fmt.Sprintf("Target result: %v", dep.Content))
}

// TestDependency

type TestDependency struct{}

var _ targets.Target = (*TestDependency)(nil)
var _ targets.Modifiable = (*TestDependency)(nil)

func (d *TestDependency) TargetID() string {
	return "test-dependency"
}

func (d *TestDependency) Build(bc targets.BuildContext) targets.Result {
	fmt.Println("Dependency: Build()")
	return targets.OK("Dependency result")
}

func (d *TestDependency) IsModified(since time.Time) bool {
	return true
}

// Tests

func ExampleKeepAlive() {
	t1 := targets.KeepAlive(&TestTarget{}, 1*time.Minute)
	res := builder.Build(t1)
	fmt.Printf("Built (1): %v (%v)\n", res.Content, res.Err)

	t2 := targets.KeepAlive(&TestTarget{}, 1*time.Minute)
	res = builder.Build(t2)
	fmt.Printf("Built (2): %v (%v)\n", res.Content, res.Err)
	// Output:
	// Dependency: Build()
	// Target: Build()
	// Built (1): Target result: Dependency result (<nil>)
	// Built (2): Target result: Dependency result (<nil>)
}
