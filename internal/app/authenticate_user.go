package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-nextjs-dashboard/internal/auth"
	"go-nextjs-dashboard/internal/hashing"
	"go-nextjs-dashboard/internal/user"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthenticateUser struct {
	auth   *auth.Service
	usrSvc *user.Service
	hash   *hashing.Hash
}

func NewAuthenticateUser(usrSvc *user.Service, auth *auth.Service, hash *hashing.Hash) *AuthenticateUser {
	return &AuthenticateUser{
		usrSvc: usrSvc,
		auth:   auth,
		hash:   hash,
	}
}

func (u AuthenticateUser) Execute(ctx context.Context, username, password string) (AccessToken, error) {
	var accessToken AccessToken

	usr, err := u.usrSvc.GetByEmail(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			return accessToken, user.ErrUserNotFound
		default:
			return accessToken, fmt.Errorf("get user by email: %w", err)
		}
	}

	match, err := u.hash.Check(password, usr.Password)
	if err != nil {
		return accessToken, fmt.Errorf("check password hash: %w", err)
	}
	if !match {
		return accessToken, auth.ErrPasswordIncorrect
	}

	newAccess, exp, err := u.auth.NewJWT(usr.ID)
	if err != nil {
		return accessToken, fmt.Errorf("new access token: %w", err)
	}

	refresh, err := u.auth.CreateRefreshToken(ctx, usr.ID)
	if err != nil {
		return accessToken, fmt.Errorf("create refresh token: %w", err)
	}

	accessToken = AccessToken{
		AccessToken:  newAccess,
		RefreshToken: refresh,
		ExpiresIn:    int(time.Until(exp).Seconds()),
	}

	return accessToken, nil
}
