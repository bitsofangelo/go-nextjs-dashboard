package bus

import (
	"reflect"

	customerevent "go-nextjs-dashboard/internal/customer/event"
	"go-nextjs-dashboard/internal/event"
)

func RegisterAll() map[reflect.Type]event.Publisher {
	buses := make(map[reflect.Type]event.Publisher)

	custCreatedBus := newBus[customerevent.Created]()
	custCreatedT := reflect.TypeOf(customerevent.Created{})
	{
		custCreatedBus.Subscribe(customerevent.SendWelcomeMessage)
		buses[custCreatedT] = custCreatedBus
	}

	return buses
}
