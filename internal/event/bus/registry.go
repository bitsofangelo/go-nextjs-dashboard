package bus

import (
	customerevent "go-nextjs-dashboard/internal/customer/event"
	"go-nextjs-dashboard/internal/event"
)

func RegisterAll() map[string]event.Publisher {
	buses := make(map[string]event.Publisher)

	custBus := newBus[customerevent.Created]()
	{
		custBus.Subscribe(customerevent.SendWelcomeMessage)
		buses[custBus.Name()] = custBus
	}

	return buses
}
