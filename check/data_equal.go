package check

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

// UnexpectedDataValueError is returned when a check fails.
type UnexpectedDataValueError struct {
	// DataID is the ID of the data.
	DataID string
	// Expected is the expected value.
	Expected interface{}
	// Actual is the actual value.
	Actual interface{}
}

// Error returns an error string.
func (e *UnexpectedDataValueError) Error() string {
	return fmt.Sprintf("unexpected data value for %s: expected %v, got %v", e.DataID, e.Expected, e.Actual)
}

// DataEqualChecker is used to check whether or not the value stored in the context data is equal to the given value
type DataEqualChecker struct {
	DataID string
	Value  interface{}
}

// Check performs the DataEqual check
func (c *DataEqualChecker) Check(ctx context.Context, response *http.Response) error {
	if exp, got := c.Value, DataIDFromContext(ctx, c.DataID); !reflect.DeepEqual(exp, got) {
		return &UnexpectedDataValueError{
			DataID:   c.DataID,
			Expected: exp,
			Actual:   got,
		}
	}

	return nil
}
