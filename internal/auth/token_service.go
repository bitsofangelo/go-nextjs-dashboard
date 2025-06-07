package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	refreshStore RefreshStore
	// logger       logger.Logger
}

func NewToken(refreshStore RefreshStore) *Token {
	return &Token{
		refreshStore: refreshStore,
		// logger:       logger.With("component", "auth"),
	}
}

func (a *Token) GetRefresh(ctx context.Context, id uuid.UUID) (RefreshSession, error) {
	r, err := a.refreshStore.Get(ctx, id)
	if err != nil {
		return RefreshSession{}, fmt.Errorf("get refresh: %w", err)
	}
	return r, nil
}

func (a *Token) CreateRefresh(ctx context.Context, uid uuid.UUID) (RefreshSession, error) {
	r, err := a.refreshStore.Insert(ctx, RefreshSession{
		ID:        uuid.New(),
		UserID:    uid,
		Used:      false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	})

	if err != nil {
		return RefreshSession{}, fmt.Errorf("create refresh session: %w", err)
	}

	return r, nil
}

func (a *Token) ExchangeRefresh(ctx context.Context, currRefresh RefreshSession) (RefreshSession, error) {
	if currRefresh.ExpiresAt.Before(time.Now()) {
		return RefreshSession{}, ErrRefreshTokenExpired
	}

	if currRefresh.Used {
		return RefreshSession{}, ErrRefreshTokenUsed
	}

	currRefresh.Used = true
	if err := a.refreshStore.Update(ctx, currRefresh); err != nil {
		return RefreshSession{}, fmt.Errorf("update refresh session: %w", err)
	}

	newRefresh, err := a.CreateRefresh(ctx, currRefresh.UserID)
	if err != nil {
		return RefreshSession{}, fmt.Errorf("create refresh token: %w", err)
	}

	return newRefresh, nil
}
