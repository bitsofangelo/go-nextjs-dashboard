package handler

import (
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/model"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OverviewHandler struct {
	DB *gorm.DB
}

func NewOverviewHandler() *OverviewHandler {
	return &OverviewHandler{DB: config.DB}
}

func (h *OverviewHandler) GetOverviewData(c *fiber.Ctx) error {
	// time.Sleep(time.Second)
	var invoiceCount int64
	var customerCount int64
	var invoiceStatus struct {
		Paid    float64
		Pending float64
	}
	var wg sync.WaitGroup
	var err1, err2, err3 error

	wg.Add(3)

	go func() {
		defer wg.Done()
		// time.Sleep(time.Millisecond * 300)
		err1 = h.DB.Model(&model.Invoice{}).Count(&invoiceCount).Error
	}()

	go func() {
		defer wg.Done()
		// time.Sleep(time.Millisecond * 300)
		err2 = h.DB.Model(&model.Customer{}).Count(&customerCount).Error
	}()

	go func() {
		defer wg.Done()
		// time.Sleep(time.Millisecond * 300)
		err3 = h.DB.Model(&model.Invoice{}).Select(`
			SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END) AS "paid",
			SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END) AS "pending"
		`).Scan(&invoiceStatus).Error
	}()

	wg.Wait()

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}

	return c.Status(200).JSON(fiber.Map{"data": fiber.Map{
		"invoice_count":  invoiceCount,
		"customer_count": customerCount,
		"invoice_status": fiber.Map{
			"pending": invoiceStatus.Pending,
			"paid":    invoiceStatus.Paid,
		},
	}})
}
