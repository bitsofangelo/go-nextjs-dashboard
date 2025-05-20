package bus

import (
	"context"
	"fmt"
	"reflect"

	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/event"
)

func RegisterAll() map[string]event.Publisher {
	buses := make(map[string]event.Publisher)

	custCreatedBus := newBus[customer.Created]()
	custCreatedKey := reflect.TypeOf(customer.Created{}).String()
	{
		custCreatedBus.Subscribe(SendWelcomeMessage)
		buses[custCreatedKey] = custCreatedBus
	}

	return buses
}

func SendWelcomeMessage(_ context.Context, e customer.Created) error {
	go func() {
		fmt.Println("Welcome message received", e.ID)
	}()
	return nil
}
