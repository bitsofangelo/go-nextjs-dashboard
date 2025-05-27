package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-nextjs-dashboard/internal/logger"
)

type userModel struct {
	ID       uuid.UUID `json:"id" gormstore:"type:char(36);not nullable;unique;primary_key"`
	Name     string    `json:"name" gormstore:"type:varchar(255);not nullable"`
	Email    string    `json:"email" gormstore:"type:varchar(255);not nullable;unique"`
	Password string    `json:"-" gormstore:"type:text;not nullable"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// user.ID = uuid.New()
	// hashedPassword, err := auth.HashPassword(user.Password)
	// if err != nil {
	// 	return
	// }
	// user.Password = hashedPassword
	return
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

func (s GormStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User

	if err := s.db.First(&u, "email = ?", email).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("query by email: %w", err)
		}
	}

	return &u, nil
}
