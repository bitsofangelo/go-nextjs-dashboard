package gormstore

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/user"
)

type Store struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ user.Store = (*Store)(nil)

func New(db *gorm.DB, log logger.Logger) *Store {
	return &Store{
		db:     db,
		logger: log,
	}
}

func (s Store) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User

	if err := s.db.First(&u, "email = ?", email).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, user.ErrUserNotFound
		default:
			return nil, fmt.Errorf("query by email: %w", err)
		}
	}

	return &u, nil
}
