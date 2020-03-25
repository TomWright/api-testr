package check

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"reflect"
)

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
		return fmt.Errorf("could not read response body: %s", err)
	}

	j := gjson.ParseBytes(body)

	r := j.Get(c.Query)

	if !r.Exists() {
		return fmt.Errorf("json query element does not exist: %s`", c.Query)
	}

	if got := r.Value(); !reflect.DeepEqual(c.Value, got) {
		return fmt.Errorf("json query element at `%s` does not match expected value. expected (%T)`%v`, got (%T)`%v`", c.Query, c.Value, c.Value, got, got)
	}

	return ContextWithOptionalDataID(ctx, c.DataID, r.Value())
}
