package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/customer"
)

type Service struct {
	repo   customer.Store
	logger *slog.Logger
}

func New(repo customer.Store, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) List(ctx context.Context) ([]customer.Customer, error) {
	customers, err := s.repo.List(ctx)

	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}

	return customers, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error) {
	c, err := s.repo.Find(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get customer by id: %w", err)
	}

	return c, nil
}

func (s *Service) Create(ctx context.Context, c customer.Customer) (*customer.Customer, error) {
	exists, err := s.repo.ExistsByEmail(ctx, c.Email)
	if err != nil {
		return nil, fmt.Errorf("exists by email: %w", err)
	}

	if exists {
		s.logger.WarnContext(ctx, "email already taken", slog.String("email", c.Email))
		return nil, customer.ErrEmailAlreadyTaken
	}

	var cust *customer.Customer

	if cust, err = s.repo.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("save customer: %w", err)
	}
	return cust, nil
}

func (s *Service) SearchWithInvoiceTotals(ctx context.Context, search string) ([]customer.WithInvoiceTotals, error) {
	result, err := s.repo.SearchWithInvoiceTotals(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("search with invoice totals: %w", err)
	}

	return result, nil
}
