package parse

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tomwright/api-testr/testr"
	"io/ioutil"
)

type version struct {
	Version int `json:"version"`
}

func File(ctx context.Context, path string) (*testr.Test, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read test file: %s", err)
	}
	return Parse(ctx, data)
}

func Parse(ctx context.Context, data []byte) (*testr.Test, error) {
	v := version{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("could not unmarshal version data: %s", err)
	}

	var t *testr.Test
	var err error

	switch v.Version {
	case 1:
		t, err = V1(ctx, data)
	case 0:
		fallthrough
	default:
		return nil, fmt.Errorf("unhandled test version `%d`", v.Version)
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}
