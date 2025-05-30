package bootstrap

import (
	"fmt"
	"time"

	"github.com/google/wire"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/auth"
	"go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/dashboard"
	"go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/event"
	"go-nextjs-dashboard/internal/event/bus"
	"go-nextjs-dashboard/internal/hashing"
	"go-nextjs-dashboard/internal/http"
	"go-nextjs-dashboard/internal/http/validation"
	"go-nextjs-dashboard/internal/http/validation/gp"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/logger/slog"
	"go-nextjs-dashboard/internal/user"
)

var AppProviders = wire.NewSet(
	// CONFIG
	config.Load,
	setTimezone,

	// DB
	db.Open,
	db.NewTxManager,
	wire.Bind(new(db.TxManager), new(*db.GormTxManager)),

	// LOGGER
	slog.New,
	wire.Bind(new(logger.Logger), new(*slog.Logger)),

	// EVENT
	event.NewBroker,
	wire.Bind(new(event.Publisher), new(*event.Broker)),
	bus.RegisterAll,

	// VALIDATOR
	gp.New,
	wire.Bind(new(validation.Validator), new(*gp.Validator)),

	// HASHING
	hashing.NewArgon2IDHasher,
	wire.Bind(new(hashing.Hasher), new(*hashing.Argon2IDHasher)),
	hashing.New,

	// AUTH
	auth.NewPasswordProvider,
	auth.NewGoogleProvider,
	auth.NewManager,
	wire.Bind(new(auth.Auth), new(*auth.Manager)),
	auth.NewGormRefreshStore,
	wire.Bind(new(auth.RefreshStore), new(*auth.GormRefreshStore)),
	auth.NewGOJWT,
	wire.Bind(new(auth.JWT), new(*auth.GOJWT)),
	auth.NewToken,

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
	app.NewAuthenticateUser,
	app.NewCreateInvoice,
)

var HTTPProviders = wire.NewSet(
	// HANDLERS
	http.NewAuthHandler,
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
