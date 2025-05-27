package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRefreshSessionNotFound = errors.New("refresh session not found")
)

type RefreshSession struct {
	ID        uuid.UUID // refresh token sent to the client
	UserID    uuid.UUID
	ExpiresAt time.Time
	Used      bool // single-use rotation flag
	CreatedAt time.Time
}

type RefreshStore interface {
	Get(context.Context, uuid.UUID) (RefreshSession, error)
	Insert(context.Context, RefreshSession) (RefreshSession, error)
	Update(context.Context, RefreshSession) error
}
