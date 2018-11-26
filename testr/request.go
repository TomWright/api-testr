package testr

import "net/http"

// RequestInitFunc defines the structure of the functions that can be used to initialise a request
type RequestInitFunc func(req *http.Request, data map[string]interface{}) (*http.Request, error)
