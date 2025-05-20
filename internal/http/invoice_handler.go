package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/optional"
)

type invoiceHandler struct {
	invSvc        *invoice.Service
	createInvoice *app.CreateInvoice
	logger        logger.Logger
}

func newInvoiceHandler(invSvc *invoice.Service, createInvoice *app.CreateInvoice, logger logger.Logger) *invoiceHandler {
	return &invoiceHandler{
		invSvc:        invSvc,
		createInvoice: createInvoice,
		logger:        logger.With("component", "http.invoice"),
	}
}

func (h *invoiceHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid id.")
	}

	inv, err := h.invSvc.GetByID(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, invoice.ErrInvoiceNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Invoice not found.")
		default:
			return fmt.Errorf("get invoice by id: %w", err)
		}
	}

	return c.JSON(Response{
		Data: toInvoiceResponse(inv),
	})
}

func (h *invoiceHandler) Create(c fiber.Ctx) error {
	var req createInvoiceRequest

	if err := c.Bind().Body(&req); err != nil {
		return fmt.Errorf("create invoice bind request body: %w", err)
	}

	if err := validator.StructCtx(c.Context(), &req); err != nil {
		return fmt.Errorf("create invoice validation: %w", err)
	}

	custID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid customer id.")
	}

	reqInv := invoice.Invoice{
		CustomerID: &custID,
		Amount:     req.Amount,
		Status:     req.Status,
		Date:       req.Date,
	}

	inv, err := h.createInvoice.Execute(c.Context(), reqInv)
	if err != nil {
		switch {
		case errors.Is(err, customer.ErrCustomerNotFound):
			return fiber.NewError(fiber.StatusNotFound, "customer not found.")
		default:
			return fmt.Errorf("create invoice: %w", err)
		}
	}

	return c.JSON(Response{
		Data: toInvoiceResponse(inv),
	})
}

func (h *invoiceHandler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid id.")
	}

	var req updateInvoiceRequest
	if err = c.Bind().Body(&req); err != nil {
		return fmt.Errorf("update invoice bind request body: %w", err)
	}

	if err = req.Validate(c.Context()); err != nil {
		return err
	}

	updateReq, err := req.ToInvoiceUpdateReq()
	if err != nil {
		return err
	}

	inv, err := h.invSvc.Update(c.Context(), id, updateReq)
	if err != nil {
		return fmt.Errorf("update invoice by id: %w", err)
	}

	return c.JSON(Response{
		Data: toInvoiceResponse(inv),
	})
}

type invoiceResponse struct {
	ID         uuid.UUID  `json:"id"`
	CustomerID *uuid.UUID `json:"customer_id"`
	Amount     float64    `json:"amount"`
	Status     string     `json:"status"`
	Date       time.Time  `json:"date"`
}

func toInvoiceResponse(inv *invoice.Invoice) invoiceResponse {
	return invoiceResponse{
		ID:         inv.ID,
		CustomerID: inv.CustomerID,
		Amount:     inv.Amount,
		Status:     inv.Status,
		Date:       inv.Date,
	}
}

type createInvoiceRequest struct {
	CustomerID string    `json:"customer_id" validate:"required,uuid4"`
	Amount     float64   `json:"amount" validate:"required"`
	Status     string    `json:"status" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
}

type updateInvoiceRequest struct {
	CustomerID optional.Optional[string] `json:"customer_id" validate:"omitempty,uuid4"`
	Amount     float64                   `json:"amount" validate:"min=0,max=100"`
	Status     string                    `json:"status" validate:"required"`
	Date       string                    `json:"date" validate:"required,rfc3339"`
}

func (req *updateInvoiceRequest) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, req); err != nil {
		return fmt.Errorf("update invoice validation: %w", err)
	}
	return nil
}

func (req *updateInvoiceRequest) ToInvoiceUpdateReq() (*invoice.UpdateRequest, error) {
	customerID, err := optional.StringToUUID(req.CustomerID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "invalid customer id.")
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "invalid date.")
	}

	return &invoice.UpdateRequest{
		CustomerID: customerID,
		Amount:     req.Amount,
		Status:     req.Status,
		Date:       date,
	}, nil
}
