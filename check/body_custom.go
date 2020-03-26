package check

import (
	"context"
	"errors"
	"net/http"
)

var (
	// ErrMissingCheckFunc is returned when a BodyCustomChecker is checked without a valid func.
	ErrMissingCheckFunc = errors.New("missing check func")
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
		return err
	}

	if c.CheckBody == nil {
		return ErrMissingCheckFunc
	}

	return c.CheckBody(body)
}
