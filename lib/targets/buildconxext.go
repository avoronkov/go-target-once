package targets

type BuildContext interface {
	GetDependency(dep string) (interface{}, error)
}
