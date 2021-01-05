package main

import (
	"dont-repeat-twice/lib/builder"
	"dont-repeat-twice/lib/id"
	"dont-repeat-twice/lib/targets"
	"encoding/json"
	"fmt"
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
	path, _ := bc.GetDependency(0)

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
	meta["environment"] = environment()
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

/*
var metaTarget = &targets.Custom{
	Id:   "meta-target",
	Deps: []targets.Target{targets.NewFile("build.json")},
	DoBuild: func(this *targets.Custom, args ...id.Interface) (interface{}, time.Time) {
		now := time.Now()
		path, _ := this.Deps[0].Build()

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
		meta["environment"] = environment()
		data, err := json.Marshal(meta)
		if err != nil {
			log.Printf("[err] json.Marshal failed: %v", err)
			return "500", now
		}
		return string(data), now
	},
}
*/

var metaTarget = &BuildInfo{
	Common: targets.Common{
		Id: "buildinfo",
		Deps: []targets.Target{
			targets.NewFile("build.json"),
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
