package handlers

import (
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
	"finance-processing/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(s *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func parseFilter(c *fiber.Ctx) repository.DashboardFilter {
	var filter repository.DashboardFilter
	_ = c.QueryParser(&filter)
	return filter
}

func parsePagination(c *fiber.Ctx) (limit int, offset int, err error) {
	limit, err = strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}
	offset, err = strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset, nil
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}
	filter := parseFilter(c)
	limit, offset, _ := parsePagination(c)

	data, err := h.service.GetDashboard(c.Context(), actor, filter, limit, offset)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// GET /api/dashboard/summary?from=2026-01-01&to=2026-03-31&type=expense

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}

	summary, err := h.service.GetSummary(c.Context(), actor, parseFilter(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": summary})
}

// GET /api/dashboard/categories?type=expense&from=2026-01-01
func (h *DashboardHandler) GetCategoryTotals(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}

	totals, err := h.service.GetCategoryTotals(c.Context(), actor, parseFilter(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": totals})
}

// GET /api/dashboard/trends?from=2026-01-01&category=salary&category=rent
func (h *DashboardHandler) GetMonthlyTrends(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}

	trends, err := h.service.GetMonthlyTrends(c.Context(), actor, parseFilter(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": trends})
}

// GET /api/dashboard/recent?limit=10&offset=0
func (h *DashboardHandler) GetRecent(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}
	filter := parseFilter(c)
	limit, offset, _ := parsePagination(c)

	recent, err := h.service.GetRecent(c.Context(), actor, filter, limit, offset)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": recent})
}

// GET /api/dashboard/analytics?from=2026-01-01&to=2026-03-31
func (h *DashboardHandler) GetAnalytics(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}
	filter := parseFilter(c)
	limit, offset, _ := parsePagination(c)

	data, err := h.service.GetAnalytics(c.Context(), actor, filter, limit, offset)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}
