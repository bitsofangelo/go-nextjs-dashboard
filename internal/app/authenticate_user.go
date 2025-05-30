package app

import (
	"context"
	"fmt"
	"time"

	"go-nextjs-dashboard/internal/auth"
	"go-nextjs-dashboard/internal/hashing"
)

type AccessToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type AuthenticateUser struct {
	auth  auth.Auth
	token *auth.Token
	hash  *hashing.Hash
}

func NewAuthenticateUser(auth auth.Auth, token *auth.Token, hash *hashing.Hash) *AuthenticateUser {
	return &AuthenticateUser{
		auth:  auth,
		token: token,
		hash:  hash,
	}
}

func (u *AuthenticateUser) Execute(ctx context.Context, provider auth.ProviderType, creds auth.Credentials) (AccessToken, error) {
	var accessToken AccessToken

	usr, err := u.auth.Provider(provider).Authenticate(ctx, creds)
	if err != nil {
		return accessToken, fmt.Errorf("authenticate: %w", err)
	}

	jwt, exp, err := u.token.SignJWT(usr.ID)
	if err != nil {
		return accessToken, fmt.Errorf("sign jwt: %w", err)
	}

	refresh, err := u.token.CreateRefresh(ctx, usr.ID)
	if err != nil {
		return accessToken, fmt.Errorf("create refresh token: %w", err)
	}

	accessToken = AccessToken{
		AccessToken:  jwt,
		RefreshToken: refresh,
		ExpiresIn:    int(time.Until(exp).Seconds()),
	}

	return accessToken, nil
}
