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

	// if c.NullValue {
	// 	if r.Value() == nil {
	// 		return nil
	// 	}
	// 	return fmt.Errorf("json query element at `%v` does not match expected: %v", c.Query, nil)
	// }
	//
	// switch c.ValueType {
	// case "string":
	// 	err = c.checkString(fmt.Sprint(c.Value), r)
	//
	// case "int":
	// 	switch expect := c.Value.(type) {
	// 	case string:
	// 		i, _ := strconv.Atoi(expect)
	// 		err = c.checkInt(int64(i), r)
	// 	case []byte:
	// 		i, _ := strconv.Atoi(string(expect))
	// 		err = c.checkInt(int64(i), r)
	// 	case int:
	// 		err = c.checkInt(int64(expect), r)
	// 	case int32:
	// 		err = c.checkInt(int64(expect), r)
	// 	case int64:
	// 		err = c.checkInt(expect, r)
	// 	case float32:
	// 		err = c.checkInt(int64(expect), r)
	// 	case float64:
	// 		err = c.checkInt(int64(expect), r)
	// 	}
	//
	// case "float":
	// 	switch expect := c.Value.(type) {
	// 	case string:
	// 		i, _ := strconv.ParseFloat(expect, 64)
	// 		err = c.checkFloat(i, r)
	// 	case []byte:
	// 		i, _ := strconv.ParseFloat(string(expect), 64)
	// 		err = c.checkFloat(i, r)
	// 	case int:
	// 		err = c.checkFloat(float64(expect), r)
	// 	case int32:
	// 		err = c.checkFloat(float64(expect), r)
	// 	case int64:
	// 		err = c.checkFloat(float64(expect), r)
	// 	case float32:
	// 		err = c.checkFloat(float64(expect), r)
	// 	case float64:
	// 		err = c.checkFloat(expect, r)
	// 	}
	//
	// default:
	// 	switch expect := c.Value.(type) {
	// 	case string:
	// 		err = c.checkString(expect, r)
	// 	case []byte:
	// 		err = c.checkString(string(expect), r)
	// 	case int:
	// 		err = c.checkInt(int64(expect), r)
	// 	case int32:
	// 		err = c.checkInt(int64(expect), r)
	// 	case int64:
	// 		err = c.checkInt(expect, r)
	// 	case float32:
	// 		err = c.checkFloat(float64(expect), r)
	// 	case float64:
	// 		err = c.checkFloat(float64(expect), r)
	// 	}
	// }
	//
	// if err != nil {
	// 	return fmt.Errorf("json query %s element at `%s` does not match expected: %s", c.ValueType, c.Query, err)
	// }

	return nil
}

// func (c *BodyJSONQueryEqualChecker) checkString(exp string, gotResult gjson.Result) error {
// 	if got := gotResult.String(); exp != got {
// 		return fmt.Errorf("expected string `%s`, got `%s`", exp, got)
// 	}
// 	return nil
// }
//
// func (c *BodyJSONQueryEqualChecker) checkInt(exp int64, gotResult gjson.Result) error {
// 	if got := gotResult.Int(); exp != got {
// 		return fmt.Errorf("expected int `%d`, got `%d`", exp, got)
// 	}
// 	return nil
// }
//
// func (c *BodyJSONQueryEqualChecker) checkFloat(exp float64, gotResult gjson.Result) error {
// 	if got := gotResult.Float(); exp != got {
// 		return fmt.Errorf("expected int `%v`, got `%v`", exp, got)
// 	}
// 	return nil
// }
