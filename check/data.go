package check

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	dataCtxKey ctxKey = "ctxData"
)

// DataFromContext returns a map of check data from the context.
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

// DataIDFromContext returns data with the given ID from the context.
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

// ContextWithData embeds the given data in the context.
func ContextWithData(ctx context.Context, data map[string]interface{}) context.Context {
	return context.WithValue(ctx, dataCtxKey, data)
}

// ContextWithDataID embeds the given data item in the context.
func ContextWithDataID(ctx context.Context, id string, val interface{}) error {
	data := DataFromContext(ctx)
	if data == nil {
		return fmt.Errorf("could not store data id `%s` in context: data map is not present", id)
	}
	data[id] = val
	return nil
}

// ContextWithDataID embeds the given data item in the context, if the ID is not empty.
func ContextWithOptionalDataID(ctx context.Context, id string, val interface{}) error {
	if id == "" {
		return nil
	}
	return ContextWithDataID(ctx, id, val)
}
