package http

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
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
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid ID.")
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
		return fmt.Errorf("creaate invoice bind request body: %w", err)
	}

	if err := Validator.StructCtx(c.Context(), &req); err != nil {
		return fmt.Errorf("create invoice validation: %w", err)
	}

	var reqInv invoice.Invoice
	reqInv.CustomerID, _ = uuid.Parse(req.CustomerID)
	reqInv.Amount = req.Amount
	reqInv.Status = req.Status
	reqInv.Date = req.Date

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

type invoiceResponse struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Amount     float32   `json:"amount"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
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
	Amount     float32   `json:"amount" validate:"required"`
	Status     string    `json:"status" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
}
