package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gelozr/go-dash/internal/db"
	"github.com/gelozr/go-dash/internal/logger"
)

type GormRefreshStore struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewGormRefreshStore(db *gorm.DB, logger logger.Logger) *GormRefreshStore {
	return &GormRefreshStore{
		db:     db,
		logger: logger,
	}
}

func (s *GormRefreshStore) DB(ctx context.Context) *gorm.DB {
	if gormDB, ok := db.FromCtx(ctx); ok {
		return gormDB
	}

	return s.db
}

func (s *GormRefreshStore) Get(ctx context.Context, id uuid.UUID) (RefreshSession, error) {
	var r RefreshSession

	if err := s.DB(ctx).First(&r, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return r, ErrRefreshSessionNotFound
		default:
			return r, fmt.Errorf("query refresh session: %w", err)
		}
	}

	return r, nil
}

func (s *GormRefreshStore) Insert(ctx context.Context, refreshSession RefreshSession) (RefreshSession, error) {
	start := time.Now()
	if err := s.DB(ctx).Create(&refreshSession).Error; err != nil {
		return RefreshSession{}, fmt.Errorf("create refresh session: %w", err)
	}
	s.logger.Debug("insert refresh session", "elapsed", time.Since(start).String())

	return refreshSession, nil
}

func (s *GormRefreshStore) Update(ctx context.Context, refreshSession RefreshSession) error {
	// TODO change to Updates()
	if err := s.DB(ctx).Save(&refreshSession).Error; err != nil {
		return fmt.Errorf("save refresh session: %w", err)
	}
	return nil
}
