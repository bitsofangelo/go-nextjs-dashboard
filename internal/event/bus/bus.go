package bus

import (
	"context"
	"fmt"
	"sync"

	"go-nextjs-dashboard/internal/event"
)

type eventBus[T event.Event] struct {
	mu       sync.RWMutex
	handlers []event.Handler[T]
	evt      T
}

func newBus[T event.Event]() *eventBus[T] {
	return &eventBus[T]{}
}

func (b *eventBus[T]) Subscribe(fn event.Handler[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, fn)
}

func (b *eventBus[T]) Publish(ctx context.Context, evt event.Event) error {
	b.mu.RLock()
	handlers := append([]event.Handler[T]{}, b.handlers...)
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
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
