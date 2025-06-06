package response

import (
	"github.com/gelozr/go-dash/internal/auth"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func ToAccessToken(token auth.AccessToken) AccessToken {
	return AccessToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}
}
