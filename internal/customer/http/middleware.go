package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func rateLimiter(max ...int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max: max[0],
	})
}
