package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-nextjs-dashboard/internal/config"
	customerstore "go-nextjs-dashboard/internal/customer/gormstore"
	customersvc "go-nextjs-dashboard/internal/customer/service"
	database "go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/server"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// setup logger
	rootLog := logger.New(cfg)
	slog.SetDefault(rootLog) // optional global fallback

	// open db connection
	db, err := database.Open(cfg)
	if err != nil {
		rootLog.Error("failed to open database", "error", err)
		os.Exit(1)
	}

	// wire dependencies
	custStore := customerstore.New(db, rootLog.With("component", "gormstore"))
	custSvc := customersvc.New(custStore, rootLog.With("component", "customer-service"))

	// handle signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// instantiate server
	srv := server.New(ctx, cfg, rootLog.With("component", "http"), custSvc)

	// start server
	if err = srv.Run(); err != nil && !errors.Is(err, context.Canceled) {
		rootLog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	rootLog.Info("server exited")
}
