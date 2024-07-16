package handler

import (
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/request"
	"go-nextjs-dashboard/response"
	"go-nextjs-dashboard/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceHandler struct {
	invoiceService *service.InvoiceService
	DB             *gorm.DB
}

func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{DB: config.DB, invoiceService: service.NewInvoiceService()}
}

func (h *InvoiceHandler) GetLatestInvoices(c *fiber.Ctx) error {
	invoices, err := h.invoiceService.GetLatestInvoices()
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.NewLatestInvoicesResponse(invoices))
}

func (h *InvoiceHandler) GetFilteredInvoices(c *fiber.Ctx) error {
	// time.Sleep(time.Second)
	dto, err := request.NewGetFilteredInvoicesRequest(c)
	if err != nil {
		return err
	}

	invoices, err := h.invoiceService.GetFilteredInvoices(*dto)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.NewFiltelredInvoiceResponse(invoices))
}

func (h *InvoiceHandler) GetInvoicePages(c *fiber.Ctx) error {
	dto, err := request.NewGetInvoicePagesRequest(c)
	if err != nil {
		return err
	}

	totalPages, err := h.invoiceService.GetInvoicePages(*dto)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": fiber.Map{
		"pages": totalPages,
	}})
}

func (h *InvoiceHandler) GetInvoiceByID(c *fiber.Ctx) error {
	invoiceId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}

	invoice, err := h.invoiceService.GetInvoiceByID(invoiceId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"message": "Invoice not found"})
		}

		return err
	}

	result := struct {
		ID         uuid.UUID `json:"id"`
		CustomerID uuid.UUID `json:"customer_id"`
		Amount     float32   `json:"amount"`
		Status     string    `json:"status"`
	}{
		ID:         invoice.ID,
		CustomerID: *invoice.CustomerID,
		Amount:     invoice.Amount,
		Status:     invoice.Status,
	}

	return c.Status(200).JSON(fiber.Map{"data": result})
}

func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	time.Sleep(time.Millisecond * 1000)
	invoiceData, err := request.NewCreateInvoiceRequest(c)
	if err != nil {
		return err
	}

	invoice, err := h.invoiceService.CreateInvoice(*invoiceData)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewCreateInvoiceResponse(*invoice))
}

func (h *InvoiceHandler) UpdateInvoice(c *fiber.Ctx) error {
	dto, err := request.NewUpdateInvoiceRequest(c)
	if err != nil {
		return err
	}

	invoice, err := h.invoiceService.UpdateInvoice(*dto)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"message": "Invoice not found"})
		}

		return err
	}

	return c.Status(200).JSON(response.NewUpdateInvoiceResponse(*invoice))
}

func (h *InvoiceHandler) DeleteInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := config.Validate.Var(id, "required,uuid4"); err != nil {
		return fiber.NewError(400, "Invalid Invoice ID in path parameters")
	}

	invoiceID, _ := uuid.Parse(id)

	if err := h.invoiceService.DeleteInvoice(invoiceID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"message": "Invoice not found"})
		}

		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": true})
}
