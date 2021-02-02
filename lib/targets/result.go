package targets

import "time"

type Result struct {
	Content interface{}
	Time    time.Time
	Err     error
}

func OK(content interface{}) Result {
	return Result{
		Content: content,
		Time:    time.Now(),
	}
}

func OKTime(content interface{}, tm time.Time) Result {
	return Result{
		Content: content,
		Time:    tm,
	}
}

func Failed(err error) Result {
	return Result{
		Err: err,
	}
}
