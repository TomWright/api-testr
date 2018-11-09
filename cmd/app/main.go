package main

import (
	"github.com/tomwright/api-testr/testr"
	"github.com/tomwright/api-testr/testr/parse"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	baseAddr := "https://jsonplaceholder.typicode.com"
	testDir := "tests"

	testFiles, err := filepath.Glob(testDir + "/*.json")
	if err != nil {
		panic(err)
	}

	tests := make([]*testr.Test, 0)
	for _, testFile := range testFiles {
		t, err := parse.File(testFile, baseAddr)
		if err != nil {
			log.Printf("could not parse test file `%s`: %s", testFile, err)
			continue
		}
		tests = append(tests, t)
	}

	testr.RunAll(testr.RunAllArgs{
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
		MaxConcurrentTests: 5,
	}, tests...)
	log.Println("tests finished")
}
