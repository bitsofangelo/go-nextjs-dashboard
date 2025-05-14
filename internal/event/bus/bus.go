package bus

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go-nextjs-dashboard/internal/event"
)

var ErrNoHandlers = errors.New("no handlers registered")

type Bus[T event.Event] interface {
	event.Publisher
	Subscribe(fn Handler[T])
	Name() string
}

type Handler[T event.Event] func(context.Context, T) error

type eventBus[T event.Event] struct {
	mu       sync.RWMutex
	handlers []Handler[T]
	evt      T
}

func newBus[T event.Event]() *eventBus[T] {
	return &eventBus[T]{}
}

func (b *eventBus[T]) Subscribe(fn Handler[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, fn)
}

func (b *eventBus[T]) Publish(ctx context.Context, evt event.Event) error {
	b.mu.RLock()
	handlers := append([]Handler[T]{}, b.handlers...)
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return ErrNoHandlers
	}

	e, ok := evt.(T)
	if !ok {
		return fmt.Errorf("invalid event type: %T", evt)
	}

	var errs []error
	for _, handle := range handlers {
		if err := handle(ctx, e); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// TODO: wrap all errs into one
		return errs[0]
	}
	return nil
}

func (b *eventBus[T]) Name() string {
	return b.evt.Key()
}
