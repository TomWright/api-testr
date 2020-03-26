package check

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
)

// JSONQueryValueMissingError is returned when a check fails.
type JSONQueryValueMissingError struct {
	// Query is the JSON query.
	Query string
}

// Error returns an error string.
func (e *JSONQueryValueMissingError) Error() string {
	return fmt.Sprintf("value at %v is missing", e.Query)
}

// BodyJSONQueryExistsChecker queries the http response body JSON using `Query` and ensures a value exists there
type BodyJSONQueryExistsChecker struct {
	Query  string
	DataID string
}

// Check performs the BodyJSONQueryExists check
func (c *BodyJSONQueryExistsChecker) Check(ctx context.Context, response *http.Response) error {
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

	return ContextWithOptionalDataID(ctx, c.DataID, r.Value())
}
