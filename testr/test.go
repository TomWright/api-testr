package testr

import (
	"github.com/tomwright/api-testr/testr/check"
	"net/http"
)

type Test struct {
	Name     string
	Checks   []check.Checker
	Request  *http.Request
	Response *http.Response
}
