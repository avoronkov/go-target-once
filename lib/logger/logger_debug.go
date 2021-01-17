// +build debug

package logger

import (
	"log"
	"os"
)

func init() {
	Logger = log.New(os.Stderr, "[go-target-once] ", log.LstdFlags)
}
