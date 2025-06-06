//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
)

func InitApp() (*App, error) {
	wire.Build(
		AppServiceProviders,
		HTTPProviders,
		AppProvider,
	)
	return nil, nil
}
