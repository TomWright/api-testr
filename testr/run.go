package testr

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	DefaultMaxConcurrentTests = 5
)

func Run(t *Test, httpClient *http.Client, logger *log.Logger) error {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if logger != nil {
		logger.Printf("running test: %s\n", t.Name)
	}

	var err error

	t.Response, err = httpClient.Do(t.Request)
	if err != nil {
		return fmt.Errorf("could not execute request: %s", err)
	}

	for _, c := range t.Checks {
		err := c.Check(t.Response)
		if err != nil {
			return fmt.Errorf("failed `%T` check: %s", c, err)
		}
	}

	return nil
}

type RunAllArgs struct {
	HTTPClient         *http.Client
	MaxConcurrentTests int
	Logger             *log.Logger
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
		if args.Logger != nil {
			args.Logger.Printf("running group %s\n", groupName)
		}
		if len(groupTests.tests) == 0 {
			if args.Logger != nil {
				args.Logger.Println("no tests")
			}
			continue groupLoop
		}

	groupOrderLoop:
		for i := groupTests.minOrder; i <= groupTests.maxOrder; i++ {
			if args.Logger != nil {
				args.Logger.Printf("running order %d\n", i)
			}

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

					err := Run(t, args.HTTPClient, args.Logger)

					resMu.Lock()
					defer resMu.Unlock()

					res.Executed++
					if err != nil {
						res.Failed++
						if args.Logger != nil {
							args.Logger.Printf("test `%s` failed: %s\nRequest:\n%s\nResponse:\n%s\n", t.Name, err, fmtRequest(t.Request), fmtResponse(t.Response))
						}
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

func fmtRequest(r *http.Request) string {
	reqBodyStr, _ := ioutil.ReadAll(r.Body)
	requestInfo := fmt.Sprintf("%s %s", r.Method, r.URL.String())
	if reqBodyStr != nil && len(reqBodyStr) > 0 {
		requestInfo = fmt.Sprintf("%s\nBody:", requestInfo)
		requestInfo = fmt.Sprintf("%s\n%s", requestInfo, string(reqBodyStr))
	}
	if r.Header != nil && len(r.Header) > 0 {
		requestInfo = fmt.Sprintf("%s%s", requestInfo, fmtHeader(r.Header))
	}
	return requestInfo
}

func fmtResponse(r *http.Response) string {
	respBodyStr, _ := ioutil.ReadAll(r.Body)
	responseInfo := fmt.Sprintf("%s", r.Status)
	if respBodyStr != nil && len(respBodyStr) > 0 {
		responseInfo = fmt.Sprintf("%s\nBody:", responseInfo)
		responseInfo = fmt.Sprintf("%s\n%s", responseInfo, string(respBodyStr))
	}
	if r.Header != nil && len(r.Header) > 0 {
		responseInfo = fmt.Sprintf("%s%s", responseInfo, fmtHeader(r.Header))
	}
	return responseInfo
}

func fmtHeader(header http.Header) string {
	resp := ""
	if header != nil && len(header) > 0 {
		resp = fmt.Sprintf("%s\nHeaders:", resp)
		for headerName, headers := range header {
			for _, headerVal := range headers {
				resp = fmt.Sprintf("%s\n\t%s: %s", resp, headerName, headerVal)
			}
		}
	}
	return resp
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
