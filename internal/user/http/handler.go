package http

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"

	"go-nextjs-dashboard/internal/http"
	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/user"
	"go-nextjs-dashboard/internal/user/service"
)

func RegisterHTTP(r fiber.Router, svc *service.Service, log logger.Logger) {
	h := newHandler(svc, log)
	r.Get("/users/email/:email", h.GetByEmail)
}

type handler struct {
	svc    *service.Service
	logger logger.Logger
}

func newHandler(svc *service.Service, log logger.Logger) *handler {
	return &handler{
		svc:    svc,
		logger: log,
	}
}

func (h *handler) GetByEmail(c fiber.Ctx) error {
	email := c.Path("email")

	u, err := h.svc.GetByEmail(c.Context(), email)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			return fiber.NewError(fiber.StatusNotFound, "User not found.")
		default:
			return fmt.Errorf("get user by email: %w", err)
		}
	}

	return c.JSON(http.Response{
		Data: toResponse(u),
	})
}
