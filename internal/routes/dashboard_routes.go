package routes

import (
	"finance-processing/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(api fiber.Router, h *handlers.Handlers) {

	dashboard := api.Group("/dashboard")

	dashboard.Get("/", h.Dashboard.GetDashboard)

	// GET /api/dashboard/summary?from=2026-01-01&to=2026-03-31&type=expense
	dashboard.Get("/summary", h.Dashboard.GetSummary)

	// GET /api/dashboard/categories?type=expense&type=rent
	dashboard.Get("/categories", h.Dashboard.GetCategoryTotals)

	// /api/dashboard/trends?from=2026-01-01&category=salary&category=rent
	dashboard.Get("/trends", h.Dashboard.GetMonthlyTrends)

	// Paginated recent activity: GET /api/dashboard/recent?limit=10&offset=0
	dashboard.Get("/recent", h.Dashboard.GetRecent)
	dashboard.Get("/analytics", h.Dashboard.GetAnalytics)
}
