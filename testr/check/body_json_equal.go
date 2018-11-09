package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

type BodyJSONChecker struct {
	Value interface{}
}

func (c *BodyJSONChecker) Check(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)
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
