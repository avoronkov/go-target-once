package targets

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
)

type File struct {
	path string
}

var _ Target = (*File)(nil)
var _ Modifiable = (*File)(nil)

func NewFile(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) TargetId() string {
	return f.path
}

func (f *File) Build(bc BuildContext) (content interface{}, t time.Time, e error) {
	file, err := os.Open(f.path)
	if err != nil {
		e = err
		return
	}

	fi, err := file.Stat()
	if err != nil {
		e = err
		return
	}

	content, err = ioutil.ReadAll(file)
	if err != nil {
		e = err
		return
	}

	return content, fi.ModTime(), nil
}

func (f *File) IsModified(since time.Time) bool {
	fi, err := os.Stat(f.path)
	if err != nil {
		logger.Warningf("Cannot stat file: %v", err)
		return true
	}
	return fi.ModTime().After(since)
}
