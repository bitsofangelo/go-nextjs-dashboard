package http

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/logger"
)

type customerHandler struct {
	svc    *customer.Service
	logger logger.Logger
}

func newCustomerHandler(svc *customer.Service, log logger.Logger) *customerHandler {
	return &customerHandler{
		svc:    svc,
		logger: log.With("component", "http.customer"),
	}
}

func (h *customerHandler) List(c fiber.Ctx) error {
	customers, err := h.svc.List(c.Context())
	if err != nil {
		switch {
		case errors.Is(err, customer.ErrCustomerNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Customer not found.")
		default:
			return fmt.Errorf("retrieve customer: %w", err)
		}
	}

	rs := make([]customerResponse, len(customers))
	for i, cust := range customers {
		rs[i] = toCustomerResponse(&cust)
	}

	return c.JSON(Response{
		Data: rs,
	})
}

func (h *customerHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid ID.")
	}

	cust, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, customer.ErrCustomerNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Customer not found.")
		default:
			return fmt.Errorf("get customer by id: %w", err)
		}
	}

	return c.JSON(Response{
		Data: toCustomerResponse(cust),
	})
}

func (h *customerHandler) Create(c fiber.Ctx) error {
	var req createRequest

	if err := c.Bind().Body(&req); err != nil {
		return fmt.Errorf("creaate customer bind request body: %w", err)
	}

	if err := req.Validate(c.Context()); err != nil {
		return fmt.Errorf("customer create request validation: %w", err)
	}

	reqCust := req.toCustomer()

	cust, err := h.svc.Create(c.Context(), *reqCust)

	if err != nil {
		switch {
		case errors.Is(err, customer.ErrEmailAlreadyTaken):
			h.logger.ErrorContext(c.Context(), "email already taken")
			return fiber.NewError(fiber.StatusConflict, "Email already taken.")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("create customer: %s", err))
		}
	}

	return c.JSON(Response{Data: toCustomerResponse(cust)})
}

func (h *customerHandler) SearchWithInvoiceTotals(c fiber.Ctx) error {
	search := c.Query("search")
	result, err := h.svc.SearchWithInvoiceTotals(c.Context(), search)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("retrieve customer: %s", err))
	}

	return c.JSON(Response{
		Data: toResponseWithInvoicesTotals(result),
	})
}

type customerResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL *string   `json:"image_url"`
}

func toCustomerResponse(customer *customer.Customer) customerResponse {
	return customerResponse{
		ID:       customer.ID,
		Name:     customer.Name,
		Email:    customer.Email,
		ImageURL: customer.ImageURL,
	}
}

type createRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	ImageURL *string `json:"image_url"`
}

func (req *createRequest) Validate(ctx context.Context) error {
	if err := Validator.StructCtx(ctx, req); err != nil {
		return fmt.Errorf("validate createRequest: %w", err)
	}

	return nil
}

func (req *createRequest) toCustomer() *customer.Customer {
	return &customer.Customer{
		Name:     req.Name,
		Email:    req.Email,
		ImageURL: req.ImageURL,
	}
}

type withInvoiceTotalsResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ImageURL      *string   `json:"image_url"`
	TotalInvoices int64     `json:"total_invoices"`
	TotalPending  float64   `json:"total_pending"`
	TotalPaid     float64   `json:"total_paid"`
}

func toResponseWithInvoiceTotals(it customer.WithInvoiceTotals) withInvoiceTotalsResponse {
	return withInvoiceTotalsResponse{
		ID:            it.ID,
		Name:          it.Name,
		Email:         it.Email,
		ImageURL:      it.ImageURL,
		TotalInvoices: it.TotalInvoices,
		TotalPending:  it.TotalPending,
		TotalPaid:     it.TotalPaid,
	}
}

func toResponseWithInvoicesTotals(its []customer.WithInvoiceTotals) []withInvoiceTotalsResponse {
	rs := make([]withInvoiceTotalsResponse, len(its))
	for i, it := range its {
		rs[i] = toResponseWithInvoiceTotals(it)
	}
	return rs
}
