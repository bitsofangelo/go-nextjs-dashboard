package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-nextjs-dashboard/internal/bootstrap"
)

func main() {
	// handle signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// init app
	app, err := bootstrap.InitializeApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app %v", err)
	}
	defer func() {
		if err = app.Close(); err != nil {
			log.Fatalf("failed to clean up %v", err)
		}
	}()

	// run app
	if err = app.Run(); err != nil && !errors.Is(err, context.Canceled) {
		app.Logger().Error("failed to start server", "error", err)
		os.Exit(1)
	}

	app.Logger().Info("server exited")
}
