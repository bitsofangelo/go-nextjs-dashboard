package invoice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/listing"
	"go-nextjs-dashboard/internal/optional"
)

var ErrInvoiceNotFound = errors.New("invoice not found")

type Store interface {
	List(context.Context, listing.SortOrder) ([]Invoice, error)
	Search(context.Context, SearchFilter, listing.Page) ([]Invoice, int64, error)
	Find(context.Context, uuid.UUID) (*Invoice, error)
	Exists(context.Context, uuid.UUID) (bool, error)
	Insert(context.Context, Invoice) (*Invoice, error)
	Update(context.Context, uuid.UUID, UpdateInput) error
	Delete(context.Context, uuid.UUID) error

	ListWithCustomerInfo(context.Context, listing.SortOrder) ([]WithCustomerInfo, error)
}

type SearchFilter struct {
	Text string
	Sort listing.SortOrder
}

type UpdateInput struct {
	CustomerID optional.Optional[uuid.UUID]
	Amount     float64
	Status     string
	Date       optional.Optional[time.Time]
	IsActive   optional.Optional[bool]
}

type WithCustomerInfo struct {
	Invoice
	CustomerName     string
	CustomerEmail    string
	CustomerImageURL string
}
