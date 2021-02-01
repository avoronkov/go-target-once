package targets

import "time"

type Result struct {
	Content interface{}
	Time    time.Time
	Err     error
}

func ResultOk(content interface{}) Result {
	return Result{
		Content: content,
		Time:    time.Now(),
	}
}

func ResultOkTime(content interface{}, tm time.Time) Result {
	return Result{
		Content: content,
		Time:    tm,
	}
}

func ResultFailed(err error) Result {
	return Result{
		Err: err,
	}
}
