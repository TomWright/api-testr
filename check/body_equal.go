package check

import (
	"context"
	"fmt"
	"net/http"
)

// UnexpectedValueError is returned when a check fails.
type UnexpectedValueError struct {
	// Expected is the expected value.
	Expected string
	// Actual is the actual value.
	Actual string
}

// Error returns an error string.
func (e *UnexpectedValueError) Error() string {
	return fmt.Sprintf("unexpected value: expected %v, got %v", e.Expected, e.Actual)
}

// BodyEqualChecker is used to validate the http response body string exactly matches `Value`
type BodyEqualChecker struct {
	Value string
}

// Check performs the BodyEqual check
func (c *BodyEqualChecker) Check(ctx context.Context, response *http.Response) error {
	body, err := readResponseBody(response)
	if err != nil {
		return err
	}

	if exp, got := c.Value, string(body); exp != got {
		return &UnexpectedValueError{
			Expected: exp,
			Actual:   got,
		}
	}

	return nil
}
