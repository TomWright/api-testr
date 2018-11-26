package testr

import (
	"context"
	"github.com/tomwright/api-testr/testr/check"
)

type ctxKey string

const (
	ctxBaseURLKey      ctxKey = "baseUrl"
	ctxCustomCheckKey  ctxKey = "customBodyCheck_"
	ctxRequestInitFunc ctxKey = "requestInitFunc_"
)

// ContextWithBaseURL stores the given base URL in the context
func ContextWithBaseURL(ctx context.Context, baseURL string) context.Context {
	return context.WithValue(ctx, ctxBaseURLKey, baseURL)
}

// BaseURLFromContext returns the base URL to be used in tests, as stored in the given context
func BaseURLFromContext(ctx context.Context) string {
	val := ctx.Value(ctxBaseURLKey)
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// ContextWithCustomBodyCheck stores a BodyCustomCheckerFunc into the context under the given id
func ContextWithCustomBodyCheck(ctx context.Context, checkID string, checkFunc check.BodyCustomCheckerFunc) context.Context {
	return context.WithValue(ctx, ctxCustomCheckKey+ctxKey(checkID), checkFunc)
}

// CustomBodyCheckFromContext retrieves a BodyCustomCheckerFunc from the context by id
func CustomBodyCheckFromContext(ctx context.Context, checkID string) check.BodyCustomCheckerFunc {
	val := ctx.Value(ctxCustomCheckKey + ctxKey(checkID))
	if val == nil {
		return nil
	}
	if checkFunc, ok := val.(check.BodyCustomCheckerFunc); ok {
		return checkFunc
	}
	return nil
}

// ContextWithRequestInitFunc stores a RequestInitFunc into the context under the given id
func ContextWithRequestInitFunc(ctx context.Context, initFuncID string, requestInitFunc RequestInitFunc) context.Context {
	return context.WithValue(ctx, ctxRequestInitFunc+ctxKey(initFuncID), requestInitFunc)
}

// RequestInitFuncFromContext retrieves a RequestInitFunc from the context by id
func RequestInitFuncFromContext(ctx context.Context, initFuncID string) RequestInitFunc {
	val := ctx.Value(ctxRequestInitFunc + ctxKey(initFuncID))
	if val == nil {
		return nil
	}
	if requestInitFunc, ok := val.(RequestInitFunc); ok {
		return requestInitFunc
	}
	return nil
}
