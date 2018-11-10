package check

import (
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
)

type BodyJSONQueryRegexMatchChecker struct {
	Query  string
	Regexp *regexp.Regexp
}

func (c *BodyJSONQueryRegexMatchChecker) Check(response *http.Response) error {
	body, err := readResponseBody(response)
	if err != nil {
		return fmt.Errorf("could not read response body: %s", err)
	}

	j := gjson.ParseBytes(body)

	r := j.Get(c.Query)

	if !r.Exists() {
		return fmt.Errorf("json query element does not exist: %s`", c.Query)
	}

	if !c.Regexp.MatchString(r.String()) {
		return fmt.Errorf("json query element at `%s` with a value of `%s` does not match regex pattern `%s`", c.Query, r.String(), c.Regexp.String())
	}

	return nil
}
