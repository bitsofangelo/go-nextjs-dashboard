package event

import "github.com/google/uuid"

type Created struct {
	ID uuid.UUID
}

func (Created) Key() string {
	return "customer.created"
}
