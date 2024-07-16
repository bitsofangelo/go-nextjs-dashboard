package request

type Payload map[string]any

type Validator interface {
	Validate(payload Payload, rules map[string]any) error
}

type Preparer[T any] interface {
	Prepare(dto T, payload Payload) T
}

type RequestContext interface {
	Params(key string, defaultValue ...string) string
	BodyParser(interface{}) error
	Query(key string, defaultValue ...string) string
	FormValue(key string, defaultValue ...string) string
}
