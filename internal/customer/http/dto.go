package http

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/customer"
)

type Response struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL *string   `json:"image_url"`
}

func toResponse(customer *customer.Customer) Response {
	return Response{
		ID:       customer.ID,
		Name:     customer.Name,
		Email:    customer.Email,
		ImageURL: customer.ImageURL,
	}
}

func toResponses(customers []customer.Customer) []Response {
	rs := make([]Response, len(customers))
	for i, c := range customers {
		rs[i] = toResponse(&c)
	}
	return rs
}

type CreateRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	ImageURL *string `json:"image_url"`
}

func (req *CreateRequest) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, req); err != nil {
		return fmt.Errorf("validate CreateRequest: %w", err)
	}

	return nil
}

func (req *CreateRequest) toCustomer() *customer.Customer {
	return &customer.Customer{
		Name:     req.Name,
		Email:    req.Email,
		ImageURL: req.ImageURL,
	}
}
