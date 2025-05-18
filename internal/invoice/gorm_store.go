package invoice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/logger"
)

type invoiceModel struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Amount     float32
	Status     string
	Date       time.Time
}

func (i *invoiceModel) BeforeCreate(*gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}

	return
}

func (*invoiceModel) TableName() string {
	return "invoices"
}

func toModel(i Invoice) invoiceModel {
	return invoiceModel{
		ID:         i.ID,
		CustomerID: i.CustomerID,
		Amount:     i.Amount,
		Status:     i.Status,
		Date:       i.Date,
	}
}

func toEntity(i invoiceModel) Invoice {
	return Invoice{
		ID:         i.ID,
		CustomerID: i.CustomerID,
		Amount:     i.Amount,
		Status:     i.Status,
		Date:       i.Date,
	}
}

type GormStore struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ Store = (*GormStore)(nil)

func NewStore(db *gorm.DB, logger logger.Logger) *GormStore {
	return &GormStore{
		db:     db,
		logger: logger.With("component", "store.gorm.invoice"),
	}
}

func (s *GormStore) DB(ctx context.Context) *gorm.DB {
	if gormDB, ok := db.FromCtx(ctx); ok {
		return gormDB.WithContext(ctx)
	}
	return s.db.WithContext(ctx)
}

func (s *GormStore) Find(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	var i invoiceModel

	if err := s.DB(ctx).First(&i, "id = ?", id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrInvoiceNotFound
		default:
			return nil, fmt.Errorf("query invoice: %w", err)
		}
	}

	inv := toEntity(i)
	return &inv, nil
}

func (s *GormStore) Save(ctx context.Context, i Invoice) (*Invoice, error) {
	invModel := toModel(i)

	if err := s.DB(ctx).Create(&invModel).Error; err != nil {
		return nil, fmt.Errorf("store invoice: %w", err)
	}

	i = toEntity(invModel)
	return &i, nil
}
