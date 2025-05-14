package event

import (
	"context"
	"errors"
	"fmt"
)

type Mode int

const (
	ModeSync Mode = iota
	ModeAsync
)

type Event interface {
	Key() string
}

type Publisher interface {
	Publish(ctx context.Context, evt Event) error
}

type Broker struct {
	buses map[string]Publisher
}

func NewBroker(buses map[string]Publisher) *Broker {
	return &Broker{
		buses: buses,
	}
}

func (r *Broker) Publish(ctx context.Context, evt Event) error {
	bus, ok := r.buses[evt.Key()]
	if !ok {
		return errors.New(fmt.Sprintf("eventBus [%s] not registered", evt.Key()))
	}

	return bus.Publish(ctx, evt)
}
