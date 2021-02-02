package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"time"

	"github.com/avoronkov/go-target-once/lib/builder"
	"github.com/avoronkov/go-target-once/lib/targets"
)

// Declare the Gopher target
type GopherTarget struct{}

var _ targets.Target = (*GopherTarget)(nil)
var _ targets.WithDependencies = (*GopherTarget)(nil)
var _ targets.Cacheable = (*GopherTarget)(nil)

// Declare TargetID: it is used as a key in cache storage.
func (g *GopherTarget) TargetID() string {
	return "gopher-target"
}

func (g *GopherTarget) Cacheable() bool {
	return true
}

// Declare the dependencies: URL to the gopher image
func (g *GopherTarget) Dependencies() map[string]targets.Target {
	return map[string]targets.Target{
		"gopher": targets.NewUrl("https://blog.golang.org/gopher/gopher.png"),
	}
}

// Build action.
func (g *GopherTarget) Build(bc targets.BuildContext) (interface{}, time.Time, error) {
	fmt.Println("Calculating gopher image size...")

	// Dependency content is accessible by the name "gopher"
	gopherData, err := bc.GetDependency("gopher")
	if err != nil {
		return nil, time.Time{}, err
	}
	img, _, err := image.Decode(bytes.NewReader(gopherData.([]byte)))
	if err != nil {
		return nil, time.Time{}, err
	}
	bounds := img.Bounds()
	res := fmt.Sprintf("Gopher image size: %v x %v", bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y)
	return res, time.Now(), nil
}

func main() {
	target := new(GopherTarget)

	// Build target the first time
	result, tm, err := builder.Build(target)
	fmt.Printf("[1] %v (time=%v, err=%v)\n", result, tm, err)

	// Rebuild the target.
	// GopherTarget.Build will not be called if image on server is unchanged.
	result, tm, err = builder.Build(target)
	fmt.Printf("[2] %v (time=%v, err=%v)\n", result, tm, err)
}
