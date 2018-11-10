package check

import (
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"reflect"
)

type BodyJSONQueryEqualChecker struct {
	Query     string
	ValueType string
	Value     interface{}
	NullValue bool
}

func (c *BodyJSONQueryEqualChecker) Check(response *http.Response) error {
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

	return nil
}
