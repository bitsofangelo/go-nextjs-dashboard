package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-nextjs-dashboard/internal/config"
	customerstore "go-nextjs-dashboard/internal/customer/gormstore"
	customersvc "go-nextjs-dashboard/internal/customer/service"
	database "go-nextjs-dashboard/internal/db"
	sloglogger "go-nextjs-dashboard/internal/logger/slog"
	"go-nextjs-dashboard/internal/server"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// set timezone
	loc, err := time.LoadLocation(cfg.AppTimezone)
	if err != nil {
		log.Fatalf("failed to set timezone: %v", err)
	}
	time.Local = loc

	// init logger
	logger, err := sloglogger.New(cfg)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func() {
		if err = logger.Close(); err != nil {
			logger.Error("failed to close logger", "error", err)
		}
	}()

	// open db
	db, err := database.Open(cfg)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}

	// wire dependencies
	custStore := customerstore.New(db, logger.With("component", "store.gorm"))
	custSvc := customersvc.New(custStore, logger.With("component", "service.customer"))

	// handle signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// build and run server
	srv := server.New(ctx, cfg, logger.With("component", "http"), custSvc)
	if err = srv.Run(); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited")
}
