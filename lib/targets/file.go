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

// Content = []byte
func (f *File) Build(bc BuildContext) Result {
	file, err := os.Open(f.path)
	if err != nil {
		return ResultFailed(err)
	}

	fi, err := file.Stat()
	if err != nil {
		return ResultFailed(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return ResultFailed(err)
	}

	return ResultOkTime(content, fi.ModTime())
}

func (f *File) IsModified(since time.Time) bool {
	fi, err := os.Stat(f.path)
	if err != nil {
		logger.Warningf("Cannot stat file: %v", err)
		return true
	}
	return fi.ModTime().After(since)
}
