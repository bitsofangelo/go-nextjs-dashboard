package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/model"
)

type OverviewHandler struct {
	DB *gorm.DB
}

func NewOverviewHandler() *OverviewHandler {
	return &OverviewHandler{DB: config.DB}
}

func (h *OverviewHandler) GetOverviewData(c *fiber.Ctx) error {

	ctx := c.UserContext()
	g, ctx := errgroup.WithContext(ctx)

	// time.Sleep(time.Second)

	var (
		invoiceCount  int64
		customerCount int64
		invoiceStatus struct{ Paid, Pending float64 }
	)

	g.Go(func() error {
		// time.Sleep(time.Millisecond * 300)
		return h.DB.WithContext(ctx).Model(&model.Invoice{}).Count(&invoiceCount).Error
	})

	g.Go(func() error {
		// time.Sleep(time.Millisecond * 300)
		return h.DB.WithContext(ctx).Model(&model.Customer{}).Count(&customerCount).Error
	})

	g.Go(func() error {
		// time.Sleep(time.Millisecond * 300)
		return h.DB.WithContext(ctx).Model(&model.Invoice{}).Select(`
			SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END) AS "paid",
			SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END) AS "pending"
		`).Scan(&invoiceStatus).Error
	})

	if err := g.Wait(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
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
