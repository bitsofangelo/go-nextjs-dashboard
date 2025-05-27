package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/http/response"
	"go-nextjs-dashboard/internal/http/validation"
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
			var vErrs validation.Errors
			var jsonErr *json.UnmarshalTypeError

			switch {
			case errors.As(err, &vErrs):
				return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationError{
					Message: "The given data was invalid.",
					Errors:  vErrs,
				})

			case errors.As(err, &jsonErr):
				return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationError{
					Message: "The given data was invalid.",
					Errors: map[string][]string{
						jsonErr.Field: {
							fmt.Sprintf("%s must be a of type %s", jsonErr.Field, jsonErr.Type.String()),
						},
					},
				})
			}
		}

		return err
	}
}

type ctxKey string

var reqIDKey ctxKey = "req_id"
var reqLocale ctxKey = "req_locale"

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

func RequestLocale() fiber.Handler {
	return func(c fiber.Ctx) error {
		const hdr = "Accept-Language"

		l := c.Get(hdr)
		if l == "" {
			l = "en"
		}

		ctx := context.WithValue(c.Context(), reqLocale, l)
		c.SetContext(ctx)

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
