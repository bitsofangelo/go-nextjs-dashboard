package bootstrap

import (
	"fmt"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/gelozr/go-dash/internal/app"
	"github.com/gelozr/go-dash/internal/auth"
	"github.com/gelozr/go-dash/internal/config"
	"github.com/gelozr/go-dash/internal/customer"
	"github.com/gelozr/go-dash/internal/dashboard"
	"github.com/gelozr/go-dash/internal/db"
	"github.com/gelozr/go-dash/internal/event"
	"github.com/gelozr/go-dash/internal/event/registry"
	"github.com/gelozr/go-dash/internal/hashing"
	"github.com/gelozr/go-dash/internal/http"
	"github.com/gelozr/go-dash/internal/http/validation"
	"github.com/gelozr/go-dash/internal/http/validation/gp"
	"github.com/gelozr/go-dash/internal/invoice"
	"github.com/gelozr/go-dash/internal/logger"
	"github.com/gelozr/go-dash/internal/logger/slog"
	"github.com/gelozr/go-dash/internal/mail"
	"github.com/gelozr/go-dash/internal/user"
)

var AppServiceProviders = wire.NewSet(
	// CONFIG
	config.Load,

	// LOGGER
	slog.New,
	wire.Bind(new(logger.Logger), new(*slog.Logger)),

	// DB
	db.Open,
	db.NewTxManager,
	wire.Bind(new(db.TxManager), new(*db.GormTxManager)),

	// EVENT
	event.NewBroker,
	wire.Bind(new(event.Publisher), new(*event.Broker)),
	registry.RegisterAll,

	// VALIDATOR
	gp.New,
	wire.Bind(new(validation.Validator), new(*gp.Validator)),

	// HASHING
	hashing.NewManager,

	// MAIL
	mail.NewManager,
	wire.Bind(new(mail.Mailer), new(mail.Manager)),

	// AUTH
	auth.NewGormRefreshStore,
	wire.Bind(new(auth.RefreshStore), new(*auth.GormRefreshStore)),
	auth.NewToken,
	auth.NewDBUserProvider,
	auth.NewJWTDriver,
	AuthProvider,
	wire.Bind(new(auth.Authenticator), new(*auth.Provider)),
	wire.Bind(new(auth.LoginHandler), new(*auth.Provider)),
	wire.Bind(new(auth.LogoutHandler), new(*auth.Provider)),
	wire.Bind(new(auth.Checker), new(*auth.Provider)),
	wire.Bind(new(auth.TokenRefresher), new(*auth.Provider)),
	wire.Bind(new(auth.Auth), new(*auth.Provider)),

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
	http.NewAuthHandler,
	http.NewDashboardHandler,
	http.NewUserHandler,
	http.NewCustomerHandler,
	http.NewInvoiceHandler,

	// ENGINE
	http.NewFiberServer,
	// wire.Bind(new(Server), new(*http.FiberServer)),

	// ROUTES
	http.SetupFiberRoutes,
)

func AuthProvider(dbUserProvider *auth.DBUserProvider, jwtDriver *auth.JWTDriver) (*auth.Provider, error) {
	a := auth.New()

	if err := a.Extend("jwt", auth.GuardOption{Driver: jwtDriver, UserProvider: dbUserProvider}); err != nil {
		return nil, fmt.Errorf("auth extend: %w", err)
	}

	if err := a.SetDefaultGuard("jwt"); err != nil {
		return nil, fmt.Errorf("set default guard: %w", err)
	}

	return a, nil
}

func AppProvider(
	cfg *config.Config,
	db *gorm.DB,
	logger logger.Logger,
	fiberServer *http.FiberServer,
	_ registry.RegisterInitializer,
	_ http.RouteInitializer,
) (*App, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to sql db: %w", err)
	}

	return NewApp(cfg, sqlDB, logger, fiberServer), nil
}
