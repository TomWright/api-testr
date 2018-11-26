package check

import (
	"fmt"
	"net/http"
)

// StatusCodeEqualChecker is used to validate the status code in the response
type StatusCodeEqualChecker struct {
	Value int
}

// Check performs the StatusCodeEqual check
func (c *StatusCodeEqualChecker) Check(response *http.Response) error {
	if response.StatusCode != c.Value {
		return fmt.Errorf("expected status code `%d`, got `%d`", c.Value, response.StatusCode)
	}
	return nil
}
