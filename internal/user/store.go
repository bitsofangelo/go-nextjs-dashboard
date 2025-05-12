package user

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("customer not found")

type Store interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}
