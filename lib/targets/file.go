package targets

import (
	"dont-repeat-twice/lib/id"
	"log"
	"os"
	"time"
)

type File struct {
	path string
}

var _ Target = (*File)(nil)

func NewFile(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) TargetId() string {
	return f.path
}

func (f *File) Build(args ...id.Interface) (content interface{}, t time.Time) {
	return f.path, time.Now()
}

func (f *File) IsModified(since time.Time) bool {
	fi, err := os.Stat(f.path)
	if err != nil {
		log.Printf("[WARN] Cannot stat file: %v", err)
		return true
	}
	return fi.ModTime().After(since)
}

func (f *File) Dependencies() []Target {
	return nil
}
