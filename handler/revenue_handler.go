package handler

import (
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/model"
	"go-nextjs-dashboard/response"
	"sort"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RevenueHandler struct {
	DB *gorm.DB
}

func NewRevenueHandler() *RevenueHandler {
	return &RevenueHandler{DB: config.DB}
}

func (h *RevenueHandler) GetRevenues(c *fiber.Ctx) error {
	// time.Sleep(time.Millisecond * 200)
	var revenues model.Revenues

	h.DB.Find(&revenues)
	sort.Sort(revenues)

	return c.Status(fiber.StatusOK).JSON(response.NewRevenueResponse(revenues))
}
