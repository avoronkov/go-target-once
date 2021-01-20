// +build debug

package logger

import (
	"log"
	"os"
)

func init() {
	Logger = log.New(os.Stderr, "[go-target-once] ", log.LstdFlags)
}

func Debugf(format string, v ...interface{}) {
	Logger.Printf("[debug] "+format, v...)
}

func Warningf(format string, v ...interface{}) {
	Logger.Printf("[warning] "+format, v...)
}
