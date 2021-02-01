package targets

type BuildContext interface {
	GetDependency(dep string) Result
	Build(t Target) Result
	Builds(ts ...Target) []Result
}
