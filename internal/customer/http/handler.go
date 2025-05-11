package http

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/customer/service"
	"go-nextjs-dashboard/internal/http"
)

var validator = http.Validator

func RegisterHTTP(r fiber.Router, svc *service.Service) {
	h := newHandler(svc)
	r.Get("/customers", h.List)
	r.Get("/customers/:id", h.Get)
	r.Post("/customers", h.Create, rateLimiter(5))
}

type handler struct {
	svc *service.Service
}

func newHandler(svc *service.Service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) List(c fiber.Ctx) error {
	cust, err := h.svc.List(c.Context())
	if err != nil {
		switch {
		case errors.Is(err, customer.ErrCustomerNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Customer not found.")
		default:
			return fmt.Errorf("retrieve customer: %w", err)
		}
	}

	return c.JSON(fiber.Map{
		"data": toResponses(cust),
	})
}

func (h *handler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid ID.")
	}

	cust, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, customer.ErrCustomerNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Customer not found.")
		default:
			return fmt.Errorf("retrieve customer: %w", err)
		}
	}

	return c.JSON(fiber.Map{
		"data": toResponse(cust),
	})
}

func (h *handler) Create(c fiber.Ctx) error {
	var req CreateRequest

	if err := c.Bind().Body(&req); err != nil {
		return fmt.Errorf("creaate customer bind request body: %w", err)
	}

	if err := req.Validate(c.Context()); err != nil {
		return fmt.Errorf("customer create request validation: %w", err)
	}

	reqCust := req.toCustomer()

	cust, err := h.svc.Create(c.Context(), *reqCust)

	if err != nil {
		switch {
		case errors.Is(err, customer.ErrEmailAlreadyTaken):
			return fiber.NewError(fiber.StatusConflict, "Email already taken.")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create customer: %s", err))
		}
	}

	return c.JSON(fiber.Map{"data": toResponse(cust)})
}
