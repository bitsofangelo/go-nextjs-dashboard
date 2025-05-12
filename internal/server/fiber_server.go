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

	"go-nextjs-dashboard/internal/config"
	customerhttp "go-nextjs-dashboard/internal/customer/http"
	customersvc "go-nextjs-dashboard/internal/customer/service"
	"go-nextjs-dashboard/internal/http"
	"go-nextjs-dashboard/internal/logger"
)

type Server struct {
	app    *fiber.App
	ctx    context.Context
	cfg    *config.Config
	logger logger.Logger
}

func New(ctx context.Context, cfg *config.Config, logger logger.Logger, custSvc *customersvc.Service) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: errHandler(logger),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// global middlewares
	// app.Use(logger.New())
	app.Use(http.RequestID())
	app.Use(cors.New(cors.Config{}))
	app.Use(limiter.New(limiter.Config{Max: 10}))
	app.Use(http.ValidationResponse())
	app.Use(recover.New(recover.Config{}))

	// route registration
	api := app.Group("/api")
	customerhttp.RegisterHTTP(api, custSvc, logger)

	return &Server{
		app:    app,
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() error {
	srvErr := make(chan error)

	// run fiber in goroutine
	go func() {
		if err := s.app.Listen(":" + s.cfg.AppPort); err != nil {
			srvErr <- fmt.Errorf("start server: %w", err)
		}
	}()

	s.logger.Info(fmt.Sprintf("server is running at %s", s.cfg.AppPort))

	select {
	case err := <-srvErr:
		return err
	case <-s.ctx.Done(): // block until shutdown signal is received
		s.logger.Info("shutdown signal received")

		// give other goroutines time to finish (DB, jobs, etc.)
		shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// graceful shutdown
		if err := s.app.ShutdownWithContext(shutCtx); err != nil {
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
