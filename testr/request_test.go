package testr_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tomwright/api-testr/testr"
	"github.com/tomwright/api-testr/testr/check"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func headerToStr(header http.Header) string {
	res := ""
	for k, v := range header {
		for _, v2 := range v {
			res += fmt.Sprintf("%s=%s&", k, v2)
		}
	}
	return strings.TrimRight(res, "&")
}

func bodyToStr(r *http.Request) string {
	if r.Body == nil {
		return ""
	}
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return string(data)
}

func TestRequestReplacements(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		desc            string
		replacements    map[string]interface{}
		url             string
		body            []byte
		headers         map[string]string
		expectedUrl     string
		expectedBody    []byte
		expectedHeaders map[string]string
		ctx             func(ctx context.Context) context.Context
	}{
		{
			desc: "domain replacements work",
			replacements: map[string]interface{}{
				"old.com": "new.com",
			},
			url:         "https://old.com",
			expectedUrl: "https://new.com",
		},
		{
			desc: "scheme replacements work",
			replacements: map[string]interface{}{
				"https": "http",
			},
			url:         "https://example.com",
			expectedUrl: "http://example.com",
		},
		{
			desc: "standard string replacements work",
			replacements: map[string]interface{}{
				":name:": "Tom",
			},
			url:             "https://example.com/users/:name:?name=:name:",
			body:            []byte("hello :name:"),
			headers:         map[string]string{"Custom-Name": ":name:"},
			expectedUrl:     "https://example.com/users/Tom?name=Tom",
			expectedBody:    []byte("hello Tom"),
			expectedHeaders: map[string]string{"Custom-Name": "Tom"},
		},
		{
			desc: "variable replacements work",
			replacements: map[string]interface{}{
				":name:": "$.replacementTest.name",
			},
			url:             "https://example.com/users/:name:?name=:name:",
			body:            []byte("hello :name:"),
			headers:         map[string]string{"Custom-Name": ":name:"},
			expectedUrl:     "https://example.com/users/Tom?name=Tom",
			expectedBody:    []byte("hello Tom"),
			expectedHeaders: map[string]string{"Custom-Name": "Tom"},
			ctx: func(ctx context.Context) context.Context {
				ctx = check.ContextWithData(ctx, make(map[string]interface{}))
				if err := check.ContextWithDataID(ctx, "replacementTest.name", "Tom"); err != nil {
					panic(err)
				}
				return ctx
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tc.ctx != nil {
				ctx = tc.ctx(ctx)
			}

			inputReq, _ := http.NewRequest("POST", tc.url, bytes.NewBuffer(tc.body))
			for k, v := range tc.headers {
				inputReq.Header.Add(k, v)
			}

			expectedReq, _ := http.NewRequest("POST", tc.expectedUrl, bytes.NewBuffer(tc.expectedBody))
			for k, v := range tc.expectedHeaders {
				expectedReq.Header.Add(k, v)
			}

			resultReq, err := testr.RequestReplacements(ctx, inputReq, tc.replacements)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if exp, got := expectedReq.URL.String(), resultReq.URL.String(); exp != got {
				t.Errorf("expected url of `%s`. got `%s`", exp, got)
			}

			if exp, got := headerToStr(expectedReq.Header), headerToStr(resultReq.Header); exp != got {
				t.Errorf("expected headers of `%s`. got `%s`", exp, got)
			}

			if exp, got := bodyToStr(expectedReq), bodyToStr(resultReq); exp != got {
				t.Errorf("expected body of `%s`. got `%s`", exp, got)
			}
		})
	}
}
