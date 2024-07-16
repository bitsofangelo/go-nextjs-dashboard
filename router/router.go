package router

import (
	"go-nextjs-dashboard/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(r *fiber.App) {
	// registerUserRoutes(r)
	// registerCustomerRoutes(r)
	// registerInvoiceRoutes(r)
	// registerRevenueRoutes(r)

	overviewHander := handler.NewOverviewHandler()
	revenueHandler := handler.NewRevenueHandler()
	invoiceHandler := handler.NewInvoiceHandler()
	customerHandler := handler.NewCustomerHandler()
	userHandler := handler.NewUserHandler()

	api := r.Group("/api")
	{
		api.Get("/overview", overviewHander.GetOverviewData)
		api.Get("/revenues", revenueHandler.GetRevenues)

		api.Get("/invoices/latest", invoiceHandler.GetLatestInvoices)
		api.Get("/invoices/filtered", invoiceHandler.GetFilteredInvoices)
		api.Get("/invoices/total-pages", invoiceHandler.GetInvoicePages)
		api.Get("/invoices/:id", invoiceHandler.GetInvoiceByID)
		api.Post("/invoices", invoiceHandler.CreateInvoice)
		api.Put("/invoices/:id", invoiceHandler.UpdateInvoice)
		api.Delete("/invoices/:id", invoiceHandler.DeleteInvoice)

		api.Get("/customers", customerHandler.GetCustomers)
		api.Get("/customers/filtered", customerHandler.GetFilteredCustomers)

		api.Get("/users/email/:email", userHandler.GetUserByEmail)
	}
}
