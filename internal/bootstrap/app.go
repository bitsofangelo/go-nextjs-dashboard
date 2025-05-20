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

	return InitializeApp(cfg, logger, db)
}

func NewApp(cfg *config.Config,
	logger *sloglogger.Logger,

	// domain services
	custSvc *customer.Service,
	userSvc *user.Service,
	dashSvc *dashboard.Service,
	invSvc *invoice.Service,

	// use cases
	createInvoice *app.CreateInvoice,
) (*App, error) {
	return &App{
		Config: cfg,
		Logger: logger,

		// domain services
		CustSvc: custSvc,
		UserSvc: userSvc,
		DashSvc: dashSvc,
		InvSvc:  invSvc,

		// use cases
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
