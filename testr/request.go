package testr

import "net/http"

type RequestInitFunc func(req *http.Request, data map[string]interface{}) (*http.Request, error)
