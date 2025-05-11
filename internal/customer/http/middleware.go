package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func rateLimiter(max ...int) fiber.Handler {
	var m int
	if len(max) > 0 {
		m = max[0]
	}

	return limiter.New(limiter.Config{Max: m})
}
