package http

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/user"
)

type userHandler struct {
	svc    *user.Service
	logger logger.Logger
}

func newUserHandler(svc *user.Service, log logger.Logger) *userHandler {
	return &userHandler{
		svc:    svc,
		logger: log.With("component", "http.user"),
	}
}

func (h *userHandler) GetByEmail(c fiber.Ctx) error {
	email := c.Params("email")

	u, err := h.svc.GetByEmail(c.Context(), email)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			return fiber.NewError(fiber.StatusNotFound, "User not found.")
		default:
			return fmt.Errorf("get user by email: %w", err)
		}
	}

	res := struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}

	return c.JSON(Response{
		Data: res,
	})
}
