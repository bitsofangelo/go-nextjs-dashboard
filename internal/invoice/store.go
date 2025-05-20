package invoice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/optional"
)

var ErrInvoiceNotFound = errors.New("invoice not found")

type Store interface {
	Find(context.Context, uuid.UUID) (*Invoice, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Insert(context.Context, Invoice) (*Invoice, error)
	Update(context.Context, uuid.UUID, *UpdateRequest) error
}

type UpdateRequest struct {
	CustomerID optional.Optional[uuid.UUID]
	Amount     float64
	Status     string
	Date       time.Time
}
