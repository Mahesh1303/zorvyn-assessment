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
	actor := c.Locals("user").(policy.User)
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
	actor := c.Locals("user").(policy.User)

	summary, err := h.service.GetSummary(c.Context(), actor, parseFilter(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": summary})
}

// GET /api/dashboard/categories?type=expense&from=2026-01-01
func (h *DashboardHandler) GetCategoryTotals(c *fiber.Ctx) error {
	actor := c.Locals("user").(policy.User)

	totals, err := h.service.GetCategoryTotals(c.Context(), actor, parseFilter(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": totals})
}

// GET /api/dashboard/trends?from=2026-01-01&category=salary&category=rent
func (h *DashboardHandler) GetMonthlyTrends(c *fiber.Ctx) error {
	actor := c.Locals("user").(policy.User)

	trends, err := h.service.GetMonthlyTrends(c.Context(), actor, parseFilter(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": trends})
}
