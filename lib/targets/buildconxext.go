package targets

type BuildContext interface {
	GetDependency(dep string) (interface{}, error)
	Build(t Target) (content interface{}, err error)
}
