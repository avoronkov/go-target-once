package logger

import "log"

var Logger *log.Logger

func Debugf(format string, v ...interface{}) {
	Logger.Printf("[debug] "+format, v...)
}

func Warningf(format string, v ...interface{}) {
	Logger.Printf("[warning] "+format, v...)
}
