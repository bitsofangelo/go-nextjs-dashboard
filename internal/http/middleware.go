package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
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
			var ves validator.ValidationErrors

			if errors.As(err, &ves) {
				trans, found := Uni.GetTranslator(c.Get("Accept-Language"))
				if !found {
					trans, _ = Uni.GetTranslator("en")
				}

				out := make(map[string]string, len(ves))
				for _, e := range ves {
					out[e.Field()] = e.Translate(trans)
				}

				return c.Status(http.StatusUnprocessableEntity).JSON(struct {
					Message string            `json:"message"`
					Errors  map[string]string `json:"errors"`
				}{
					Message: "The given data was invalid.",
					Errors:  out,
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

		c.Request().Header.Set(hdr, id)
		c.Response().Header.Set(hdr, id)

		return c.Next()
	}
}

// RequestLogger makes logger with request context
func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		log := logger.With("req_id", c.Get("X-Request-Id"), "path", c.Path())
		c.Locals("log", log)

		return c.Next()
	}
}
