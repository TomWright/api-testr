package apitestr_test

import (
	"context"
	"encoding/json"
	"github.com/tomwright/apitestr"
	"github.com/tomwright/apitestr/parse"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, _ := json.Marshal(map[string]interface{}{
			"completed": false,
			"id":        1,
			"title":     "delectus aut autem",
			"userId":    1,
		})
		_, _ = w.Write(res)
	}))
	defer ts.Close()

	ctx := apitestr.ContextWithBaseURL(context.Background(), ts.URL)

	te, err := parse.File(ctx, "tests/example.json")
	if err != nil {
		t.Errorf("unexpected error parsing file: %s", err)
		return
	}

	if err := apitestr.Run(ctx, te, nil, nil); err != nil {
		t.Errorf("unexpected error in test: %s", err)
		return
	}
}
