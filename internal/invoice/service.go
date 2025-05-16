package invoice

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/logger"
)

type Service struct {
	store  Store
	logger logger.Logger
}

func NewService(store Store, logger logger.Logger) *Service {
	return &Service{
		store:  store,
		logger: logger.With("component", "service.invoice"),
	}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	inv, err := s.store.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find invoice: %w", err)
	}
	return inv, nil
}

// func (s *Service) Create(ctx context.Context, inv invoice.Invoice) (*invoice.Invoice, error) {
// 	customerExists, err := s.custStore.Exists(ctx, inv.CustomerID)
// 	if err != nil {
//
// 	}
// }
