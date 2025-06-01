package http

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/gelozr/go-dash/internal/app"
	"github.com/gelozr/go-dash/internal/auth"
	"github.com/gelozr/go-dash/internal/http/request"
	"github.com/gelozr/go-dash/internal/http/response"
	"github.com/gelozr/go-dash/internal/http/validation"
	"github.com/gelozr/go-dash/internal/user"
)

type AuthHandler struct {
	authUser           *app.AuthenticateUser
	refreshAccessToken *app.RefreshAccessToken
	validator          validation.Validator
}

func NewAuthHandler(authUser *app.AuthenticateUser, refreshAccessToken *app.RefreshAccessToken, validator validation.Validator) *AuthHandler {
	return &AuthHandler{
		authUser:           authUser,
		refreshAccessToken: refreshAccessToken,
		validator:          validator,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req request.Login

	if err := c.Bind().Body(&req); err != nil {
		return fmt.Errorf("parsing login request: %w", err)
	}

	ctx := c.Context()

	if err := h.validator.ValidateStruct(ctx, req); err != nil {
		return fmt.Errorf("login request validation: %w", err)
	}

	creds := auth.PasswordCredentials{
		Email:    req.Username,
		Password: req.Password,
	}

	accessToken, err := h.authUser.Execute(ctx, creds)
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

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req request.Refresh

	if err := c.Bind().Body(&req); err != nil {
		return fmt.Errorf("parsing refresh request: %w", err)
	}

	ctx := c.Context()

	if err := h.validator.ValidateStruct(ctx, req); err != nil {
		return fmt.Errorf("refresh request validation: %w", err)
	}

	tokenID, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		return response.NewError("invalid refresh token", fiber.StatusUnauthorized, err)
	}

	accessToken, err := h.refreshAccessToken.Execute(ctx, tokenID)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrRefreshSessionNotFound),
			errors.Is(err, auth.ErrRefreshTokenUserMismatch),
			errors.Is(err, auth.ErrRefreshTokenUsed):
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		default:
			return fmt.Errorf("refresh access token: %w", err)
		}
	}

	return c.JSON(
		response.New(response.ToAccessToken(accessToken)),
	)
}
