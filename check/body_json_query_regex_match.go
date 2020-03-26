package check

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
)

// UnexpectedJSONQueryRegexValueError is returned when a check fails.
type UnexpectedJSONQueryRegexValueError struct {
	// Pattern is the regex pattern.
	Pattern string
	// Query is the JSON query.
	Query string
	// Actual is the actual value.
	Actual string
}

// Error returns an error string.
func (e *UnexpectedJSONQueryRegexValueError) Error() string {
	return fmt.Sprintf("unexpected value at %v: does not match pattern %v: got %v", e.Query, e.Pattern, e.Actual)
}

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
		return err
	}

	j := gjson.ParseBytes(body)

	r := j.Get(c.Query)

	if !r.Exists() {
		return &JSONQueryValueMissingError{
			Query: c.Query,
		}
	}

	if !c.Regexp.MatchString(r.String()) {
		return &UnexpectedJSONQueryRegexValueError{
			Pattern: c.Regexp.String(),
			Query:   c.Query,
			Actual:  r.String(),
		}
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
