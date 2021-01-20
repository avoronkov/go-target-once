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


## Debugging

Build with `go build -tags=debug` to enable debug logging.

## TODO

- Cache intermediate results in BuildContext.
- Improve resolving "diamond" dependencies.
- Build methods should return time.Time.
- Add methods "IsModified" and "BuildIfModified" to supports handling "If-Modified-Since" HTTP header.
