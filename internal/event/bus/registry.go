package bus

import (
	"context"
	"fmt"

	"go-dash/internal/customer"
	"go-dash/internal/event"
	"go-dash/internal/logger"
)

type RegisterInitializer struct{}

func RegisterAll(
	broker *event.Broker,
	custSvc *customer.Service,
	logger logger.Logger,
) RegisterInitializer {
	log := logger.With("component", "event/bus")

	custCreatedBus := newBus[customer.Created]()
	{
		custCreatedBus.Subscribe(SendWelcomeMessage(custSvc, log))
		custCreatedBus.Subscribe(SendVerifyEmailMessage(custSvc, log))
		broker.RegisterBus(custCreatedBus)
	}

	return RegisterInitializer{}
}

func SendWelcomeMessage(custSvc *customer.Service, log logger.Logger) Handler[customer.Created] {
	return func(ctx context.Context, e customer.Created) error {
		log = log.With("handler", "SendWelcomeMessage")

		safeGo(ctx, log, func() {
			cust, err := custSvc.GetByID(ctx, e.ID)
			if err != nil {
				log.Error("get customer by id", "error", err.Error())
				return
			}

			fmt.Println("Welcome message received", e.ID, cust.Name)
		})

		return nil
	}
}

func SendVerifyEmailMessage(custSvc *customer.Service, log logger.Logger) Handler[customer.Created] {
	return func(ctx context.Context, e customer.Created) error {
		log = log.With("handler", "SendVerifyEmailMessage")

		safeGo(ctx, log, func() {
			cust, err := custSvc.GetByID(ctx, e.ID)
			if err != nil {
				log.Error("get customer by id", "error", err.Error())
				return
			}

			fmt.Println("Email verification sent", e.ID, cust.Name)
		})

		return nil
	}
}
