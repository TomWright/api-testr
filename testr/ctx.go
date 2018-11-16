package testr

import (
	"context"
	"github.com/tomwright/api-testr/testr/check"
)

type ctxKey string

const (
	ctxBaseURLKey     ctxKey = "baseUrl"
	ctxCustomCheckKey ctxKey = "customBodyCheck_"
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
