# go-target-once

## Disclaimer

This is not my origin idea - I read it somewhere on the internet but unfortunately have lost the link :(
(If you know what article I'm talking about please let me know).

## Idea

Imagine you wrote a web application.

At first some user performs a request to your application then application starts to work:
it performs a number of requests to database, to some external services, reads some files,
then it processes/combines the fetched data and forms a response for the user (e.g. HTML page or JSON data).
Performing request takes some time.

Then user hits F5, performs the same request, waits for the same time and receives the response.
Here the user may wonder: Why?
Why the processing takes the same amount of time if no data has changed and all required information is already processed?
Why not to return previously calculated result? It should be really much faster.

This idea is rarely implemented in web applications but it is widely in another domain area - in build systems.
Build systems (such as make, buck etc) use the conception of targets and dependencies to describe the build process.
And they do not rebuild target if the target itself and its dependencies were not changed.

This library allows to use this conception in Go programs.
You just describe some calculations (e.g. http handlers) as a set of targets with dependencies
and then use "lib/builder" package to "Build" them.

Features:

- Parallel build of independent targets.
- Optional caching of build results.
- Cached results are not rebuilt if the target or its dependencies were not changed.

## Examples

Example of how to fetch image by URL and calculate its size using Targets.

```Go
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
type GopherTarget struct {}

// Make sure that required methods are implemented.
var _ targets.Target = (*GopherTarget)(nil)
var _ targets.WithDependencies = (*GopherTarget)(nil)
var _ targets.Cachable = (*GopherTarget)(nil)

// Declare TargetId: it is used as a key in cache storage.
func (g *GopherTarget) TargetId() string {
	return "gopher-target"
}

func (g *GopherTarget) Cachable() bool {
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
```

See the [examples](examples/) for more.

## Debugging

Build with `go build -tags=debug` to enable debug logging.

## TODO

- Cache intermediate results in BuildContext.
- Improve resolving "diamond" dependencies.
- Build methods should return time.Time.
- Add methods "IsModified" and "BuildIfModified" to support handling "If-Modified-Since" HTTP header.
