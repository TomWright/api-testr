package testr

import (
	"context"
	"net/http"
)

// RequestInitFunc defines the structure of the functions that can be used to initialise a request
type RequestInitFunc func(ctx context.Context, req *http.Request, data map[string]interface{}) (*http.Request, error)
