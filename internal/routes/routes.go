package routes

import (
	"finance-processing/internal/handlers"
	"finance-processing/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	h *handlers.Handlers,
	mw *middleware.Middleware,
) {

	public := app.Group("/")

	public.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	public.Get("/login", h.AuthHandler.Login)

	// protected Routes for all Admin, Analyst and Viewer
	api := app.Group("/api", mw.Auth())

	// for transactions
	transactions := api.Group("/transaction")

	transactions.Post("/", h.Transaction.CreateTransaction)
	transactions.Get("/", h.Transaction.ListTransactions)
	transactions.Get("/:id", h.Transaction.GetTransaction)
	transactions.Put("/:id", h.Transaction.UpdateTransaction)
	transactions.Delete("/:id", h.Transaction.DeleteTransaction)

	// for dashboard
	dashboard := api.Group("/dashboard")

	dashboard.Get("/summary", h.Dashboard.GetSummary)
	dashboard.Get("/categories", h.Dashboard.GetCategoryTotals)
	dashboard.Get("/trends", h.Dashboard.GetMonthlyTrends)
	dashboard.Get("/recent", h.Dashboard.GetRecent)

	// for Users
	users := api.Group("/users")

	users.Post("/", h.User.CreateUser)
	users.Patch("/:id/role", h.User.ChangeRole)
	users.Patch("/:id/status", h.User.SetActive)

}
