package check

import (
	"fmt"
	"net/http"
)

type StatusCodeEqualChecker struct {
	Value int
}

func (c *StatusCodeEqualChecker) Check(response *http.Response) error {
	if response.StatusCode != c.Value {
		return fmt.Errorf("expected status code `%d`, got `%d`", c.Value, response.StatusCode)
	}
	return nil
}
