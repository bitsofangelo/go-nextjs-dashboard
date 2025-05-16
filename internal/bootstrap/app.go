package bootstrap

import (
	"fmt"
	"os"
	"time"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/dashboard"
	database "go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/event"
	"go-nextjs-dashboard/internal/event/bus"
	"go-nextjs-dashboard/internal/invoice"
	sloglogger "go-nextjs-dashboard/internal/logger/slog"
	"go-nextjs-dashboard/internal/user"
)

type App struct {
	Config *config.Config
	Logger *sloglogger.Logger

	// services
	CustSvc *customer.Service
	UserSvc *user.Service
	DashSvc *dashboard.Service
	InvSvc  *invoice.Service

	// use cases
	CreateInvoice *app.CreateInvoice
}

func New() (*App, error) {
	// load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// set timezone
	loc, err := time.LoadLocation(cfg.AppTimezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}
	time.Local = loc

	// init logger
	logger, err := sloglogger.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// open db
	db, err := database.Open(cfg, logger.With("component", "db"))
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	txm := database.NewTxManager(db)

	// init events and handlers
	buses := bus.RegisterAll()
	eventBroker := event.NewBroker(buses)

	// wire dependencies
	custStore := customer.NewStore(db, logger)
	custSvc := customer.NewService(custStore, txm, eventBroker, logger)
	userStore := user.NewStore(db, logger)
	userSvc := user.NewService(userStore, logger)
	dashStore := dashboard.NewStore(db, logger)
	dashSvc := dashboard.NewService(dashStore, logger)
	invStore := invoice.NewStore(db, logger)
	invSvc := invoice.NewService(invStore, logger)

	// use cases
	createInvoice := app.NewCreateInvoice(custStore, invStore, txm, logger)

	return &App{
		Config:  cfg,
		Logger:  logger,
		CustSvc: custSvc,
		UserSvc: userSvc,
		DashSvc: dashSvc,
		InvSvc:  invSvc,

		CreateInvoice: createInvoice,
	}, nil
}

func (app *App) Close() error {
	if err := app.Logger.Close(); err != nil {
		app.Logger.Error("failed to close logger", "error", err)
		return fmt.Errorf("failed to close logger: %w", err)
	}
	return nil
}
