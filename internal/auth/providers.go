package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/gelozr/go-dash/internal/hashing"
	"github.com/gelozr/go-dash/internal/user"
)

type DBProvider[U any] struct {
	userSvc *user.Service
	hash    hashing.Hasher
}

var _ Authenticator[any] = (*DBProvider[any])(nil)

func NewDBProvider[U any](userSvc *user.Service, hash hashing.Hasher) *DBProvider[U] {
	return &DBProvider[U]{
		userSvc: userSvc,
		hash:    hash,
	}
}

func (p DBProvider[U]) Authenticate(ctx context.Context, creds PasswordCredentials) (*U, error) {
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

	return any(u).(*U), nil
}
