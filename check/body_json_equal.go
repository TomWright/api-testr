package check

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// UnexpectedJSONBodyError is returned when a check fails.
type UnexpectedJSONBodyError struct {
	// Expected is the expected value.
	Expected interface{}
	// Actual is the actual value.
	Actual interface{}
}

// Error returns an error string.
func (e *UnexpectedJSONBodyError) Error() string {
	return fmt.Sprintf("unexpected value: expected %v, got %v", e.Expected, e.Actual)
}

// BodyJSONChecker is used to validate http response body can be JSON decoded and is equal to `Value`
type BodyJSONChecker struct {
	Value interface{}
}

// Check performs the BodyJSON check
func (c *BodyJSONChecker) Check(ctx context.Context, response *http.Response) error {
	body, err := readResponseBody(response)
	if err != nil {
		return err
	}

	var got map[string]interface{}

	err = json.Unmarshal(body, &got)
	if err != nil {
		return fmt.Errorf("could not unmarshal actual response: %w", err)
	}

	if !reflect.DeepEqual(c.Value, got) {
		return &UnexpectedJSONBodyError{
			Expected: c.Value,
			Actual:   got,
		}
	}

	return nil
}
