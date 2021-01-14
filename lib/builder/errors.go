package builder

import "fmt"

type DependencyNotFound struct {
	dep string
}

func NewDependencyNotFound(name string) *DependencyNotFound {
	return &DependencyNotFound{dep: name}
}

func (d *DependencyNotFound) Error() string {
	return fmt.Sprintf("Depndency not found: %v", d.dep)
}
