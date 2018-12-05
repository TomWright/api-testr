package main

import (
	"context"
	"flag"
	"github.com/tomwright/api-testr/testr"
	"github.com/tomwright/api-testr/testr/parse"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultMaxConcurrentTests = 5
	defaultHTTPTimeout        = 5
)

func main() {
	var baseAddr string
	var testDirs string
	var maxConcurrentTests int
	var httpTimeout int

	flag.StringVar(&baseAddr, "base", "", "the base address used in http requests")
	flag.StringVar(&testDirs, "tests", "", "the directory that tests are located in")
	flag.IntVar(&maxConcurrentTests, "maxConcurrentTests", defaultMaxConcurrentTests, "the maximum number of tests that can be run concurrently")
	flag.IntVar(&httpTimeout, "httpTimeout", defaultHTTPTimeout, "the http timeout duration in seconds")

	flag.Parse()

	logger := log.New(os.Stderr, "", log.LstdFlags)

	ctx := context.Background()
	ctx = testr.ContextWithBaseURL(ctx, baseAddr)

	tests := make([]*testr.Test, 0)

	for _, testDir := range strings.Split(testDirs, ",") {
		if logger != nil {
			logger.Printf("searching directory for tests: %s", testDir)
		}
		testFiles, err := filepath.Glob(testDir + "/[^_]*.json")
		if err != nil {
			panic(err)
		}
		for _, testFile := range testFiles {
			t, err := parse.File(ctx, testFile)
			if err != nil {
				if logger != nil {
					logger.Printf("could not parse test file `%s`: %s", testFile, err)
				}
				continue
			}
			tests = append(tests, t)
		}
	}

	res := testr.RunAll(ctx, testr.RunAllArgs{
		Logger: logger,
		HTTPClient: &http.Client{
			Timeout: time.Second * time.Duration(httpTimeout),
		},
		MaxConcurrentTests: maxConcurrentTests,
	}, tests...)

	if logger != nil {
		logger.Printf("tests finished\nexecuted: %d\npassed: %d\nfailed: %d", res.Executed, res.Passed, res.Failed)
	}

	if res.Failed > 0 {
		os.Exit(1)
	}

	os.Exit(0)
}
