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
	"strings"
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

func (b *BuildInfo) Build(bc targets.BuildContext, args ...id.Interface) (content interface{}, t time.Time, err error) {
	now := time.Now()
	//
	buildinfo, err := bc.GetDependency(0)
	if err != nil {
		return nil, time.Time{}, err
	}

	meta := make(map[string]interface{})
	if err := json.Unmarshal(buildinfo.([]byte), &meta); err != nil {
		log.Printf("[err] json.Decode failed: %v", err)
		return nil, time.Time{}, err
	}

	envContent, err := bc.GetDependency(1)
	if err != nil {
		return nil, time.Time{}, err
	}
	environment := strings.TrimSpace(string(envContent.([]byte)))
	meta["environment"] = string(environment)

	data, err := json.Marshal(meta)
	if err != nil {
		log.Printf("[err] json.Marshal failed: %v", err)
		return nil, time.Time{}, err
	}
	return string(data), now, nil
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
	data, err := builder.Build(metaTarget)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	fmt.Fprintf(w, "%v\n", data)
}

func main() {
	http.HandleFunc("/", http.HandlerFunc(metaHandler))
	log.Printf("Lintening to http://localhost:8080/ ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
