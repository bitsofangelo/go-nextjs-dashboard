package http

import (
	"context"
)

func ReqID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(reqIDKey).(string)
	return v, ok
}
