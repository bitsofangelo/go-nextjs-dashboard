package gormstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/logger"
)

type customerModel struct {
	ID       uuid.UUID `gorm:"type:char(36);not null;unique;primary_key"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"type:varchar(255);not null;unique"`
	ImageURL *string   `gorm:"type:varchar(255)"`
}

func (c *customerModel) BeforeCreate(_ *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	return
}

func (c *customerModel) TableName() string {
	return "customers"
}

func toModel(c customer.Customer) customerModel {
	return customerModel{
		ID:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		ImageURL: c.ImageURL,
	}
}

func toEntity(c customerModel) customer.Customer {
	return customer.Customer{
		ID:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		ImageURL: c.ImageURL,
	}
}

type Store struct {
	db     *gorm.DB
	logger logger.Logger
}

// compile‑time check that we satisfy the port
var _ customer.Store = (*Store)(nil)

func New(db *gorm.DB, logger logger.Logger) *Store {
	return &Store{
		db:     db,
		logger: logger,
	}
}

func (s *Store) List(ctx context.Context) ([]customer.Customer, error) {
	var models []customerModel

	if err := s.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("query customers: %w", err)
	}

	customers := make([]customer.Customer, len(models))

	for i, m := range models {
		customers[i] = toEntity(m)
	}

	return customers, nil
}

func (s *Store) Find(ctx context.Context, u uuid.UUID) (*customer.Customer, error) {
	var model customerModel

	if err := s.db.WithContext(ctx).First(&model, "id = ?", u).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, customer.ErrCustomerNotFound
		default:
			return nil, fmt.Errorf("query customer: %w", err)
		}
	}

	c := toEntity(model)
	return &c, nil
}

func (s *Store) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	tx := s.db.WithContext(ctx).
		Model(&customerModel{}).
		Where("email = ?", email)

	exists, err := db.RecordExists(tx)
	if err != nil {
		return false, fmt.Errorf("exists by email: %w", err)
	}

	return exists, nil
}

func (s *Store) Save(ctx context.Context, c customer.Customer) (*customer.Customer, error) {
	model := toModel(c)

	if err := s.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, fmt.Errorf("store customer: %w", err)
	}

	c = toEntity(model)

	return &c, nil
}

func (s *Store) SearchWithInvoiceTotals(ctx context.Context, search string) ([]customer.WithInvoiceTotals, error) {
	var out []customer.WithInvoiceTotals

	start := time.Now()

	err := s.db.WithContext(ctx).
		Model(&customerModel{}).
		Select(`
            customers.id,
            customers.name,
            customers.email,
            customers.image_url,
            COUNT(invoices.id) AS total_invoices,
            SUM(CASE WHEN invoices.status = 'pending' THEN invoices.amount ELSE 0 END) AS total_pending,
            SUM(CASE WHEN invoices.status = 'paid'    THEN invoices.amount ELSE 0 END) AS total_paid
        `).
		Joins("LEFT JOIN invoices ON customers.id = invoices.customer_id").
		Where("customers.name LIKE @s OR customers.email LIKE @s", sql.Named("s", "%"+search+"%")).
		Group("customers.id, customers.name, customers.email, customers.image_url").
		Order("customers.name").
		Scan(&out).Error

	if err != nil {
		return nil, fmt.Errorf("invoice totals query: %w", err)
	}

	s.logger.DebugContext(ctx, "fetch customer with invoice", "elapsed", time.Since(start).String())

	return out, nil
}
