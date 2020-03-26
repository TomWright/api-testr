package apitestr

import (
	"context"
	"fmt"
	"github.com/tomwright/apitestr/check"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	DefaultMaxConcurrentTests = 5
)

// Run executes a single test
func Run(ctx context.Context, t *Test, httpClient *http.Client, logger *log.Logger) error {
	testData := check.DataFromContext(ctx)
	if testData == nil {
		testData = make(map[string]interface{})
		ctx = check.ContextWithData(ctx, testData)
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if logger != nil {
		logger.Printf("running test: %s\n", t.Name)
	}

	var err error

	for i, initFunc := range t.RequestInitFuncs {
		initFuncData := t.RequestInitFuncsData[i]
		t.Request, err = initFunc(ctx, t.Request, initFuncData)
		if err != nil {
			return fmt.Errorf("request init func failed: %w", err)
		}
	}

	t.Response, err = httpClient.Do(t.Request)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}

	for _, c := range t.Checks {
		err := c.Check(ctx, t.Response)
		if err != nil {
			return fmt.Errorf("failed `%T` check: %w", c, err)
		}
	}

	return nil
}

// RunAllArgs defines which arguments are available to give to RunAll
type RunAllArgs struct {
	HTTPClient         *http.Client
	MaxConcurrentTests int
	Logger             *log.Logger
	// If Groups is not nil, only tests belonging to one of these groups will be executed
	Groups []string
	// If IgnoreGroups is not nil, any tests belonging to one of these groups will not be executed
	IgnoreGroups []string
	// IgnoreAllOnFailure should be true if when a test fails you want no more tests to be executed
	IgnoreAllOnFailure bool
	// IgnoreGroupOnFailure should be true if when a test fails you want no more tests in the failed group to be executed
	IgnoreGroupOnFailure bool
}

// RunAllResult is the response given from RunAll
type RunAllResult struct {
	Executed int
	Passed   int
	Failed   int
	Skipped  int
}

// RunAll runs the set of given tests
func RunAll(ctx context.Context, args RunAllArgs, tests ...*Test) RunAllResult {
	if args.HTTPClient == nil {
		args.HTTPClient = http.DefaultClient
	}
	if args.MaxConcurrentTests == 0 {
		args.MaxConcurrentTests = DefaultMaxConcurrentTests
	}

	testData := check.DataFromContext(ctx)
	if testData == nil {
		testData = make(map[string]interface{})
		ctx = check.ContextWithData(ctx, testData)
	}

	sem := make(chan struct{}, args.MaxConcurrentTests)

	groupedTests := groupTests(tests...)

	if args.Groups != nil {
		for groupName := range groupedTests {
			found := false
		foundLoop:
			for _, g := range args.Groups {
				if groupName == g {
					found = true
					break foundLoop
				}
			}
			if !found {
				delete(groupedTests, groupName)
			}
		}
	}

	if args.IgnoreGroups != nil {
		for _, g := range args.IgnoreGroups {
			delete(groupedTests, g)
		}
	}

	overallResMu := sync.Mutex{}
	overallRes := &RunAllResult{}

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

		groupResMu := sync.Mutex{}
		groupRes := &RunAllResult{}

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

					skip := false
					if args.IgnoreGroupOnFailure {
						groupResMu.Lock()
						if groupRes.Failed > 0 {
							skip = true
						}
						groupResMu.Unlock()
					}
					if !skip && args.IgnoreAllOnFailure {
						overallResMu.Lock()
						groupResMu.Lock()
						if overallRes.Failed > 0 || groupRes.Failed > 0 {
							skip = true
						}
						overallResMu.Unlock()
						groupResMu.Unlock()
					}
					if skip {
						groupResMu.Lock()
						groupRes.Skipped++
						if args.Logger != nil {
							args.Logger.Printf("test `%s` skipped", t.Name)
						}
						groupResMu.Unlock()
						return
					}

					err := Run(ctx, t, args.HTTPClient, args.Logger)

					groupResMu.Lock()
					defer groupResMu.Unlock()

					groupRes.Executed++
					if err != nil {
						groupRes.Failed++
						if args.Logger != nil {
							args.Logger.Printf("test `%s` failed: %s\nRequest:\n%s\nResponse:\n%s\n", t.Name, err, fmtRequest(t.Request), fmtResponse(t.Response))
						}
					} else {
						groupRes.Passed++
					}
				}(t)
			}

			wg.Wait()
		}

		overallResMu.Lock()
		groupResMu.Lock()

		overallRes.Executed += groupRes.Executed
		overallRes.Passed += groupRes.Passed
		overallRes.Failed += groupRes.Failed
		overallRes.Skipped += groupRes.Skipped

		overallResMu.Unlock()
		groupResMu.Unlock()
	}

	overallResMu.Lock()
	defer overallResMu.Unlock()

	return *overallRes
}

func fmtRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
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
	if r == nil {
		return ""
	}
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
