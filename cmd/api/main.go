package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gelozr/go-dash/internal/bootstrap"
)

func main() {
	// init app
	app, err := bootstrap.InitApp()
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
		app.Logger().Error("failed to start app", "error", err)
		os.Exit(1)
	}

	app.Logger().Info("app exited")
}
