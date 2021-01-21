package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/avoronkov/go-target-once/lib/builder"
	"github.com/avoronkov/go-target-once/lib/targets"
)

type AsciiImage struct {
	url    string
	target targets.Target
}

var _ targets.Target = (*AsciiImage)(nil)
var _ targets.WithDependencies = (*AsciiImage)(nil)

func NewAsciiImage(url string) *AsciiImage {
	return &AsciiImage{
		url:    url,
		target: targets.NewUrl(url),
	}
}

func (g *AsciiImage) TargetId() string {
	return fmt.Sprintf("web-resource-%v", g.url)
}

func (g *AsciiImage) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"image": g.target,
	}
}

func (g *AsciiImage) Build(bc targets.BuildContext) (content interface{}, t time.Time, err error) {
	now := time.Now()
	data, err := bc.GetDependency("image")
	if err != nil {
		return nil, time.Time{}, err
	}
	ascii, err := image2ascii(data.([]byte))
	if err != nil {
		return nil, time.Time{}, err
	}
	return []byte(ascii), now, nil
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
	data, _, err := builder.Build(resource)
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
