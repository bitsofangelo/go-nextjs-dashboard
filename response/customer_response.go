package response

import (
	"go-nextjs-dashboard/model"
	"go-nextjs-dashboard/service"
	"time"

	"github.com/google/uuid"
)

type CustomerInvoiceRes struct {
	ID     uuid.UUID `json:"id"`
	Amount float32   `json:"amount"`
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
}

type CustomerRes struct {
	ID         uuid.UUID            `json:"id"`
	Name       string               `json:"name"`
	Email      string               `json:"email"`
	ImageURL   string               `json:"image_url"`
	InvoiceRes []CustomerInvoiceRes `json:"invoices"`
}

func NewCustomerResponse(customer model.Customer) CustomerRes {
	response := CustomerRes{
		ID:       customer.ID,
		Name:     customer.Name,
		Email:    customer.Email,
		ImageURL: customer.ImageURL,
	}

	response.InvoiceRes = make([]CustomerInvoiceRes, len(customer.Invoices))

	if len(customer.Invoices) > 0 {
		for i, invoice := range customer.Invoices {
			response.InvoiceRes[i] = CustomerInvoiceRes{
				ID:     invoice.ID,
				Amount: invoice.Amount,
				Status: invoice.Status,
				Date:   invoice.Date,
			}
		}
	}

	// rels := strings.Split(include, ",")

	// for _, rel := range rels {
	// 	if rel == "invoices" {
	// 		invoiceResponse := NewInvoicesResponse(customer.Invoices, "")
	// 		response.Invoices = &invoiceResponse
	// 	}
	// }

	return response
}

func NewCustomersResponse(customers []model.Customer) map[string]any {
	type respStruct struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	response := make([]respStruct, len(customers))

	for i, customer := range customers {
		response[i] = respStruct{
			ID:   customer.ID,
			Name: customer.Name,
		}
	}

	return map[string]any{"data": response}
}

func NewFilteredCustomersResponse(filteredCustomers []service.FilteredCustomerDTO) map[string]any {
	type respStruct struct {
		ID            uuid.UUID `json:"id"`
		Name          string    `json:"name"`
		Email         string    `json:"email"`
		ImageURL      string    `json:"image_url"`
		TotalInvoices uint32    `json:"total_invoices"`
		TotalPending  float64   `json:"total_pending"`
		TotalPaid     float64   `json:"total_paid"`
	}

	response := make([]respStruct, len(filteredCustomers))

	for i, customer := range filteredCustomers {
		response[i] = respStruct{
			ID:            customer.ID,
			Name:          customer.Name,
			Email:         customer.Email,
			ImageURL:      customer.ImageURL,
			TotalInvoices: customer.TotalInvoices,
			TotalPending:  customer.TotalPending,
			TotalPaid:     customer.TotalPaid,
		}
	}

	return map[string]any{"data": response}
}
