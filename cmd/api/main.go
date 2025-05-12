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
	dashboardstore "go-nextjs-dashboard/internal/dashboard/gormstore"
	dashboardservice "go-nextjs-dashboard/internal/dashboard/service"
	database "go-nextjs-dashboard/internal/db"
	sloglogger "go-nextjs-dashboard/internal/logger/slog"
	"go-nextjs-dashboard/internal/server"
	userstore "go-nextjs-dashboard/internal/user/gormstore"
	userservice "go-nextjs-dashboard/internal/user/service"
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
	db, err := database.Open(cfg, logger.With("component", "db"))
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}

	// wire dependencies
	gormStoreLogger := logger.With("component", "store.gorm")
	custStore := customerstore.New(db, gormStoreLogger)
	custSvc := customersvc.New(custStore, logger.With("component", "service.customer"))
	userStore := userstore.New(db, gormStoreLogger)
	userSvc := userservice.New(userStore, logger.With("component", "service.user"))
	dashStore := dashboardstore.New(db, gormStoreLogger)
	dashSvc := dashboardservice.New(dashStore, logger.With("component", "service.dashboard"))

	// handle signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// build and run server
	srv := server.New(ctx, cfg, logger.With("component", "http"), custSvc, userSvc, dashSvc)
	if err = srv.Run(); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited")
}
