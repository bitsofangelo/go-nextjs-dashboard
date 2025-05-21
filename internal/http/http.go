package http

import (
	"context"
	"strconv"
)

type Response struct {
	Data any `json:"data"`
}

func ReqID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(reqIDKey).(string)
	return v, ok
}

type NumberType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func getDefaultNum[T NumberType](value string, def T) T {
	switch any(def).(type) {
	case int, int8, int16, int32, int64:
		i, err := strconv.Atoi(value)
		if err != nil {
			return def
		}
		return T(i)
	case float32, float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return T(f)
	default:
		return def
	}
}
