package event

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Publisher interface {
	Publish(ctx context.Context, event any) error
}

type KeyedPublisher interface {
	Publisher
	Key() string
}

type Broker struct {
	buses map[string]Publisher
}

func NewBroker() *Broker {
	return &Broker{
		buses: make(map[string]Publisher),
	}
}

func (r *Broker) RegisterBus(bus KeyedPublisher) {
	r.buses[bus.Key()] = bus
}

func (r *Broker) Publish(ctx context.Context, event any) error {
	t := reflect.TypeOf(event).String()
	bus, ok := r.buses[t]
	if !ok {
		return errors.New(fmt.Sprintf("eventBus [%s] not registered", t))
	}

	return bus.Publish(ctx, event)
}
