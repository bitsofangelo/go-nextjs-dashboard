package response

import (
	"github.com/google/uuid"

	"github.com/gelozr/go-dash/internal/customer"
)

type Customer struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL *string   `json:"image_url"`
}

func ToCustomer(customer customer.Customer) Customer {
	return Customer{
		ID:       customer.ID,
		Name:     customer.Name,
		Email:    customer.Email,
		ImageURL: customer.ImageURL,
	}
}

func ToCustomers(data []customer.Customer) []Customer {
	return ToList(data, ToCustomer)
}

type CustomerWithInvoiceInfo struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ImageURL      *string   `json:"image_url"`
	TotalInvoices int64     `json:"total_invoices"`
	TotalPending  float64   `json:"total_pending"`
	TotalPaid     float64   `json:"total_paid"`
}

func ToCustomerWithInvoiceInfo(c customer.WithInvoiceInfo) CustomerWithInvoiceInfo {
	return CustomerWithInvoiceInfo{
		ID:            c.ID,
		Name:          c.Name,
		Email:         c.Email,
		ImageURL:      c.ImageURL,
		TotalInvoices: c.TotalInvoices,
		TotalPending:  c.TotalPending,
		TotalPaid:     c.TotalPaid,
	}
}

func ToCustomerWithInvoiceInfoList(data []customer.WithInvoiceInfo) []CustomerWithInvoiceInfo {
	return ToList(data, ToCustomerWithInvoiceInfo)
}
