package id

type Interface interface {
	Id() string
}

func Of(args ...Interface) (res string) {
	for i, arg := range args {
		if i > 0 {
			res += "//"
		}
		res += arg.Id()
	}
	return
}
