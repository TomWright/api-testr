package testr

import (
	"github.com/tomwright/api-testr/testr/check"
	"net/http"
)

type Test struct {
	// Name is the name of the test
	Name string
	// Group is the group that the test belongs to
	Group string
	// Order specified the order in which it will be run. Tests with the same order will be executed at the same time
	Order int
	// Checks contains all checks contained in this test
	Checks []check.Checker
	// Request contains the http request being made
	Request *http.Request
	// Response contains the http response
	Response *http.Response
}
