package service

import (
	"context"

	"github.com/google/uuid"
)

type InvoiceTotals struct {
	ID            uuid.UUID
	Name          string
	Email         string
	ImageURL      *string
	TotalInvoices int64
	TotalPending  int64
	TotalPaid     int64
}

type InvoiceTotalsReader interface {
	ListWithInvoiceTotals(ctx context.Context, search string) ([]InvoiceTotals, error)
}
