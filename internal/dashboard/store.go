package dashboard

import "context"

type Store interface {
	GetOverview(ctx context.Context) (*Overview, error)
}
