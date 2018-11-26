package check

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// BodyJSONChecker is used to validate http response body can be JSON decoded and is equal to `Value`
type BodyJSONChecker struct {
	Value interface{}
}

// Check performs the BodyJSON check
func (c *BodyJSONChecker) Check(response *http.Response) error {
	body, err := readResponseBody(response)
	if err != nil {
		return fmt.Errorf("could not read response body: %s", err)
	}

	var got map[string]interface{}

	err = json.Unmarshal(body, &got)
	if err != nil {
		return fmt.Errorf("could not unmarshal actual response: %s", err)
	}

	if !reflect.DeepEqual(c.Value, got) {
		return fmt.Errorf("expected response body of `%v`, got `%v`", c.Value, got)
	}

	return nil
}
