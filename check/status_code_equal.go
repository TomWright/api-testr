package check

import (
	"context"
	"fmt"
	"net/http"
)

// UnexpectedStatusCode is returned when a check fails.
type UnexpectedStatusCode struct {
	// Expected is the expected status code.
	Expected int
	// Actual is the actual status code.
	Actual int
}

// Error returns an error string.
func (e *UnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: expected %v, got %v", e.Expected, e.Actual)
}

// StatusCodeEqualChecker is used to validate the status code in the response
type StatusCodeEqualChecker struct {
	Value int
}

// Check performs the StatusCodeEqual check
func (c *StatusCodeEqualChecker) Check(ctx context.Context, response *http.Response) error {
	if response.StatusCode != c.Value {
		return &UnexpectedStatusCode{
			Expected: c.Value,
			Actual:   response.StatusCode,
		}
	}
	return nil
}
