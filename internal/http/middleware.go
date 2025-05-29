package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/auth"
	"go-nextjs-dashboard/internal/http/response"
	"go-nextjs-dashboard/internal/http/validation"
)

func init() {
	limiter.ConfigDefault.LimitReached = func(ctx fiber.Ctx) error {
		return fiber.NewError(http.StatusTooManyRequests, "Too Many Requests")
	}
}

type ctxKey string

var (
	reqIDKey  ctxKey = "req_id"
	reqLocale ctxKey = "req_locale"
)

func AuthMiddleware(token *auth.Token) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenStr := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if tokenStr == "" {
			return fiber.NewError(http.StatusUnauthorized, "missing authorization header")
		}

		claims, err := token.ParseJWT(tokenStr)
		if err != nil {
			switch {
			case errors.Is(err, auth.ErrJWTInvalid):
				return fiber.NewError(http.StatusUnauthorized, "invalid token")
			case errors.Is(err, auth.ErrJWTExpired):
				return fiber.NewError(http.StatusUnauthorized, "expired token")
			default:
				return fmt.Errorf("parse token: %w", err)
			}
		}

		c.Locals("user", claims.UserID)

		return c.Next()
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
