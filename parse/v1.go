package parse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tomwright/apitestr"
	"github.com/tomwright/apitestr/check"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type v1 struct {
	Name    string    `json:"name"`
	Group   string    `json:"group"`
	Order   int       `json:"order"`
	Request v1Request `json:"request"`
	Checks  []v1Check `json:"checks"`
}

type v1Request struct {
	Base     string                            `json:"base"`
	Method   string                            `json:"method"`
	Path     string                            `json:"path"`
	Body     interface{}                       `json:"body"`
	Headers  map[string]string                 `json:"headers"`
	InitFunc map[string]map[string]interface{} `json:"init"`
}

type v1Check struct {
	Type string `json:"type"`
	Data *data  `json:"data"`
}

func V1(ctx context.Context, data []byte) (*apitestr.Test, error) {
	v := v1{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("could not unmarshal v1 test data: %w", err)
	}

	if v.Request.Base == "" {
		v.Request.Base = apitestr.BaseURLFromContext(ctx)
	}

	var requestBody []byte

	if contentType, found := v.Request.Headers["Content-Type"]; found && strings.Contains(contentType, "application/json") {
		var err error
		requestBody, err = json.Marshal(v.Request.Body)
		if err != nil {
			return nil, fmt.Errorf("could not marshal request body: %w", err)
		}
	} else {
		switch requestBodyVal := v.Request.Body.(type) {
		case nil:
			break
		case string:
			requestBody = []byte(requestBodyVal)
		case []byte:
			requestBody = requestBodyVal
		default:
			return nil, fmt.Errorf("cannot handle type `%T` for body", v.Request.Body)
		}
	}

	req, err := http.NewRequest(v.Request.Method, v.Request.Base+v.Request.Path, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	if v.Request.Headers != nil {
		for headerName, headerVal := range v.Request.Headers {
			req.Header.Add(headerName, headerVal)
		}
	}

	requestInitFuncs := make([]apitestr.RequestInitFunc, 0)
	requestInitFuncsData := make([]map[string]interface{}, 0)

	if v.Request.InitFunc != nil {
		for initFuncID, initFuncData := range v.Request.InitFunc {
			initFunc := apitestr.RequestInitFuncFromContext(ctx, initFuncID)
			if initFunc == nil {
				return nil, fmt.Errorf("no request init func found with id of `%s`", initFuncID)
			}

			requestInitFuncs = append(requestInitFuncs, initFunc)
			requestInitFuncsData = append(requestInitFuncsData, initFuncData)
		}
	}

	t := &apitestr.Test{
		Name:                 v.Name,
		Group:                v.Group,
		Order:                v.Order,
		Request:              req,
		Checks:               make([]check.Checker, len(v.Checks)),
		RequestInitFuncs:     requestInitFuncs,
		RequestInitFuncsData: requestInitFuncsData,
	}

	if t.Name == "" {
		t.Name = "unknown"
	}
	if t.Group == "" {
		t.Group = "default"
	}
	if t.Order < 0 {
		t.Order = 0
	}

	for cIndex, c := range v.Checks {
		checker, err := V1Check(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("could not parse v1 check [%d]: %w", cIndex, err)
		}

		t.Checks[cIndex] = checker
	}

	return t, nil
}

func V1Check(ctx context.Context, c v1Check) (check.Checker, error) {
	switch c.Type {
	case "bodyEqual":
		value, ok := c.Data.string("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.BodyEqualChecker{Value: value}, nil

	case "dataEqual":
		id, ok := c.Data.string("id")
		if !ok {
			return nil, fmt.Errorf("missing required data `id`")
		}
		value, ok := c.Data.get("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.DataEqualChecker{Value: value, DataID: id}, nil

	case "jsonBodyEqual":
		value, ok := c.Data.get("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.BodyJSONChecker{Value: value}, nil

	case "jsonBodyQueryExists":
		query, ok := c.Data.string("query")
		if !ok {
			return nil, fmt.Errorf("missing required data `query`")
		}
		dataID, _ := c.Data.string("dataId")
		return &check.BodyJSONQueryExistsChecker{Query: query, DataID: dataID}, nil

	case "jsonBodyQueryEqual":
		query, ok := c.Data.string("query")
		if !ok {
			return nil, fmt.Errorf("missing required data `query`")
		}
		value, ok := c.Data.get("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		dataID, _ := c.Data.string("dataId")
		return &check.BodyJSONQueryEqualChecker{Query: query, Value: value, NullValue: value == nil, DataID: dataID}, nil

	case "jsonBodyQueryRegexMatch":
		query, ok := c.Data.string("query")
		if !ok {
			return nil, fmt.Errorf("missing required data `query`")
		}
		pattern, ok := c.Data.string("pattern")
		if !ok {
			return nil, fmt.Errorf("missing required data `pattern`")
		}
		r, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("could not compile regex pattern `%s`: %w", pattern, err)
		}

		var dataIDs map[int]string
		if dataIDsInterface, ok := c.Data.get("dataIds"); ok {
			switch dataIDsOfType := dataIDsInterface.(type) {
			case map[int]string:
				dataIDs = dataIDsOfType
			case map[string]string:
				dataIDs = make(map[int]string, len(dataIDsOfType))
				for k, v := range dataIDsOfType {
					intK, err := strconv.Atoi(k)
					if err != nil {
						return nil, fmt.Errorf("could not parse `dataIds` key `%v` to int: %w", k, err)
					}
					dataIDs[intK] = v
				}
			case map[string]interface{}:
				dataIDs = make(map[int]string, len(dataIDsOfType))
				for k, interfaceVal := range dataIDsOfType {
					intK, err := strconv.Atoi(k)
					if err != nil {
						return nil, fmt.Errorf("could not parse `dataIds` key `%v` to int: %w", k, err)
					}

					switch valOfType := interfaceVal.(type) {
					case string:
						dataIDs[intK] = valOfType
					case []byte:
						dataIDs[intK] = string(valOfType)
					default:
						return nil, fmt.Errorf("could not parse `dataIds` value for `%d` to string", intK)
					}
				}
			default:
				return nil, fmt.Errorf("could not parse `dataIds` data. expected type of `map[int]string` or `map[string]string`, got %T", dataIDsInterface)
			}
		} else {
			dataIDs = make(map[int]string)
		}
		if dataID, _ := c.Data.string("dataId"); dataID != "" {
			dataIDs[0] = dataID
		}
		return &check.BodyJSONQueryRegexMatchChecker{Query: query, Regexp: r, DataIDs: dataIDs}, nil

	case "statusCodeEqual":
		value, ok := c.Data.int("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.StatusCodeEqualChecker{Value: value}, nil

	case "bodyCustom":
		value, ok := c.Data.string("id")
		if !ok {
			return nil, fmt.Errorf("missing required data `id`")
		}
		checkFunc := apitestr.CustomBodyCheckFromContext(ctx, value)
		if checkFunc == nil {
			return nil, fmt.Errorf("no custom body check found with id of `%s`", value)
		}
		return &check.BodyCustomChecker{CheckBody: checkFunc}, nil

	default:
		return nil, fmt.Errorf("unhandled type `%s`", c.Type)
	}
}
