// internal/handlers/dashboard_handler.go
package handlers

import (
	"finance-processing/internal/policy"
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

// 🔹 Get Summary
func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {

	actor := c.Locals("user").(policy.User)
	summary, err := h.service.GetSummary(c.Context(), actor)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch summary",
		})
	}

	return c.JSON(summary)
}

func (h *DashboardHandler) GetCategoryTotals(c *fiber.Ctx) error {
	actor := c.Locals("user").(policy.User)

	data, err := h.service.GetCategoryTotals(c.Context(), actor)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *DashboardHandler) GetMonthlyTrends(c *fiber.Ctx) error {
	actor := c.Locals("user").(policy.User)

	data, err := h.service.GetMonthlyTrends(c.Context(), actor)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *DashboardHandler) GetRecent(c *fiber.Ctx) error {
	actor := c.Locals("user").(policy.User)

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid limit"})
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid offset"})
	}

	data, err := h.service.GetRecent(c.Context(), actor, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}
