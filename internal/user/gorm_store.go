package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-dash/internal/db"
	"go-dash/internal/logger"
)

type userModel struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);not nullable;unique;primary_key"`
	Name     string    `json:"name" gorm:"type:varchar(255);not nullable"`
	Email    string    `json:"email" gorm:"type:varchar(255);not nullable;unique"`
	Password string    `json:"-" gorm:"type:text;not nullable"`
}

type GormStore struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ Store = (*GormStore)(nil)

func NewStore(db *gorm.DB, log logger.Logger) *GormStore {
	return &GormStore{
		db:     db,
		logger: log.With("component", "store.gorm.user"),
	}
}

func (s *GormStore) DB(ctx context.Context) *gorm.DB {
	if gormDB, ok := db.FromCtx(ctx); ok {
		return gormDB
	}
	return s.db
}

func (s *GormStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User

	if err := s.DB(ctx).First(&u, "email = ?", email).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("query by email: %w", err)
		}
	}

	return &u, nil
}
