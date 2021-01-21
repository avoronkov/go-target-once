package builder

import "time"

type contentError struct {
	content interface{}
	tm      time.Time
	err     error
}
