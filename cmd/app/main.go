package main

import (
	"context"
	"fmt"
	"github.com/tomwright/api-testr/testr"
	"github.com/tomwright/api-testr/testr/check"
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

	ctx := context.Background()
	ctx = testr.ContextWithBaseURL(ctx, baseAddr)

	var custom123 check.BodyCustomCheckerFunc = func(bytes []byte) error {
		if string(bytes) == "" {
			return fmt.Errorf("response is empty")
		}
		return nil
	}
	ctx = testr.ContextWithCustomBodyCheck(ctx, "123check", custom123)

	tests := make([]*testr.Test, 0)
	for _, testFile := range testFiles {
		t, err := parse.File(ctx, testFile)
		if err != nil {
			log.Printf("could not parse test file `%s`: %s", testFile, err)
			continue
		}
		tests = append(tests, t)
	}

	res := testr.RunAll(testr.RunAllArgs{
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
		MaxConcurrentTests: 5,
	}, tests...)

	log.Printf("tests finished\n\texecuted: %d\n\tpassed: %d\n\tfailed: %d", res.Executed, res.Passed, res.Failed)
}
