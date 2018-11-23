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

func ContextWithBaseURL(ctx context.Context, baseURL string) context.Context {
	return context.WithValue(ctx, ctxBaseURLKey, baseURL)
}

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

func ContextWithCustomBodyCheck(ctx context.Context, checkID string, checkFunc check.BodyCustomCheckerFunc) context.Context {
	return context.WithValue(ctx, ctxCustomCheckKey+ctxKey(checkID), checkFunc)
}

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

func ContextWithRequestInitFunc(ctx context.Context, initFuncID string, requestInitFunc RequestInitFunc) context.Context {
	return context.WithValue(ctx, ctxRequestInitFunc+ctxKey(initFuncID), requestInitFunc)
}

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
