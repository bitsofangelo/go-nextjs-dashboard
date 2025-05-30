package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"

	"go-dash/internal/config"
	"go-dash/internal/event/bus"
	"go-dash/internal/http"
	"go-dash/internal/logger"
)

type Server interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	logger logger.Logger
	server Server
}

func NewApp(
	cfg *config.Config,
	db *gorm.DB,
	logger logger.Logger,
	server Server,
	_ timezoneInitializer,
	_ bus.RegisterInitializer,
	_ http.RouteInitializer,
) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
		server: server,
	}
}

func (a *App) Run() error {
	// handle signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srvErr := make(chan error)

	go func() {
		if err := a.server.Serve(); err != nil {
			srvErr <- fmt.Errorf("start server: %w", err)
		}
	}()

	select {
	case err := <-srvErr:
		return err
	case <-ctx.Done(): // block until shutdown signal is received
		a.logger.With("component", "server").Info("shutdown signal received")

		// give other goroutines time to finish (DB, jobs, etc.)
		shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// shutdown fiber
		if err := a.server.Shutdown(shutCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed, forcing close: %w", err)
		}
		return nil
	}
}

func (a *App) Logger() logger.Logger {
	return a.logger
}

func (a *App) Close() error {
	if err := a.logger.Close(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}

	if db, err := a.db.DB(); err == nil {
		if err = db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}

	return nil
}
