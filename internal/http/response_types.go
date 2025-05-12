package http

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
