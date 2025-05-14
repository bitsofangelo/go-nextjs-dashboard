package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"go-nextjs-dashboard/internal/bootstrap"
	customerhttp "go-nextjs-dashboard/internal/customer/http"
	dashboardhttp "go-nextjs-dashboard/internal/dashboard/http"
	"go-nextjs-dashboard/internal/http"
	"go-nextjs-dashboard/internal/logger"
	userhttp "go-nextjs-dashboard/internal/user/http"
)

type Server struct {
	app    *bootstrap.App
	router *fiber.App
	ctx    context.Context
}

func New(
	ctx context.Context,
	app *bootstrap.App,
) *Server {
	router := fiber.New(fiber.Config{
		ErrorHandler: errHandler(app.Logger.With("component", "http")),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// global middlewares
	// router.Use(logger.New())
	router.Use(http.RequestID())
	router.Use(cors.New(cors.Config{}))
	router.Use(limiter.New(limiter.Config{Max: 10}))
	router.Use(http.ValidationResponse())
	router.Use(recover.New(recover.Config{}))

	// routes registration
	api := router.Group("/api")
	customerhttp.RegisterHTTP(api, app.CustSvc, app.Logger.With("component", "http.customer"))
	userhttp.RegisterHTTP(api, app.UserSvc, app.Logger.With("component", "http.user"))
	dashboardhttp.RegisterHTTP(api, app.DashSvc, app.Logger.With("component", "http.dashboard"))

	return &Server{
		app:    app,
		router: router,
		ctx:    ctx,
	}
}

func (s *Server) Run() error {
	serverLogger := s.app.Logger.With("component", "server")
	srvErr := make(chan error)

	// run fiber in goroutine
	go func() {
		if err := s.router.Listen(":" + s.app.Config.AppPort); err != nil {
			srvErr <- fmt.Errorf("start server: %w", err)
		}
	}()

	select {
	case err := <-srvErr:
		return err
	case <-s.ctx.Done(): // block until shutdown signal is received
		serverLogger.Info("shutdown signal received")

		// give other goroutines time to finish (DB, jobs, etc.)
		shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// graceful shutdown
		if err := s.router.ShutdownWithContext(shutCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed, forcing close: %w", err)
		}
		return nil
	}
}

func errHandler(logger logger.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		var fe *fiber.Error
		if errors.As(err, &fe) {
			if fe.Code != code {
				code = fe.Code
				message = fe.Message
			}
		}

		if code >= fiber.StatusInternalServerError {
			logger.ErrorContext(c.Context(), err.Error())
		}

		return c.Status(code).JSON(http.ErrResponse{
			Message: message,
		})
	}
}
