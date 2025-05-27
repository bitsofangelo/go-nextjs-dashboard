package bootstrap

import (
	"fmt"
	"time"

	"github.com/google/wire"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/dashboard"
	database "go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/event"
	"go-nextjs-dashboard/internal/event/bus"
	"go-nextjs-dashboard/internal/http"
	"go-nextjs-dashboard/internal/http/validation"
	"go-nextjs-dashboard/internal/http/validation/gp"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
	sloglogger "go-nextjs-dashboard/internal/logger/slog"
	"go-nextjs-dashboard/internal/user"
)

var AppProviders = wire.NewSet(
	// CONFIG
	config.Load,
	setTimezone,

	// DB
	database.Open,
	database.NewTxManager,
	wire.Bind(new(database.TxManager), new(*database.GormTxManager)),

	// LOGGER
	sloglogger.New,
	wire.Bind(new(logger.Logger), new(*sloglogger.Logger)),

	// EVENT
	event.NewBroker,
	wire.Bind(new(event.Publisher), new(*event.Broker)),
	bus.RegisterAll,

	// VALIDATOR
	gp.New,
	wire.Bind(new(validation.Validator), new(*gp.Validator)),

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
)

var HTTPProviders = wire.NewSet(
	// HANDLERS
	http.NewDashboardHandler,
	http.NewUserHandler,
	http.NewCustomerHandler,
	http.NewInvoiceHandler,

	// ENGINE
	http.NewFiberServer,
	wire.Bind(new(Server), new(*http.FiberServer)),

	// ROUTES
	http.SetupFiberRoutes,
)

type timezoneInitializer struct{}

func setTimezone(cfg *config.Config) (timezoneInitializer, error) {
	loc, err := time.LoadLocation(cfg.AppTimezone)
	if err != nil {
		return timezoneInitializer{}, fmt.Errorf("failed to load timezone: %w", err)
	}
	time.Local = loc

	return timezoneInitializer{}, nil
}
