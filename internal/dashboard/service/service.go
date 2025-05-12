package service

import (
	"context"
	"fmt"

	"go-nextjs-dashboard/internal/dashboard"
	"go-nextjs-dashboard/internal/logger"
)

type Service struct {
	store  dashboard.Store
	logger logger.Logger
}

func New(store dashboard.Store, log logger.Logger) *Service {
	return &Service{
		store:  store,
		logger: log,
	}
}

func (s *Service) GetOverview(ctx context.Context) (*dashboard.Overview, error) {
	o, err := s.store.GetOverview(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieve overview: %w", err)
	}
	return o, nil
}
