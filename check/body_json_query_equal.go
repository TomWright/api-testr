package check

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"reflect"
)

// UnexpectedJSONQueryValueError is returned when a check fails.
type UnexpectedJSONQueryValueError struct {
	// Query is the JSON query.
	Query string
	// Expected is the expected value.
	Expected interface{}
	// Actual is the actual value.
	Actual interface{}
}

// Error returns an error string.
func (e *UnexpectedJSONQueryValueError) Error() string {
	return fmt.Sprintf("unexpected value at %v: expected %v, got %v", e.Query, e.Expected, e.Actual)
}

// BodyJSONQueryEqualChecker queries the http response body JSON using `Query` and ensures the value is equal to `Value`
type BodyJSONQueryEqualChecker struct {
	Query     string
	Value     interface{}
	NullValue bool
	DataID    string
}

// Check performs the BodyJSONQueryEqual check
func (c *BodyJSONQueryEqualChecker) Check(ctx context.Context, response *http.Response) error {
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

	if got := r.Value(); !reflect.DeepEqual(c.Value, got) {
		return &UnexpectedJSONQueryValueError{
			Query:    c.Query,
			Expected: c.Value,
			Actual:   got,
		}
	}

	return ContextWithOptionalDataID(ctx, c.DataID, r.Value())
}
