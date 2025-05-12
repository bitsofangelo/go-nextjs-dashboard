package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"go-nextjs-dashboard/internal/dashboard/service"
	"go-nextjs-dashboard/internal/http"
	"go-nextjs-dashboard/internal/logger"
)

var validator = http.Validator

func RegisterHTTP(r fiber.Router, svc *service.Service, log logger.Logger) {
	h := newHandler(svc, log)
	r.Get("/overview", h.GetOverview)
}

type handler struct {
	svc    *service.Service
	logger logger.Logger
}

func newHandler(svc *service.Service, log logger.Logger) *handler {
	return &handler{
		svc:    svc,
		logger: log,
	}
}

func (h *handler) GetOverview(c fiber.Ctx) error {
	o, err := h.svc.GetOverview(c.Context())
	if err != nil {
		return fmt.Errorf("get overview: %w", err)
	}

	return c.JSON(http.Response{
		Data: newOverviewResponse(o),
	})
}
