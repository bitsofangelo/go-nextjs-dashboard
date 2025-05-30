package user

import (
	"context"
	"fmt"

	"github.com/gelozr/go-dash/internal/logger"
)

type Service struct {
	store  Store
	logger logger.Logger
}

func NewService(store Store, log logger.Logger) *Service {
	return &Service{
		store:  store,
		logger: log.With("component", "service.user"),
	}
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	u, err := s.store.FindByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return u, nil
}
