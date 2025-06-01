package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gelozr/go-dash/internal/auth"
)

type RefreshAccessToken struct {
	token *auth.Token
}

func NewRefreshAccessToken(token *auth.Token) *RefreshAccessToken {
	return &RefreshAccessToken{
		token: token,
	}
}

func (u RefreshAccessToken) Execute(ctx context.Context, tokenID uuid.UUID) (auth.AccessToken, error) {
	currRefresh, err := u.token.GetRefresh(ctx, tokenID)
	if err != nil {
		return auth.AccessToken{}, fmt.Errorf("get refresh session: %w", err)
	}

	newRefresh, err := u.token.ExchangeRefresh(ctx, currRefresh)
	if err != nil {
		return auth.AccessToken{}, fmt.Errorf("exchange refresh: %w", err)
	}

	jwt, exp, err := u.token.SignJWT(newRefresh.UserID)
	if err != nil {
		return auth.AccessToken{}, fmt.Errorf("sign jwt: %w", err)
	}

	accessToken := auth.AccessToken{
		AccessToken:  jwt,
		RefreshToken: newRefresh.ID.String(),
		ExpiresIn:    int(time.Until(exp).Seconds()),
	}

	return accessToken, nil
}
