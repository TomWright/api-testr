package main

import (
	"github.com/tomwright/api-testr/testr"
	"github.com/tomwright/api-testr/testr/parse"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	baseAddr := "https://jsonplaceholder.typicode.com"
	testDir := "tests"
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	maxConcurrentTests := 1

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

	sem := make(chan struct{}, maxConcurrentTests)

	wg := &sync.WaitGroup{}
	wg.Add(len(tests))

	for _, t := range tests {
		go func(t *testr.Test) {
			defer func() {
				<-sem
				wg.Done()
			}()
			sem <- struct{}{}

			err = testr.Run(t, httpClient)
			if err != nil {
				log.Printf("test `%s` failed: %s", t.Name, err)
			}

			time.Sleep(time.Second)
		}(t)
	}

	wg.Wait()
	log.Println("tests finished")
}
