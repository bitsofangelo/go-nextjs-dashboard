package http

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	himoauth "github.com/gelozr/himo/auth2"

	"github.com/gelozr/go-dash/internal/auth"
	"github.com/gelozr/go-dash/internal/http/request"
	"github.com/gelozr/go-dash/internal/http/response"
	"github.com/gelozr/go-dash/internal/http/validation"
	"github.com/gelozr/go-dash/internal/user"
)

type AuthHandler struct {
	auth      himoauth.Auth
	validator validation.Validator
}

func NewAuthHandler(auth himoauth.Auth, validator validation.Validator) *AuthHandler {
	return &AuthHandler{
		auth:      auth,
		validator: validator,
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

	usr, err := h.auth.MustHandler("jwt").Authenticate(ctx, creds)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound), errors.Is(err, auth.ErrPasswordIncorrect):
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		default:
			return fmt.Errorf("authenticate user: %w", err)
		}
	}

	login, err := h.auth.Login(ctx, usr)
	if err != nil && !errors.Is(err, himoauth.ErrLoginNotSupported) {
		switch {
		default:
			return fmt.Errorf("login: %w", err)
		}
	}

	var r any
	switch v := login.(type) {
	case auth.AccessToken:
		r = response.ToAccessToken(v)
	default:
		r = true
	}

	return c.JSON(
		response.New(r),
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

	token, err := h.auth.RefreshToken(ctx, tokenID.String())
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrRefreshSessionNotFound),
			errors.Is(err, auth.ErrRefreshTokenUserMismatch),
			errors.Is(err, auth.ErrRefreshTokenUsed),
			errors.Is(err, auth.ErrRefreshTokenInvalid):
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		default:
			return fmt.Errorf("refresh access token: %w", err)
		}
	}

	var res any
	switch v := token.(type) {
	case auth.AccessToken:
		res = response.ToAccessToken(v)
	default:
		return fmt.Errorf("failed to refresh token")
	}

	return c.JSON(
		response.New(res),
	)
}
