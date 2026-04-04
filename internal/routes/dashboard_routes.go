package routes

import (
	"finance-processing/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(api fiber.Router, h *handlers.Handlers) {

	dashboard := api.Group("/dashboard")

	dashboard.Get("/", h.Dashboard.GetDashboard)
	dashboard.Get("/summary", h.Dashboard.GetSummary)
	// GET /api/dashboard/categories?type=expense&from=2026-01-01
	dashboard.Get("/categories", h.Dashboard.GetCategoryTotals)
	dashboard.Get("/trends", h.Dashboard.GetMonthlyTrends)
	// Paginated recent activity: GET /api/dashboard/recent?limit=10&offset=0
	dashboard.Get("/recent", h.Dashboard.GetRecent)
	dashboard.Get("/analytics", h.Dashboard.GetAnalytics)
}
