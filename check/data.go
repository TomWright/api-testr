package check

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	dataCtxKey ctxKey = "ctxData"
)

func DataFromContext(ctx context.Context) map[string]interface{} {
	val := ctx.Value(dataCtxKey)
	if val == nil {
		return nil
	}
	if mapVal, ok := val.(map[string]interface{}); ok {
		return mapVal
	}
	return nil
}

func DataIDFromContext(ctx context.Context, id string) interface{} {
	data := DataFromContext(ctx)
	if data == nil {
		return nil
	}
	if val, ok := data[id]; ok {
		return val
	}
	return nil
}

func ContextWithData(ctx context.Context, data map[string]interface{}) context.Context {
	return context.WithValue(ctx, dataCtxKey, data)
}

func ContextWithDataID(ctx context.Context, id string, val interface{}) error {
	data := DataFromContext(ctx)
	if data == nil {
		return fmt.Errorf("could not store data id `%s` in context: data map is not present", id)
	}
	data[id] = val
	return nil
}

func ContextWithOptionalDataID(ctx context.Context, id string, val interface{}) error {
	if id == "" {
		return nil
	}
	return ContextWithDataID(ctx, id, val)
}
