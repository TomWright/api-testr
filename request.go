package apitestr

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tomwright/apitestr/check"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// RequestInitFunc defines the structure of the functions that can be used to initialise a request
type RequestInitFunc func(ctx context.Context, req *http.Request, data map[string]interface{}) (*http.Request, error)

// RequestReplacements runs a find and replace in the request body, headers and url path with the given replacements in `data`
func RequestReplacements(ctx context.Context, req *http.Request, data map[string]interface{}) (*http.Request, error) {
	req, err := RequestURLReplacements(ctx, req, data)
	if err != nil {
		return req, err
	}
	req, err = RequestHeaderReplacements(ctx, req, data)
	if err != nil {
		return req, err
	}
	req, err = RequestBodyReplacements(ctx, req, data)
	if err != nil {
		return req, err
	}
	return req, err
}

// RequestSchemeReplacements runs a find and replace in the request url scheme with the given replacements in `data`
func RequestURLReplacements(ctx context.Context, req *http.Request, data map[string]interface{}) (*http.Request, error) {
	urlStr := req.URL.String()
	for k, v := range data {
		vStr, err := getReplacementValue(ctx, v)
		if err != nil {
			return nil, err
		}
		urlStr = strings.Replace(urlStr, k, vStr, -1)
	}
	var err error
	req.URL, err = url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// RequestHeaderReplacements runs a find and replace in the request headers with the given replacements in `data`
func RequestHeaderReplacements(ctx context.Context, req *http.Request, data map[string]interface{}) (*http.Request, error) {
	for k, v := range data {
		vStr, err := getReplacementValue(ctx, v)
		if err != nil {
			return nil, err
		}
		for headerIndex, headerVals := range req.Header {
			for headerValIndex, h := range headerVals {
				req.Header[headerIndex][headerValIndex] = strings.Replace(h, k, vStr, -1)
			}
		}
	}
	return req, nil
}

// RequestBodyReplacements runs a find and replace in the request body with the given replacements in `data`
func RequestBodyReplacements(ctx context.Context, req *http.Request, data map[string]interface{}) (*http.Request, error) {
	if req.Body == nil {
		return req, nil
	}
	if len(data) == 0 {
		return req, nil
	}
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	bodyStr := string(bodyData)
	for k, v := range data {
		vStr, err := getReplacementValue(ctx, v)
		if err != nil {
			return nil, err
		}
		bodyStr = strings.Replace(bodyStr, k, vStr, -1)
	}

	newBodyBuffer := bytes.NewBuffer([]byte(bodyStr))

	req.ContentLength = int64(newBodyBuffer.Len())
	buf := newBodyBuffer.Bytes()
	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(buf)
		return ioutil.NopCloser(r), nil
	}

	req.Body = ioutil.NopCloser(newBodyBuffer)

	return req, nil
}

func getReplacementValue(ctx context.Context, val interface{}) (string, error) {
	var valStr string

	switch valOfType := val.(type) {
	case string:
		valStr = valOfType
	case []byte:
		valStr = string(valOfType)
	case nil:
		valStr = ""
	default:
		return "", fmt.Errorf("unhandled replacement value type of `%T` with value `%v`", val, val)
	}

	if strings.HasPrefix(valStr, "$.") {
		dataVal := check.DataIDFromContext(ctx, strings.TrimLeft(valStr, "$."))
		switch valOfType := dataVal.(type) {
		case string:
			valStr = valOfType
		case []byte:
			valStr = string(valOfType)
		case nil:
			valStr = ""
		default:
			return "", fmt.Errorf("unhandled replacement variable value type of `%T` with value `%v`", val, val)
		}
	}
	return valStr, nil
}
