package dashboard

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
)

type GormStore struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ Store = (*GormStore)(nil)

func NewStore(db *gorm.DB, logger logger.Logger) *GormStore {
	return &GormStore{
		db:     db,
		logger: logger.With("component", "store.gorm.dash"),
	}
}

func (s GormStore) GetOverview(ctx context.Context) (*Overview, error) {
	var (
		invoiceCount  int64
		customerCount int64
		invoiceStatus InvoiceStatus
	)

	g, egCtx := errgroup.WithContext(ctx)

	start := time.Now()

	g.Go(func() error {
		if err := s.db.WithContext(egCtx).Model(&invoice.Invoice{}).Count(&invoiceCount).Error; err != nil {
			return fmt.Errorf("query invoice count: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := s.db.WithContext(egCtx).Model(&customer.Customer{}).Count(&customerCount).Error; err != nil {
			return fmt.Errorf("query customer count: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		err := s.db.WithContext(egCtx).Model(&invoice.Invoice{}).Select(`
			SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END) AS "paid",
			SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END) AS "pending"
		`).Scan(&invoiceStatus).Error

		if err != nil {
			return fmt.Errorf("query invoice status: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("query overview: %w", err)
	}

	s.logger.DebugContext(ctx, "fetch overview", "elapsed", time.Since(start).String())

	return &Overview{
		InvoiceCount:  invoiceCount,
		CustomerCount: customerCount,
		InvoiceStatus: invoiceStatus,
	}, nil
}
