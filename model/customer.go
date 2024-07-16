package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID       uuid.UUID `gorm:"type:char(36);not null;unique;primary_key"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"type:varchar(255);not null;unique"`
	ImageURL string    `gorm:"type:varchar(255)"`
	Invoices []Invoice
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	customer.ID = uuid.New()
	return
}
