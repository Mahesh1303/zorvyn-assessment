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

	public.Get("/login")
	// protected Routes for all Admin, Analyst and Viewer
	api := app.Group("/api", mw.Auth())

	// for transactions
	transactions := api.Group("/transaction")

	transactions.Post("/", h.Transaction.CreateTransaction) //create Transaction only Admin

	// for dashboard
	dashboard := api.Group("/dashboard")

	dashboard.Get("/summary", h.Dashboard.GetSummary)

	// for Users
	users := api.Group("/users")

	users.Post("/", h.User.CreateUser)

}
