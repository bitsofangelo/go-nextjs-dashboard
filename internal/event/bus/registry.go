package bus

import (
	"context"
	"fmt"
	"reflect"

	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/event"
)

func RegisterAll() map[reflect.Type]event.Publisher {
	buses := make(map[reflect.Type]event.Publisher)

	custCreatedBus := newBus[customer.Created]()
	custCreatedT := reflect.TypeOf(customer.Created{})
	{
		custCreatedBus.Subscribe(SendWelcomeMessage)
		buses[custCreatedT] = custCreatedBus
	}

	return buses
}

func SendWelcomeMessage(_ context.Context, e customer.Created) error {
	go func() {
		fmt.Println("Welcome message received", e.ID)
	}()
	return nil
}
