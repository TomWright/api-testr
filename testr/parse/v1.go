package parse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tomwright/api-testr/testr"
	"github.com/tomwright/api-testr/testr/check"
	"net/http"
	"regexp"
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

func V1(ctx context.Context, data []byte) (*testr.Test, error) {
	v := v1{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("could not unmarshal v1 test data: %s", err)
	}

	baseAddr := testr.BaseURLFromContext(ctx)

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
		checker, err := V1Check(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("could not parse v1 check [%d]: %s", cIndex, err)
		}

		t.Checks[cIndex] = checker
	}

	return t, nil
}

func V1Check(ctx context.Context, c v1Check) (check.Checker, error) {
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
		return &check.BodyJSONQueryEqualChecker{Query: query, Value: value, NullValue: value == nil}, nil

	case "jsonBodyQueryRegexMatch":
		query, ok := c.Data.String("query")
		if !ok {
			return nil, fmt.Errorf("missing required data `query`")
		}
		pattern, ok := c.Data.String("pattern")
		if !ok {
			return nil, fmt.Errorf("missing required data `pattern`")
		}
		r, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("could not compile regex pattern `%s`: %s", pattern, err)
		}
		return &check.BodyJSONQueryRegexMatchChecker{Query: query, Regexp: r}, nil

	case "statusCodeEqual":
		value, ok := c.Data.Int("value")
		if !ok {
			return nil, fmt.Errorf("missing required data `value`")
		}
		return &check.StatusCodeEqualChecker{Value: value}, nil

	case "bodyCustom":
		value, ok := c.Data.String("id")
		if !ok {
			return nil, fmt.Errorf("missing required data `id`")
		}
		checkFunc := testr.CustomBodyCheckFromContext(ctx, value)
		if checkFunc == nil {
			return nil, fmt.Errorf("no custom body check found with id of `%s`", value)
		}
		return &check.BodyCustomChecker{CheckBody: checkFunc}, nil

	default:
		return nil, fmt.Errorf("unhandled type `%s`", c.Type)
	}
}
