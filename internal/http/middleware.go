package http

import (
	"context"
	"errors"
	"net/http"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/google/uuid"
)

func init() {
	limiter.ConfigDefault.LimitReached = func(ctx fiber.Ctx) error {
		return fiber.NewError(http.StatusTooManyRequests, "Too Many Requests")
	}
}

// ValidationResponse maps the validation errors into a JSON response
func ValidationResponse() fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			var vErrs govalidator.ValidationErrors

			if errors.As(err, &vErrs) {
				trans, found := Uni.GetTranslator(c.Get("Accept-Language"))
				if !found {
					trans, _ = Uni.GetTranslator("en")
				}

				out := make(map[string]string, len(vErrs))
				for _, e := range vErrs {
					out[e.Field()] = e.Translate(trans)
				}

				return c.Status(http.StatusUnprocessableEntity).
					JSON(ValidationErrResponse{
						Message: "The given data was invalid.",
						Errors:  out,
					})
			}
		}

		return err
	}
}

type ctxKey string

var reqIDKey ctxKey = "req_id"

// RequestID extracts the request id from the request header or generates a new one
func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		const hdr = "X-Request-Id"

		id := c.Get(hdr)
		if len(id) == 0 || len(id) > 64 {
			id = uuid.NewString()
		}

		ctx := context.WithValue(c.Context(), reqIDKey, id)
		c.SetContext(ctx)

		c.Request().Header.Set(hdr, id)
		c.Response().Header.Set(hdr, id)

		return c.Next()
	}
}

func rateLimiter(max ...int) fiber.Handler {
	var m int
	if len(max) > 0 {
		m = max[0]
	}

	return limiter.New(limiter.Config{Max: m})
}

func loggerKeyMiddleware(key string) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := context.WithValue(c.Context(), "logger_key", key)
		c.SetContext(ctx)
		return c.Next()
	}
}
