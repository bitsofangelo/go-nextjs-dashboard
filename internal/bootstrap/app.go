package bootstrap

import (
	"context"
	"fmt"
	"time"

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
	logger logger.Logger
	server Server
	ctx    context.Context
}

func NewApp(
	ctx context.Context,
	cfg *config.Config,
	logger logger.Logger,
	server Server,
	_ timezoneInitializer,
	_ bus.RegisterInitializer,
	_ http.RouteInitializer,
) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
		server: server,
		ctx:    ctx,
	}
}

func (a *App) Run() error {
	srvErr := make(chan error)

	go func() {
		if err := a.server.Serve(); err != nil {
			srvErr <- fmt.Errorf("start server: %w", err)
		}
	}()

	select {
	case err := <-srvErr:
		return err
	case <-a.ctx.Done(): // block until shutdown signal is received
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
		a.logger.Error("failed to close logger", "error", err)
		return fmt.Errorf("failed to close logger: %w", err)
	}
	return nil
}
