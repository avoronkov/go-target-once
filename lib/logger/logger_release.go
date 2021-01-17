// +build !debug

package logger

import (
	"io/ioutil"
	"log"
)

func init() {
	Logger = log.New(ioutil.Discard, "[go-target-once] ", log.LstdFlags)
}
