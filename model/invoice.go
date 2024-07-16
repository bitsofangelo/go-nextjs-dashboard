package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID         uuid.UUID  `gorm:"type:char(36);not null;unique;primary_key"`
	CustomerID *uuid.UUID `gorm:"type:char(36)"`
	Amount     float32    `gorm:"type:float;not null"`
	Status     string     `gorm:"type:varchar(255);not null"`
	Date       time.Time  `gorm:"not null"`
	Customer   *Customer
}

func (invoice *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	invoice.ID = uuid.New()
	return
}
