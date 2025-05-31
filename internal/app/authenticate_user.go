package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gelozr/go-dash/internal/auth"
	"github.com/gelozr/go-dash/internal/user"
)

type AccessToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type AuthenticateUser struct {
	auth  auth.Authenticator[user.User]
	token *auth.Token
}

func NewAuthenticateUser(auth auth.Authenticator[user.User], token *auth.Token) *AuthenticateUser {
	return &AuthenticateUser{
		auth:  auth,
		token: token,
	}
}

func (u *AuthenticateUser) Execute(ctx context.Context, creds auth.PasswordCredentials) (AccessToken, error) {
	var accessToken AccessToken

	usr, err := u.auth.Authenticate(ctx, creds)
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
