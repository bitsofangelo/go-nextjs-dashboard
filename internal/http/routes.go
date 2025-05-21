package http

import (
	"github.com/gofiber/fiber/v3"

	"go-nextjs-dashboard/internal/app"
	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/dashboard"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
	"go-nextjs-dashboard/internal/user"
)

func RegisterRoutes(
	r fiber.Router,

	// domain services
	custSvc *customer.Service,
	userSvc *user.Service,
	dashSvc *dashboard.Service,
	invSvc *invoice.Service,

	// use cases
	createInvoice *app.CreateInvoice,

	logger logger.Logger,
) {
	dashHandler := newDashboardHandler(dashSvc, logger)
	usrHandler := newUserHandler(userSvc, logger)
	custHandler := newCustomerHandler(custSvc, logger)
	invHandler := newInvoiceHandler(invSvc, createInvoice, logger)

	// dashboard routes
	r.Get("/overview", dashHandler.GetOverview, loggerKeyMiddleware("http.dashboard"))

	// user routes
	r.Get("/users/email/:email", usrHandler.GetByEmail, loggerKeyMiddleware("http.user"))

	// customer routes
	cg := r.Group("/customers", loggerKeyMiddleware("http.customer"))
	{
		cg.Get("/", custHandler.List)
		cg.Get("/filtered", custHandler.SearchWithInvoiceTotals)
		cg.Get("/:id", custHandler.Get)
		cg.Post("/", custHandler.Create, rateLimiter(30))
	}

	// invoice routes
	ig := r.Group("/invoices", loggerKeyMiddleware("http.invoice"))
	{
		ig.Get("/latest", invHandler.GetLatest)
		ig.Get("/filtered", invHandler.Search)
		// ig.Get("/invoices/total-pages", invoiceHandler.GetInvoicePages)

		ig.Get("/:id", invHandler.Get)
		ig.Post("/", invHandler.Create, rateLimiter(30))
		ig.Patch("/:id", invHandler.Update, rateLimiter(30))
		ig.Delete("/:id", invHandler.Delete, rateLimiter(30))
	}
}
