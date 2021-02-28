package targets

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/avoronkov/go-target-once/lib/logger"
)

type Url struct {
	url        string
	method     string
	user, pass string
	dynamic    bool
	// TODO support multiple values
	headers map[string]string
}

var _ Target = (*Url)(nil)
var _ Modifiable = (*Url)(nil)

func NewUrl(url string) *Url {
	return &Url{
		url:     url,
		method:  http.MethodGet,
		headers: make(map[string]string),
	}
}

func (u *Url) SetBasicAuth(user, pass string) *Url {
	u.user = user
	u.pass = pass
	return u
}

func (u *Url) Dynamic() *Url {
	u.dynamic = true
	return u
}

func (u *Url) SetHeader(key, value string) *Url {
	u.headers[key] = value
	return u
}

func (u *Url) TargetID() string {
	return u.url
}

func (u *Url) Build(bc BuildContext) Result {
	req, err := http.NewRequest(u.method, u.url, nil)
	if err != nil {
		return Failed(err)
	}

	resp, err := http.DefaultClient.Do(req)
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
	if u.dynamic {
		return true
	}

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
