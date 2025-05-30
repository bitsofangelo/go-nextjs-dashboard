package bootstrap

import (
	"fmt"
	"time"

	"github.com/google/wire"

	"github.com/gelozr/go-dash/internal/app"
	"github.com/gelozr/go-dash/internal/auth"
	"github.com/gelozr/go-dash/internal/config"
	"github.com/gelozr/go-dash/internal/customer"
	"github.com/gelozr/go-dash/internal/dashboard"
	"github.com/gelozr/go-dash/internal/db"
	"github.com/gelozr/go-dash/internal/event"
	"github.com/gelozr/go-dash/internal/event/bus"
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

	// MAIL
	mail.NewSMTPMailer,
	mail.NewManager,
	wire.Bind(new(mail.Sender), new(*mail.Manager)),

	// AUTH
	AuthDBProvider,
	wire.Bind(new(auth.Authenticator[user.User]), new(*auth.DBProvider[user.User])),
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

func AuthDBProvider(userSvc *user.Service, hash *hashing.Hash) *auth.DBProvider[user.User] {
	return auth.NewDBProvider[user.User](userSvc, hash)
}
