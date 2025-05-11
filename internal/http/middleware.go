package http

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func init() {
	limiter.ConfigDefault.LimitReached = func(ctx fiber.Ctx) error {
		return fiber.NewError(http.StatusTooManyRequests, "Too Many Requests")
	}
}

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
