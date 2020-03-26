package check

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Checker is used to outline each individual response check
type Checker interface {
	Check(ctx context.Context, response *http.Response) error
}

// readResponseBody allows you to read a http response body multiple times
func readResponseBody(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}
	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
