package testr

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	DefaultMaxConcurrentTests = 5
)

func Run(t *Test, httpClient *http.Client) error {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	log.Printf("Running test: %s\n", t.Name)

	resp, err := httpClient.Do(t.Request)
	if err != nil {
		return fmt.Errorf("could not execute request: %s", err)
	}

	for _, c := range t.Checks {
		err := c.Check(resp)
		if err != nil {
			return fmt.Errorf("failed `%T` check: %s", c, err)
		}
	}

	return nil
}

type RunAllArgs struct {
	HTTPClient         *http.Client
	MaxConcurrentTests int
}

type RunAllResult struct {
	Executed int
	Passed   int
	Failed   int
}

func RunAll(args RunAllArgs, tests ...*Test) RunAllResult {
	if args.HTTPClient == nil {
		args.HTTPClient = http.DefaultClient
	}
	if args.MaxConcurrentTests == 0 {
		args.MaxConcurrentTests = DefaultMaxConcurrentTests
	}

	sem := make(chan struct{}, args.MaxConcurrentTests)

	groupedTests := groupTests(tests...)

	resMu := sync.Mutex{}
	res := &RunAllResult{}

groupLoop:
	for groupName, groupTests := range groupedTests {
		log.Printf("running group %s\n", groupName)
		if len(groupTests.tests) == 0 {
			log.Println("no tests")
			continue groupLoop
		}

	groupOrderLoop:
		for i := groupTests.minOrder; i <= groupTests.maxOrder; i++ {
			log.Printf("running order %d\n", i)

			groupOrderTests, ok := groupTests.tests[i]
			if !ok {
				continue groupOrderLoop
			}

			wg := &sync.WaitGroup{}
			wg.Add(len(groupOrderTests))

			for _, t := range groupOrderTests {
				go func(t *Test) {
					defer func() {
						<-sem
						wg.Done()
					}()
					sem <- struct{}{}

					err := Run(t, args.HTTPClient)

					resMu.Lock()
					defer resMu.Unlock()

					res.Executed++
					if err != nil {
						res.Failed++
						log.Printf("test `%s` failed: %s\n", t.Name, err)
					} else {
						res.Passed++
					}
				}(t)
			}

			wg.Wait()
		}
	}

	resMu.Lock()
	defer resMu.Unlock()

	return *res
}

func newTestsGroup() *testsGroup {
	return &testsGroup{
		minOrder: -1,
		maxOrder: -1,
		tests:    make(map[int][]*Test, 0),
	}
}

type testsGroup struct {
	minOrder int
	maxOrder int
	tests    map[int][]*Test
}

func (g *testsGroup) add(t *Test) {
	if t.Order < g.minOrder || g.minOrder == -1 {
		g.minOrder = t.Order
	}
	if t.Order > g.maxOrder || g.maxOrder == -1 {
		g.maxOrder = t.Order
	}
	if _, ok := g.tests[t.Order]; !ok {
		g.tests[t.Order] = make([]*Test, 0)
	}
	g.tests[t.Order] = append(g.tests[t.Order], t)
}

func groupTests(tests ...*Test) map[string]*testsGroup {
	res := make(map[string]*testsGroup)

	for _, t := range tests {
		if _, ok := res[t.Group]; !ok {
			res[t.Group] = newTestsGroup()
		}
		res[t.Group].add(t)
	}

	return res
}
