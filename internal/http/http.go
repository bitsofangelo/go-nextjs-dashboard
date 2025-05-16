package http

import "context"

type Response struct {
	Data any `json:"data"`
}

type ErrResponse struct {
	Message string `json:"message"`
}

type ValidationErrResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func ReqID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(reqIDKey).(string)
	return v, ok
}
