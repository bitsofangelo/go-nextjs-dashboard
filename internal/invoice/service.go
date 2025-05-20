package invoice

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/logger"
)

var ErrInvalidCustomerID = fmt.Errorf("invalid customer id")

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

func (s *Service) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := s.store.Exists(ctx, id)
	if err != nil {
		return false, fmt.Errorf("exists invoice: %w", err)
	}
	return exists, nil
}

func (s *Service) Create(ctx context.Context, inv Invoice) (*Invoice, error) {
	if inv.CustomerID == nil {
		return nil, ErrInvalidCustomerID
	}

	i, err := s.store.Insert(ctx, inv)
	if err != nil {
		return nil, fmt.Errorf("save invoice: %w", err)
	}

	return i, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*Invoice, error) {
	exists, err := s.Exists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("exists invoice: %w", err)
	}
	if !exists {
		return nil, ErrInvoiceNotFound
	}

	if err = s.store.Update(ctx, id, req); err != nil {
		return nil, fmt.Errorf("update invoice: %w", err)
	}

	inv, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get invoice: %w", err)
	}

	return inv, nil
}
