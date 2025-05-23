package request

import (
	"context"
	"fmt"

	"go-nextjs-dashboard/internal/customer"
)

type CreateCustomer struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	ImageURL *string `json:"image_url"`
}

func (req *CreateCustomer) Validate(ctx context.Context) error {
	if err := validator.StructCtx(ctx, req); err != nil {
		return fmt.Errorf("validate create customer: %w", err)
	}

	return nil
}

func (req *CreateCustomer) ToCustomer() customer.Customer {
	return customer.Customer{
		Name:     req.Name,
		Email:    req.Email,
		ImageURL: req.ImageURL,
	}
}
