package service

import (
	"database/sql"
	"fmt"
	"math"

	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/model"
	"go-nextjs-dashboard/request"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceService struct {
	DB *gorm.DB
}

func NewInvoiceService() *InvoiceService {
	return &InvoiceService{DB: config.DB}
}

func (s *InvoiceService) CreateInvoice(invoiceDto request.CreateInvoiceDTO) (*model.Invoice, error) {
	var count int64
	fmt.Println(invoiceDto.CustomerID)
	s.DB.Model(&model.Customer{}).Where("id = ?", invoiceDto.CustomerID).Count(&count)

	if count == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Customer not found.")
	}

	invoice := model.Invoice{
		CustomerID: &invoiceDto.CustomerID,
		Amount:     invoiceDto.Amount,
		Status:     invoiceDto.Status,
		Date:       invoiceDto.Date,
	}

	err := s.DB.Create(&invoice).Error
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (s *InvoiceService) GetInvoiceByID(invoiceId uuid.UUID) (*model.Invoice, error) {
	var invoice model.Invoice

	if err := s.DB.Select("id, customer_id, amount, status").Where("id = ?", invoiceId).First(&invoice).Error; err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (s *InvoiceService) GetInvoicePages(dto request.GetInvoicePagesDTO) (float64, error) {
	var count int64 = 0

	err := s.DB.Model(&model.Invoice{}).
		Joins("JOIN customers ON invoices.customer_id = customers.id").
		Where(`
			customers.name LIKE @search OR
			customers.email LIKE @search OR
			CAST(invoices.amount AS CHAR)  LIKE @search OR
			CAST(invoices.date AS CHAR) LIKE @search OR
			invoices.status LIKE @search
		`, sql.Named("search", "%"+dto.Search+"%")).
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	return math.Ceil(float64(count) / float64(dto.Limit)), nil
}

func (s *InvoiceService) GetFilteredInvoices(dto request.GetFilteredInvoicesDTO) ([]model.Invoice, error) {
	var invoices []model.Invoice

	offset := (dto.Page - 1) * dto.Limit

	err := s.DB.
		Preload("Customer").
		Joins("JOIN customers ON invoices.customer_id = customers.id").
		Where(`
			customers.name LIKE @search OR
			customers.email LIKE @search OR
			CAST(invoices.amount AS CHAR)  LIKE @search OR
			CAST(invoices.date AS CHAR) LIKE @search OR
			invoices.status LIKE @search
		`, sql.Named("search", "%"+dto.Search+"%")).
		Order("invoices.date DESC").
		Limit(dto.Limit).
		Offset(offset).
		Find(&invoices).
		Error

	if err != nil {
		return invoices, err
	}

	return invoices, nil
}

func (s *InvoiceService) GetLatestInvoices() ([]model.Invoice, error) {
	var invoices []model.Invoice
	limit := 5

	err := s.DB.Preload("Customer").Order("date DESC").Limit(limit).Find(&invoices).Error

	return invoices, err
}

func (s *InvoiceService) UpdateInvoice(dto request.UpdateInvoiceDTO) (*model.Invoice, error) {
	var invoice model.Invoice

	if err := s.DB.Where("id = ?", dto.InvoiceID).First(&invoice).Error; err != nil {
		return nil, err
	}

	invoice.CustomerID = &dto.CustomerID
	invoice.Amount = dto.Amount
	invoice.Status = dto.Status

	if err := s.DB.Save(&invoice).Error; err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (s *InvoiceService) DeleteInvoice(id uuid.UUID) error {
	var invoice model.Invoice
	if err := s.DB.Where("id = ?", id).First(&invoice).Delete(&invoice).Error; err != nil {
		return err
	}

	return nil
}
