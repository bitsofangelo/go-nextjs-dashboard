package event

import (
	"context"
	"fmt"
)

func SendWelcomeMessage(_ context.Context, e Created) error {
	go func() {
		fmt.Println("Welcome message received", e.ID)
	}()
	return nil
}
