package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID         uuid.UUID  `gormstore:"type:char(36);not null;unique;primary_key"`
	CustomerID *uuid.UUID `gormstore:"type:char(36)"`
	Amount     float32    `gormstore:"type:float;not null"`
	Status     string     `gormstore:"type:varchar(255);not null"`
	Date       time.Time  `gormstore:"not null"`
	Customer   *Customer
}

func (invoice *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	invoice.ID = uuid.New()
	return
}
