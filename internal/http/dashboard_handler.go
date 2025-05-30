package http

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/gelozr/go-dash/internal/dashboard"
	"github.com/gelozr/go-dash/internal/http/response"
	"github.com/gelozr/go-dash/internal/logger"
)

type DashboardHandler struct {
	svc    *dashboard.Service
	logger logger.Logger
}

func NewDashboardHandler(svc *dashboard.Service, log logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		svc:    svc,
		logger: log.With("component", "http.dashboard"),
	}
}

func (h *DashboardHandler) GetOverview(c fiber.Ctx) error {
	time.Sleep(5 * time.Second)
	o, err := h.svc.GetOverview(c.Context())
	if err != nil {
		return fmt.Errorf("get overview: %w", err)
	}

	return c.JSON(
		response.New(response.ToOverview(o)),
	)
}

func (h *DashboardHandler) GetMonthlyRevenues(c fiber.Ctx) error {
	revs, err := h.svc.GetMonthlyRevenues(c.Context())
	if err != nil {
		return fmt.Errorf("get monthly revenues: %w", err)
	}

	return c.JSON(
		response.New(response.ToMonthlyRevenueList(revs)),
	)
}
