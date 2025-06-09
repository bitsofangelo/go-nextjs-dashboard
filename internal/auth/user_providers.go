package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/gelozr/go-dash/internal/hashing"
	"github.com/gelozr/go-dash/internal/user"
)

type PasswordCredentials struct {
	Email    string
	Password string
}

type DBUserProvider struct {
	userSvc *user.Service
	hash    hashing.Hasher
}

// var _ auth.UserProvider = (*DBUserProvider)(nil)

func NewDBUserProvider(userSvc *user.Service, hash hashing.Manager) *DBUserProvider {
	return &DBUserProvider{
		userSvc: userSvc,
		hash:    hash,
	}
}

func (p DBUserProvider) FindByCredentials(ctx context.Context, credentials any) (any, error) {
	creds, ok := credentials.(PasswordCredentials)
	if !ok {
		return nil, errors.New("invalid credentials type")
	}

	u, err := p.userSvc.GetByEmail(ctx, creds.Email)
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
		return nil, ErrPasswordIncorrect
	}

	return u, nil
}
