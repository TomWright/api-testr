package check

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
)

// BodyJSONQueryRegexMatchChecker queries the http response body JSON using `Query` and ensures that it matches the regex pattern in `Regexp`
type BodyJSONQueryRegexMatchChecker struct {
	Query   string
	Regexp  *regexp.Regexp
	DataIDs map[int]string
}

// Check performs the BodyJSONQueryRegexMatch check
func (c *BodyJSONQueryRegexMatchChecker) Check(ctx context.Context, response *http.Response) error {
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

	if c.DataIDs != nil {
		values := c.Regexp.FindStringSubmatch(r.String())
		for i, dataID := range c.DataIDs {
			if len(values) > i {
				if err := ContextWithOptionalDataID(ctx, dataID, values[i]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
