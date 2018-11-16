package check

import (
	"fmt"
	"net/http"
)

type BodyCustomCheckerFunc func([]byte) error

type BodyCustomChecker struct {
	CheckBody BodyCustomCheckerFunc
}

func (c *BodyCustomChecker) Check(response *http.Response) error {
	body, err := readResponseBody(response)
	if err != nil {
		return fmt.Errorf("could not read response body: %s", err)
	}

	if c.CheckBody == nil {
		return fmt.Errorf("missing check func on body custom checker")
	}

	return c.CheckBody(body)
}
