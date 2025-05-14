package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-nextjs-dashboard/internal/bootstrap"
	"go-nextjs-dashboard/internal/server"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = app.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// handle signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// build and run server
	srv := server.New(ctx, app)
	if err = srv.Run(); err != nil && !errors.Is(err, context.Canceled) {
		app.Logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	app.Logger.Info("server exited")
}
