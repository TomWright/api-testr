package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tomwright/api-testr/testr"
	"github.com/tomwright/api-testr/testr/check"
	"net/http"
)

type v1 struct {
	Name    string    `json:"name"`
	Group   string    `json:"group"`
	Order   int       `json:"order"`
	Request v1Request `json:"request"`
	Checks  []v1Check `json:"checks"`
}

type v1Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}

type v1Check struct {
	Type string `json:"type"`
	Data *Data  `json:"data"`
}

func V1(data []byte, baseAddr string) (*testr.Test, error) {
	v := v1{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("could not unmarshal v1 test data: %s", err)
	}

	// var bytesBuffer *bytes.Buffer = nil
	// if len(v.Request.Body) > 0 {
	// 	bytesBuffer = bytes.NewBuffer([]byte(v.Request.Body))
	// }
	req, err := http.NewRequest(v.Request.Method, baseAddr+v.Request.Path, bytes.NewBuffer([]byte(v.Request.Body)))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	t := &testr.Test{
		Name:    v.Name,
		Group:   v.Group,
		Order:   v.Order,
		Request: req,
		Checks:  make([]check.Checker, len(v.Checks)),
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
		checker, err := V1Check(c)
		if err != nil {
			return nil, fmt.Errorf("could not parse v1 check [%d]: %s", cIndex, err)
		}

		t.Checks[cIndex] = checker
	}

	return t, nil
}

func V1Check(c v1Check) (check.Checker, error) {
	switch c.Type {
	case "bodyEqual":
		value, ok := c.Data.String("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.BodyEqualChecker{Value: value}, nil

	case "jsonBodyEqual":
		value, ok := c.Data.Get("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.BodyJSONChecker{Value: value}, nil

	case "jsonBodyQueryExists":
		query, ok := c.Data.String("query")
		if !ok {
			return nil, fmt.Errorf("missing required data `query`")
		}
		return &check.BodyJSONQueryExistsChecker{Query: query}, nil

	case "jsonBodyQueryEqual":
		query, ok := c.Data.String("query")
		if !ok {
			return nil, fmt.Errorf("missing required data `query`")
		}
		value, ok := c.Data.Get("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		valueType, _ := c.Data.String("valueType")
		return &check.BodyJSONQueryEqualChecker{Query: query, Value: value, NullValue: value == nil, ValueType: valueType}, nil

	case "statusCodeEqual":
		value, ok := c.Data.Int("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.StatusCodeEqualChecker{Value: value}, nil

	default:
		return nil, fmt.Errorf("unhandled type `%s`", c.Type)
	}
}
