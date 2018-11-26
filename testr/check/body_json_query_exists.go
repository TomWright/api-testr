package check

import (
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
)

// BodyJSONQueryExistsChecker queries the http response body JSON using `Query` and ensures a value exists there
type BodyJSONQueryExistsChecker struct {
	Query string
}

// Check performs the BodyJSONQueryExists check
func (c *BodyJSONQueryExistsChecker) Check(response *http.Response) error {
	body, err := readResponseBody(response)
	if err != nil {
		return fmt.Errorf("could not read response body: %s", err)
	}

	j := gjson.ParseBytes(body)

	r := j.Get(c.Query)

	if !r.Exists() {
		return fmt.Errorf("json query element does not exist: %s`", c.Query)
	}

	return nil
}
