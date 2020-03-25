package check

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

// DataEqualChecker is used to check whether or not the value stored in the context data is equal to the given value
type DataEqualChecker struct {
	DataID string
	Value  interface{}
}

// Check performs the DataEqual check
func (c *DataEqualChecker) Check(ctx context.Context, response *http.Response) error {
	if exp, got := c.Value, DataIDFromContext(ctx, c.DataID); !reflect.DeepEqual(exp, got) {
		return fmt.Errorf("expected data item `%s` to be `%v`, got `%v`", c.DataID, exp, got)
	}

	return nil
}
