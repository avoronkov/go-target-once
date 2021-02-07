package logger

import "log"

var Logger *log.Logger

var Debugf = func(format string, v ...interface{}) {}

var Warningf = func(format string, v ...interface{}) {}
