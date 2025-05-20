package event

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Mode int

const (
	ModeSync Mode = iota
	ModeAsync
)

type Event interface{}

type Publisher interface {
	Publish(ctx context.Context, evt Event) error
}

type Handler[T Event] func(context.Context, T) error

type Broker struct {
	buses map[string]Publisher
}

func NewBroker(buses map[string]Publisher) *Broker {
	return &Broker{
		buses: buses,
	}
}

func (r *Broker) Publish(ctx context.Context, evt Event) error {
	t := reflect.TypeOf(evt).String()
	bus, ok := r.buses[t]
	if !ok {
		return errors.New(fmt.Sprintf("eventBus [%s] not registered", t))
	}

	return bus.Publish(ctx, evt)
}
