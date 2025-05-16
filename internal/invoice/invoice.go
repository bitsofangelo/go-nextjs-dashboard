package invoice

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Amount     float32
	Status     string
	Date       time.Time
}
