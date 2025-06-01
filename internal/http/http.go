package http

import (
	"context"
	"strconv"

	"github.com/google/uuid"
)

type Response struct {
	Data any `json:"data"`
}

func ReqID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(reqIDCtxKey).(string)
	return v, ok
}

func Locale(ctx context.Context, def ...string) string {
	if v, ok := ctx.Value(reqLocaleCtxKey).(string); ok {
		return v
	}

	if len(def) > 0 {
		return def[0]
	}

	return "en"
}

func UserIDFromCtx(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(userIDCtxKey).(uuid.UUID)
	return v, ok
}

func getDefaultNum[T any](value string, def T) T {
	switch any(def).(type) {
	case int, int8, int16, int32, int64:
		i, err := strconv.Atoi(value)
		if err != nil {
			return def
		}
		return any(i).(T)
	case float32, float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return any(f).(T)
	default:
		return def
	}
}
