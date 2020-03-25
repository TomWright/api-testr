package check

import (
	"context"
	"fmt"
	"net/http"
)

// BodyCustomCheckerFunc defines the function used to perform a custom http response body check
type BodyCustomCheckerFunc func([]byte) error

// BodyCustomChecker is used to run a BodyCustomCheckerFunc against the body bytes in the http response
type BodyCustomChecker struct {
	CheckBody BodyCustomCheckerFunc
}

// Check performs the BodyCustom check
func (c *BodyCustomChecker) Check(ctx context.Context, response *http.Response) error {
	body, err := readResponseBody(response)
	if err != nil {
		return fmt.Errorf("could not read response body: %s", err)
	}

	if c.CheckBody == nil {
		return fmt.Errorf("missing check func on body custom checker")
	}

	return c.CheckBody(body)
}
