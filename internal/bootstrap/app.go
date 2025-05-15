package bootstrap

import (
	"fmt"
	"os"
	"time"

	"go-nextjs-dashboard/internal/config"
	customerstore "go-nextjs-dashboard/internal/customer/gormstore"
	customersvc "go-nextjs-dashboard/internal/customer/service"
	dashboardstore "go-nextjs-dashboard/internal/dashboard/gormstore"
	dashboardsvc "go-nextjs-dashboard/internal/dashboard/service"
	database "go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/event"
	eventbus "go-nextjs-dashboard/internal/event/bus"
	sloglogger "go-nextjs-dashboard/internal/logger/slog"
	userstore "go-nextjs-dashboard/internal/user/gormstore"
	usersvc "go-nextjs-dashboard/internal/user/service"
)

type App struct {
	Config *config.Config
	Logger *sloglogger.Logger

	CustSvc *customersvc.Service
	UserSvc *usersvc.Service
	DashSvc *dashboardsvc.Service
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
	buses := eventbus.RegisterAll()
	eventBroker := event.NewBroker(buses)

	// wire dependencies
	gormStoreLogger := logger.With("component", "store.gorm")
	custStore := customerstore.New(db, gormStoreLogger)
	custSvc := customersvc.New(custStore, txm, eventBroker, logger.With("component", "service.customer"))
	userStore := userstore.New(db, gormStoreLogger)
	userSvc := usersvc.New(userStore, logger.With("component", "service.user"))
	dashStore := dashboardstore.New(db, gormStoreLogger)
	dashSvc := dashboardsvc.New(dashStore, logger.With("component", "service.dashboard"))

	return &App{
		Config:  cfg,
		Logger:  logger,
		CustSvc: custSvc,
		UserSvc: userSvc,
		DashSvc: dashSvc,
	}, nil
}

func (app *App) Close() error {
	if err := app.Logger.Close(); err != nil {
		app.Logger.Error("failed to close logger", "error", err)
		return fmt.Errorf("failed to close logger: %w", err)
	}
	return nil
}
