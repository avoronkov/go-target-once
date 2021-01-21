package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/avoronkov/go-target-once/lib/builder"
	"github.com/avoronkov/go-target-once/lib/targets"
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

func (b *BuildInfo) Build(bc targets.BuildContext) (content interface{}, t time.Time, err error) {
	now := time.Now()
	//
	buildinfo, err := bc.GetDependency("build.json")
	if err != nil {
		return nil, time.Time{}, err
	}

	meta := make(map[string]interface{})
	if err := json.Unmarshal(buildinfo.([]byte), &meta); err != nil {
		log.Printf("[err] json.Decode failed: %v", err)
		return nil, time.Time{}, err
	}

	envContent, err := bc.GetDependency("env.txt")
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
		Id: "buildinfo-handler",
		Deps: map[string]targets.Target{
			"build.json": targets.NewFile("build.json"),
			"env.txt":    targets.NewFile("env.txt"),
		},
	},
}

func metaHandler(w http.ResponseWriter, r *http.Request) {
	data, _, err := builder.Build(metaTarget)
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
