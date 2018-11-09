package check

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Checker interface {
	Check(response *http.Response) error
}

// readResponseBody allows you to read a http response body multiple times
func readResponseBody(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %s", err)
	}
	response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %s", err)
	}
	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
