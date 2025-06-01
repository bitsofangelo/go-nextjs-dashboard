package event

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Handler[T any] func(context.Context, T) error

type Bus[T any] struct {
	mu           sync.RWMutex
	handlers     []Handler[T]
	evt          T
	asyncHandler func(Handler[T]) Handler[T]
}

func NewBus[T any]() *Bus[T] {
	return &Bus[T]{
		asyncHandler: asyncHandler[T],
	}
}

func (b *Bus[T]) Subscribe(fn Handler[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, fn)
}

func (b *Bus[T]) SubscribeAsync(fn Handler[T]) {
	newFn := b.asyncHandler(fn)
	b.Subscribe(newFn)
}

func (b *Bus[T]) SetAsyncHandler(fn func(Handler[T]) Handler[T]) error {
	if fn == nil {
		return errors.New("nil async handler")
	}

	b.asyncHandler = fn

	return nil
}

func (b *Bus[T]) Publish(ctx context.Context, evt any) error {
	b.mu.RLock()
	handlers := append([]Handler[T]{}, b.handlers...)
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
		return errors.Join(errs...)
	}
	return nil
}

func (b *Bus[T]) Key() string {
	return reflect.TypeOf(b.evt).String()
}

func asyncHandler[T any](h Handler[T]) Handler[T] {
	return func(ctx context.Context, e T) error {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			if err := h(ctx, e); err != nil {
				fmt.Println(err)
			}
		}()
		return nil
	}
}
