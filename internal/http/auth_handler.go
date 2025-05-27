package http

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/auth"
	"go-nextjs-dashboard/internal/http/request"
	"go-nextjs-dashboard/internal/http/response"
	"go-nextjs-dashboard/internal/user"
)

type AuthHandler struct {
	authSvc  *auth.Service
	authUser *app.AuthenticateUser
}

func NewAuthHandler(authSvc *auth.Service, authUser *app.AuthenticateUser) *AuthHandler {
	return &AuthHandler{
		authSvc:  authSvc,
		authUser: authUser,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var login request.Login

	if err := c.Bind().Body(&login); err != nil {
		return fmt.Errorf("error parsing login request: %w", err)
	}

	token, err := h.authUser.Execute(c.Context(), login.Username, login.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound), errors.Is(err, auth.ErrPasswordIncorrect):
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		default:
			return fmt.Errorf("error executing login: %w", err)
		}
	}

	return c.JSON(response.New(token))
}
