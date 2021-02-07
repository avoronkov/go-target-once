// +build debug

package logger

import (
	"log"
	"os"
)

func init() {
	Logger = log.New(os.Stderr, "[go-target-once] ", log.LstdFlags)
	Debugf = debugf
	Warningf = warningf
}

func debugf(format string, v ...interface{}) {
	Logger.Printf("[debug] "+format, v...)
}

func warningf(format string, v ...interface{}) {
	Logger.Printf("[warning] "+format, v...)
}
