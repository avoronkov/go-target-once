package main

import (
	"dont-repeat-twice/lib/builder"
	"dont-repeat-twice/lib/id"
	"dont-repeat-twice/lib/targets"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func environment() string {
	e, ok := os.LookupEnv("APP_ENV")
	if ok {
		return e
	}
	return "development"
}

type BuildInfo struct {
	targets.Common
}

func (b *BuildInfo) Build(bc targets.BuildContext, args ...id.Interface) (content interface{}, t time.Time) {
	now := time.Now()
	path := bc.GetDependency(0)

	f, err := os.Open(path.(string))
	if err != nil {
		log.Printf("[err] os.Open failed: %v", err)
		return "404", now
	}
	defer f.Close()

	meta := make(map[string]interface{})
	if err := json.NewDecoder(f).Decode(&meta); err != nil {
		log.Printf("[err] json.Decode failed: %v", err)
		return "500", now
	}

	envPath := bc.GetDependency(1)
	environment, _ := ioutil.ReadFile(envPath.(string))
	meta["environment"] = string(environment)

	data, err := json.Marshal(meta)
	if err != nil {
		log.Printf("[err] json.Marshal failed: %v", err)
		return "500", now
	}
	return string(data), now
}

func (b *BuildInfo) IsModified(since time.Time) bool {
	return b.Common.DepsModified(since)
}

var metaTarget = &BuildInfo{
	Common: targets.Common{
		Id: "buildinfo",
		Deps: []targets.Target{
			targets.NewFile("build.json"),
			targets.NewFile("env.txt"),
		},
	},
}

func metaHandler(w http.ResponseWriter, r *http.Request) {
	data := builder.Build(metaTarget)
	fmt.Fprintf(w, "%v\n", data)
}

func main() {
	http.HandleFunc("/", http.HandlerFunc(metaHandler))
	log.Printf("Lintening to http://localhost:8080/ ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
