package http

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/customer"
)

type response struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL *string   `json:"image_url"`
}

func toResponse(customer *customer.Customer) response {
	return response{
		ID:       customer.ID,
		Name:     customer.Name,
		Email:    customer.Email,
		ImageURL: customer.ImageURL,
	}
}

func toResponses(customers []customer.Customer) []response {
	rs := make([]response, len(customers))
	for i, c := range customers {
		rs[i] = toResponse(&c)
	}
	return rs
}

type createRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	ImageURL *string `json:"image_url"`
}

func (req *createRequest) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, req); err != nil {
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
