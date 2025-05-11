package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/internal/customer/gormstore"
	customerhttp "go-nextjs-dashboard/internal/customer/http"
	"go-nextjs-dashboard/internal/customer/service"
	database "go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/http"
)

func main() {
	// load config once
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("cannot open db: %v", err)
	}

	customerStore := gormstore.NewStore(db)
	customerService := service.NewService(customerStore)

	// Setup log file
	f, err := os.OpenFile("fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening log file: %w", err))
	}
	defer f.Close()

	fiberlog.SetOutput(f)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			var fe *fiber.Error
			if errors.As(err, &fe) {
				if fe.Code != code {
					code = fe.Code
					message = fe.Message
				}
			}

			if code == fiber.StatusInternalServerError {
				fiberlog.Error(err.Error())
			}

			return c.Status(code).JSON(fiber.Map{"message": message})
		},
	})

	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{Max: 10}))
	app.Use(http.ValidationResponse())

	api := app.Group("/api")
	customerhttp.RegisterHTTP(api, customerService)

	err = app.Listen(":" + cfg.AppPort)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
