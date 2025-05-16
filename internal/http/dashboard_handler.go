package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"go-nextjs-dashboard/internal/dashboard"
	"go-nextjs-dashboard/internal/logger"
)

type dashboardHandler struct {
	svc    *dashboard.Service
	logger logger.Logger
}

func newDashboardHandler(svc *dashboard.Service, log logger.Logger) *dashboardHandler {
	return &dashboardHandler{
		svc:    svc,
		logger: log.With("component", "http.dashboard"),
	}
}

func (h *dashboardHandler) GetOverview(c fiber.Ctx) error {
	o, err := h.svc.GetOverview(c.Context())
	if err != nil {
		return fmt.Errorf("get overview: %w", err)
	}

	return c.JSON(Response{
		Data: toOverviewResponse(o),
	})
}

type OverviewResponse struct {
	InvoiceCount  int64 `json:"invoice_count"`
	CustomerCount int64 `json:"customer_count"`
	InvoiceStatus struct {
		Paid    float64 `json:"paid"`
		Pending float64 `json:"pending"`
	} `json:"invoice_status"`
}

func toOverviewResponse(o *dashboard.Overview) OverviewResponse {
	var r OverviewResponse
	r.CustomerCount = o.CustomerCount
	r.InvoiceCount = o.InvoiceCount
	r.InvoiceStatus.Paid = o.InvoiceStatus.Paid
	r.InvoiceStatus.Pending = o.InvoiceStatus.Pending

	return r
}
