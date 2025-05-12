package response

import (
	"time"

	"go-nextjs-dashboard/model"

	"github.com/google/uuid"
)

type InvoiceCustomerResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL string    `json:"image_url"`
}

type InvoiceResponse struct {
	ID         uuid.UUID                `json:"id"`
	CustomerID *uuid.UUID               `json:"customer_id"`
	Amount     float32                  `json:"amount"`
	Status     string                   `json:"status"`
	Date       time.Time                `json:"date"`
	Customer   *InvoiceCustomerResponse `json:"customer"`
}

func NewInvoiceResponse(invoice model.Invoice) map[string]any {
	response := InvoiceResponse{
		ID:         invoice.ID,
		CustomerID: invoice.CustomerID,
		Amount:     invoice.Amount,
		Status:     invoice.Status,
		Date:       invoice.Date,
	}

	return map[string]any{"data": response}
}

func NewUpdateInvoiceResponse(invoice model.Invoice) map[string]any {
	response := struct {
		ID         uuid.UUID  `json:"id"`
		CustomerID *uuid.UUID `json:"customer_id"`
		Amount     float32    `json:"amount"`
		Status     string     `json:"status"`
		Date       time.Time  `json:"date"`
	}{
		ID:         invoice.ID,
		CustomerID: invoice.CustomerID,
		Amount:     invoice.Amount,
		Status:     invoice.Status,
		Date:       invoice.Date,
	}

	return map[string]any{"data": response}
}

type LatestInvoiceResponse struct {
	ID       uuid.UUID `json:"id"`
	Amount   float32   `json:"amount"`
	Name     *string   `json:"name"`
	Email    *string   `json:"email"`
	ImageURL *string   `json:"image_url"`
}

func NewLatestInvoicesResponse(latestInvoices []model.Invoice) map[string]any {
	latestInvoicesRes := make([]LatestInvoiceResponse, len(latestInvoices))

	for i, latestInvoice := range latestInvoices {
		latestInvoicesRes[i] = LatestInvoiceResponse{
			ID:     latestInvoice.ID,
			Amount: latestInvoice.Amount,
		}

		if latestInvoice.Customer != nil {
			latestInvoicesRes[i].Name = &latestInvoice.Customer.Name
			latestInvoicesRes[i].Email = &latestInvoice.Customer.Email
			latestInvoicesRes[i].ImageURL = &latestInvoice.Customer.ImageURL
		}
	}

	return map[string]any{"data": latestInvoicesRes}
}

type CreateInvoiceResponse struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Amount     float32   `json:"amount"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
}

func NewCreateInvoiceResponse(invoice model.Invoice) map[string]any {
	invoiceRes := CreateInvoiceResponse{
		ID:         invoice.ID,
		CustomerID: *invoice.CustomerID,
		Amount:     invoice.Amount,
		Status:     invoice.Status,
		Date:       invoice.Date,
	}

	return map[string]any{"data": invoiceRes}
}

type FilteredInvoiceResponse struct {
	ID       uuid.UUID `json:"id"`
	Amount   float32   `json:"amount"`
	Date     time.Time `json:"date"`
	Status   string    `json:"status"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL string    `json:"image_url"`
}

func NewFiltelredInvoiceResponse(invoices []model.Invoice) map[string]any {
	respData := make([]FilteredInvoiceResponse, len(invoices))

	for i, invoice := range invoices {
		respData[i] = FilteredInvoiceResponse{
			ID:       invoice.ID,
			Amount:   invoice.Amount,
			Date:     invoice.Date,
			Status:   invoice.Status,
			Name:     invoice.Customer.Name,
			Email:    invoice.Customer.Email,
			ImageURL: invoice.Customer.ImageURL,
		}
	}

	return map[string]any{"data": respData}
}
