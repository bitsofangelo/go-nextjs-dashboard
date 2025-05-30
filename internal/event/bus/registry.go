package bus

import (
	"context"
	"fmt"

	"github.com/gelozr/go-dash/internal/customer"
	"github.com/gelozr/go-dash/internal/event"
	"github.com/gelozr/go-dash/internal/logger"
	"github.com/gelozr/go-dash/internal/mail"
)

type RegisterInitializer struct{}

func RegisterAll(
	broker *event.Broker,
	custSvc *customer.Service,
	mailer mail.Sender,
	logger logger.Logger,
) RegisterInitializer {
	log := logger.With("component", "event/bus")

	custCreatedBus := newBus[customer.Created]()
	{
		custCreatedBus.Subscribe(SendWelcomeMessage(custSvc, mailer, log))
		custCreatedBus.Subscribe(SendVerifyEmailMessage(custSvc, log))
		broker.RegisterBus(custCreatedBus)
	}

	return RegisterInitializer{}
}

func SendWelcomeMessage(custSvc *customer.Service, mailer mail.Sender, log logger.Logger) Handler[customer.Created] {
	return func(ctx context.Context, e customer.Created) error {
		log = log.With("handler", "SendWelcomeMessage")

		safeGo(ctx, log, func() {
			cust, err := custSvc.GetByID(ctx, e.ID)
			if err != nil {
				log.Error("get customer by id", "error", err.Error())
				return
			}

			m := &mail.Message{
				From: mail.Address{
					Name:    "angelo",
					Address: "angelo@gwapo.dev",
				},
				To: []mail.Address{{
					Name:    cust.Name,
					Address: cust.Email,
				}},
				Subject: "Welcome subject",
				HTML:    "<p>Welcome message</p>",
				Text:    "Welcome message",
				Headers: nil,
			}

			if err = mailer.Send(ctx, m); err != nil {
				log.Error("send welcome email", "error", err.Error())
				return
			}

			fmt.Println("Welcome message sent", e.ID, cust.Name)
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
