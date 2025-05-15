package gormstore

import (
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

// var _ invoice.Store = (*Store)(nil)

// func (s *Store) FindLatestInvoices(ctx context.Context, limit int) ([]invoice.Invoice, error) {
// 	var invoices []model.Invoice
//
// 	err := s.db.Preload("Customer").Order("date DESC").Limit(limit).Find(&invoices).Error
//
// 	return invoices, err
// }

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}
