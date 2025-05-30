package bus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/gelozr/go-dash/internal/logger"
)

type Handler[T any] func(context.Context, T) error

type eventBus[T any] struct {
	mu       sync.RWMutex
	handlers []Handler[T]
	evt      T
}

func newBus[T any]() *eventBus[T] {
	return &eventBus[T]{}
}

func (b *eventBus[T]) Subscribe(fn Handler[T]) {
	// b.mu.Lock()
	// defer b.mu.Unlock()
	b.handlers = append(b.handlers, fn)
}

func (b *eventBus[T]) Publish(ctx context.Context, evt any) error {
	// no need locks since b.handlers is read-only after init
	// b.mu.RLock()
	// handlers := append([]Handler[T]{}, b.handlers...)
	// b.mu.RUnlock()

	if len(b.handlers) == 0 {
		return nil
	}

	e, ok := evt.(T)
	if !ok {
		return fmt.Errorf("invalid event type: %T", evt)
	}

	var errs []error
	for _, handle := range b.handlers {
		if err := handle(ctx, e); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (b *eventBus[T]) Key() string {
	return reflect.TypeOf(b.evt).String()
}

func safeGo(ctx context.Context, log logger.Logger, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ErrorContext(ctx, fmt.Sprintf("panic: %s", r))
			}
		}()
		fn()
	}()
}
