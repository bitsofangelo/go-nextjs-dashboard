package http

import "go-nextjs-dashboard/internal/auth"

type RouteInitializer struct{}

func SetupFiberRoutes(
	s *FiberServer,
	authSvc *auth.Token,
	authH *AuthHandler,
	dashH *DashboardHandler,
	userH *UserHandler,
	custH *CustomerHandler,
	invH *InvoiceHandler,
) RouteInitializer {

	r := s.app.Group("/api")

	// auth routes
	ag := r.Group("/auth", loggerKeyMiddleware("http.auth"))
	{
		ag.Post("/login", authH.Login)
	}

	// dashboard routes
	dg := r.Group("/dash", loggerKeyMiddleware("http.dashboard"), AuthMiddleware(authSvc))
	{
		dg.Get("/overview", dashH.GetOverview)
		dg.Get("/revenues", dashH.GetMonthlyRevenues)
	}

	// user routes
	r.Get("/users/email/:email", userH.GetByEmail, loggerKeyMiddleware("http.user"))

	// customer routes
	cg := r.Group("/customers", loggerKeyMiddleware("http.customer"))
	{
		cg.Get("/", custH.List)
		cg.Get("/filtered", custH.SearchWithInvoiceInfo)
		cg.Get("/:id", custH.Get)
		cg.Post("/", custH.Create, rateLimiter(30))
	}

	// invoice routes
	ig := r.Group("/invoices", loggerKeyMiddleware("http.invoice"))
	{
		ig.Get("/latest", invH.GetLatest)
		ig.Get("/filtered", invH.Search)

		ig.Get("/:id", invH.Get)
		ig.Post("/", invH.Create, rateLimiter(30))
		ig.Patch("/:id", invH.Update, rateLimiter(30))
		ig.Delete("/:id", invH.Delete, rateLimiter(30))
	}

	return RouteInitializer{}
}
