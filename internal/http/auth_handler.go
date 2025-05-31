package http

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/gelozr/go-dash/internal/app"
	"github.com/gelozr/go-dash/internal/auth"
	"github.com/gelozr/go-dash/internal/http/request"
	"github.com/gelozr/go-dash/internal/http/response"
	"github.com/gelozr/go-dash/internal/user"
)

type AuthHandler struct {
	authUser *app.AuthenticateUser
}

func NewAuthHandler(authUser *app.AuthenticateUser) *AuthHandler {
	return &AuthHandler{
		authUser: authUser,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var login request.Login

	if err := c.Bind().Body(&login); err != nil {
		return fmt.Errorf("error parsing login request: %w", err)
	}

	creds := auth.PasswordCredentials{
		Email:    login.Username,
		Password: login.Password,
	}

	accessToken, err := h.authUser.Execute(c.Context(), creds)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound), errors.Is(err, auth.ErrPasswordIncorrect):
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		default:
			return fmt.Errorf("error executing login: %w", err)
		}
	}

	return c.JSON(
		response.New(response.ToAccessToken(accessToken)),
	)
}
