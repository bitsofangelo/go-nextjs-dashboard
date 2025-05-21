package request

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/http/validation"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/optional"
)

var validator = validation.Validator

type CreateInvoice struct {
	CustomerID string    `json:"customer_id" validate:"required,uuid4"`
	Amount     float64   `json:"amount" validate:"required"`
	Status     string    `json:"status" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
}

func (req *CreateInvoice) Validate(ctx context.Context) error {
	return validator.StructCtx(ctx, req)
}

func (req *CreateInvoice) ToDTO() (invoice.Invoice, error) {
	custID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return invoice.Invoice{}, fiber.NewError(fiber.StatusUnprocessableEntity, "invalid customer id.")
	}

	return invoice.Invoice{
		CustomerID: &custID,
		Amount:     req.Amount,
		Status:     req.Status,
		Date:       req.Date,
	}, nil
}

type UpdateInvoice struct {
	CustomerID optional.Optional[string] `json:"customer_id" validate:"omitempty,uuid4"`
	Amount     float64                   `json:"amount" validate:"min=0,max=100"`
	Status     string                    `json:"status" validate:"required"`
	Date       string                    `json:"date" validate:"required,rfc3339"`
}

func (req *UpdateInvoice) Validate(ctx context.Context) error {
	return validator.StructCtx(ctx, req)
}

func (req *UpdateInvoice) ToDTO() (*invoice.UpdateInput, error) {
	customerID, err := optional.StringToUUID(req.CustomerID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "invalid customer id.")
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "invalid date.")
	}

	return &invoice.UpdateInput{
		CustomerID: customerID,
		Amount:     req.Amount,
		Status:     req.Status,
		Date:       date,
	}, nil
}
