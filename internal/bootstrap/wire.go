//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/dashboard"
	database "go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/event"
	"go-nextjs-dashboard/internal/event/bus"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
	sloglogger "go-nextjs-dashboard/internal/logger/slog"
	"go-nextjs-dashboard/internal/user"
)

func InitializeApp(cfg *config.Config, slogLogger *sloglogger.Logger, db *gorm.DB) (*App, error) {
	wire.Build(
		// INFRA
		database.NewTxManager,
		wire.Bind(new(database.TxManager), new(*database.GormTxManager)),

		// LOGGER
		wire.Bind(new(logger.Logger), new(*sloglogger.Logger)),

		// EVENT
		bus.RegisterAll,
		event.NewBroker,
		wire.Bind(new(event.Publisher), new(*event.Broker)),

		// STORE & SERVICES
		customer.NewStore,
		wire.Bind(new(customer.Store), new(*customer.GormStore)),
		customer.NewService,

		user.NewStore,
		wire.Bind(new(user.Store), new(*user.GormStore)),
		user.NewService,

		dashboard.NewStore,
		wire.Bind(new(dashboard.Store), new(*dashboard.GormStore)),
		dashboard.NewService,

		invoice.NewStore,
		wire.Bind(new(invoice.Store), new(*invoice.GormStore)),
		invoice.NewService,

		// USE CASES
		app.NewCreateInvoice,

		// FINAL APP CONSTRUCTOR
		NewApp,
	)

	return nil, nil
}
