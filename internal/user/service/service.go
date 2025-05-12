package service

import (
	"context"
	"fmt"

	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/user"
)

type Service struct {
	store  user.Store
	logger logger.Logger
}

func New(store user.Store, log logger.Logger) *Service {
	return &Service{
		store:  store,
		logger: log,
	}
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u, err := s.store.FindByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return u, nil
}
