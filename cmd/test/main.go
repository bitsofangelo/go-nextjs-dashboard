package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan error)

	go func() {
		time.Sleep(5 * time.Second)
		done <- errors.New("test")
	}()

	select {
	case <-ctx.Done():
		fmt.Println("ctx done", ctx.Err())
	case err := <-done:
		fmt.Println("err", err)
	}
}
