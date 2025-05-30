package request

import (
	"github.com/gelozr/go-dash/internal/customer"
)

type CreateCustomer struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	ImageURL *string `json:"image_url"`
}

func (req *CreateCustomer) ToCustomer() customer.Customer {
	return customer.Customer{
		Name:     req.Name,
		Email:    req.Email,
		ImageURL: req.ImageURL,
	}
}
