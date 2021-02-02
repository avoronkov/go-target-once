package targets

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
)

type Url struct {
	url string
}

var _ Target = (*Url)(nil)
var _ Modifiable = (*Url)(nil)

func NewUrl(url string) *Url {
	return &Url{
		url: url,
	}
}

func (u *Url) TargetID() string {
	return u.url
}

func (u *Url) Build(bc BuildContext) Result {
	resp, err := http.Get(u.url)
	if err != nil {
		return Failed(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Failed(err)
	}

	lastMdf, ok := lastModified(resp)
	if !ok {
		lastMdf = time.Now()
	}

	return OKTime(data, lastMdf)
}

var lastModifiedFmt = time.RFC1123

func (u *Url) IsModified(since time.Time) bool {
	resp, err := http.Head(u.url)
	if err != nil {
		logger.Warningf("HEAD %v failed: %v", u.url, err)
		return true
	}
	lastTime, ok := lastModified(resp)
	if !ok {
		return true
	}
	return lastTime.After(since)
}

func lastModified(resp *http.Response) (time.Time, bool) {
	last := resp.Header.Get("Last-Modified")
	if last == "" {
		logger.Warningf("Last-Modified is empty")
		return time.Time{}, false
	}
	lastTime, err := time.Parse(lastModifiedFmt, last)
	if err != nil {
		logger.Warningf("Cannot parse Last-Modified: %v (%v)", last, err)
		return time.Time{}, false
	}
	return lastTime, true
}
