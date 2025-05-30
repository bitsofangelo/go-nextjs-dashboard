//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
)

func InitApp() (*App, error) {
	wire.Build(
		AppProviders,
		HTTPProviders,
		NewApp,
	)
	return nil, nil
}
