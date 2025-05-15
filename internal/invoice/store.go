package invoice

import "context"

type Store interface {
	FindLatestInvoices(ctx context.Context, limit int) ([]Invoice, error)
}
