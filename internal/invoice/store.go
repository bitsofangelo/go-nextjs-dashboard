package invoice

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrInvoiceNotFound = errors.New("invoice not found")

type Store interface {
	Find(context.Context, uuid.UUID) (*Invoice, error)
	Save(context.Context, Invoice) (*Invoice, error)
}
