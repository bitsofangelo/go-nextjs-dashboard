package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/customer"
	customerevent "go-nextjs-dashboard/internal/customer/event"
	"go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/event"
	"go-nextjs-dashboard/internal/logger"
)

type Service struct {
	store  customer.Store
	txm    db.GormTxManager
	event  event.Publisher
	logger logger.Logger
}

func New(store customer.Store, txm db.GormTxManager, evt event.Publisher, log logger.Logger) *Service {
	return &Service{
		store:  store,
		txm:    txm,
		event:  evt,
		logger: log,
	}
}

func (s *Service) List(ctx context.Context) ([]customer.Customer, error) {
	customers, err := s.store.List(ctx)

	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}

	return customers, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error) {
	c, err := s.store.Find(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get customer by id: %w", err)
	}

	return c, nil
}

func (s *Service) Create(ctx context.Context, c customer.Customer) (*customer.Customer, error) {
	exists, err := s.store.ExistsByEmail(ctx, c.Email)
	if err != nil {
		return nil, fmt.Errorf("exists by email: %w", err)
	}

	if exists {
		s.logger.WarnContext(ctx, "email already taken", "email", c.Email)
		return nil, customer.ErrEmailAlreadyTaken
	}

	var cust *customer.Customer

	err = s.txm.Do(ctx, func(txCtx context.Context) error {
		if cust, err = s.store.Save(txCtx, c); err != nil {
			return fmt.Errorf("save customer: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err = s.event.Publish(ctx, customerevent.Created{ID: cust.ID}); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}

	return cust, nil
}

func (s *Service) SearchWithInvoiceTotals(ctx context.Context, search string) ([]customer.WithInvoiceTotals, error) {
	result, err := s.store.SearchWithInvoiceTotals(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("search with invoice totals: %w", err)
	}

	return result, nil
}
