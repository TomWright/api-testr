package check

import "net/http"

type Checker interface {
	Check(response *http.Response) error
}
