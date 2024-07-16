package handler

import (
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/response"
	"go-nextjs-dashboard/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	DB              *gorm.DB
	customerService *service.CustomerService
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{DB: config.DB, customerService: service.NewCustomerService()}
}

func (h *CustomerHandler) GetCustomers(c *fiber.Ctx) error {
	customers, err := h.customerService.GetCustomers()
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.NewCustomersResponse(customers))
}

func (h *CustomerHandler) GetFilteredCustomers(c *fiber.Ctx) error {
	search := c.Query("search")

	filteredCustomers, err := h.customerService.GetFilteredCustomers(search)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.NewFilteredCustomersResponse(filteredCustomers))
}
