//go:build wireinject
// +build wireinject

package bootstrap

import (
	"context"

	"github.com/google/wire"
)

func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(
		AppProviders,
		HTTPProviders,
		NewApp,
	)
	return nil, nil
}
