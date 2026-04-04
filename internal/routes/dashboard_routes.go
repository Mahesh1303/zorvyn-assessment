package routes

import (
	"finance-processing/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(api fiber.Router, h *handlers.Handlers) {

	dashboard := api.Group("/dashboard")

	dashboard.Get("/", h.Dashboard.GetDashboard)
	dashboard.Get("/summary", h.Dashboard.GetSummary)
	dashboard.Get("/categories", h.Dashboard.GetCategoryTotals)
	dashboard.Get("/trends", h.Dashboard.GetMonthlyTrends)
}
