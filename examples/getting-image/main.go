package main

import (
	"dont-repeat-twice/lib/builder"
	"dont-repeat-twice/lib/targets"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type AsciiImage struct {
	url    string
	target targets.Target
}

var _ targets.Target = (*AsciiImage)(nil)

func NewAsciiImage(url string) *AsciiImage {
	return &AsciiImage{
		url:    url,
		target: targets.NewUrl(url),
	}
}

func (g *AsciiImage) TargetId() string {
	return fmt.Sprintf("web-resource-%v", g.url)
}

func (g *AsciiImage) Dependencies() []targets.Target {
	return []targets.Target{
		g.target,
	}
}

func (g *AsciiImage) Build(bc targets.BuildContext) (content interface{}, t time.Time, err error) {
	now := time.Now()
	data, err := bc.GetDependency(0)
	if err != nil {
		return nil, time.Time{}, err
	}
	ascii, err := image2ascii(data.([]byte))
	if err != nil {
		return nil, time.Time{}, err
	}
	return []byte(ascii), now, nil
}

func (g *AsciiImage) IsModified(since time.Time) bool {
	return g.target.IsModified(since)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	httpsPrefix := "https:/"
	if strings.HasPrefix(path, httpsPrefix) {
		path = "https://" + path[len(httpsPrefix):]
	} else {
		path = "https://" + path
	}
	log.Printf("path (2) = %v", path)

	resource := NewAsciiImage(path)
	data, err := builder.Build(resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data.([]byte))
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("Lintening to http://localhost:8080/ ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
