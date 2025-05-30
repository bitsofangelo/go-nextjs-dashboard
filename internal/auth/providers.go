package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/gelozr/go-dash/internal/hashing"
	"github.com/gelozr/go-dash/internal/user"
)

type PasswordCredentials struct {
	Username string
	Password string
}

type PasswordProvider struct {
	userSvc *user.Service
	hash    *hashing.Hash
}

var _ ProviderDriver = (*PasswordProvider)(nil)

func NewPasswordProvider(userSvc *user.Service, hash *hashing.Hash) *PasswordProvider {
	return &PasswordProvider{
		userSvc: userSvc,
		hash:    hash,
	}
}

func (p PasswordProvider) Authenticate(ctx context.Context, credentials Credentials) (*user.User, error) {
	creds, ok := credentials.(PasswordCredentials)
	if !ok {
		return nil, errors.New("invalid credentials type")
	}

	u, err := p.userSvc.GetByEmail(ctx, creds.Username)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			return nil, user.ErrUserNotFound
		default:
			return nil, fmt.Errorf("get user by email: %w", err)
		}
	}

	match, err := p.hash.Check(creds.Password, u.Password)
	if err != nil {
		return nil, fmt.Errorf("check password hash: %w", err)
	}
	if !match {
		return u, ErrPasswordIncorrect
	}

	return u, nil
}

type GoogleCredentials struct {
	IDToken string `json:"id_token"`
}

type GoogleProvider struct {
}

var _ ProviderDriver = (*GoogleProvider)(nil)

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{}
}

func (p GoogleProvider) Authenticate(ctx context.Context, credentials Credentials) (*user.User, error) {
	return nil, errors.New("not implemented")
}
