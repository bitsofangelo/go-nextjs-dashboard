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
	"go-nextjs-dashboard/internal/optional"
)

type invoiceModel struct {
	ID         uuid.UUID
	CustomerID optional.Optional[uuid.UUID]
	Amount     float64
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
		CustomerID: optional.FromPtr(i.CustomerID),
		Amount:     i.Amount,
		Status:     i.Status,
		Date:       i.Date,
	}
}

func toEntity(i invoiceModel) Invoice {
	return Invoice{
		ID:         i.ID,
		CustomerID: i.CustomerID.Ptr(),
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

func (s *GormStore) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	tx := s.DB(ctx).Model(&invoiceModel{}).Where("id = ?", id)

	exists, err := db.RecordExists(tx)
	if err != nil {
		return false, fmt.Errorf("query invoice exists: %w", err)
	}
	return exists, nil
}

func (s *GormStore) Insert(ctx context.Context, i Invoice) (*Invoice, error) {
	invModel := toModel(i)

	if err := s.DB(ctx).Create(&invModel).Error; err != nil {
		return nil, fmt.Errorf("store invoice: %w", err)
	}

	i = toEntity(invModel)
	return &i, nil
}

func (s *GormStore) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) error {
	err := s.DB(ctx).
		Model(&invoiceModel{}).
		Where("id = ?", id).
		Updates(req).Error

	if err != nil {
		return fmt.Errorf("update invoice: %w", err)
	}

	return nil
}
