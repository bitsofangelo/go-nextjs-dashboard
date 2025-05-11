package http

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func Logger(c fiber.Ctx) *slog.Logger {
	if lg, ok := c.Locals("log").(*slog.Logger); ok {
		return lg
	}
	return slog.Default() // fallback; never panic
}
