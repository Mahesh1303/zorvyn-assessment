// internal/handlers/dashboard_handler.go
package handlers

import (
	"finance-processing/internal/services"

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
	userID := c.Locals("user_id").(string)

	summary, err := h.service.GetSummary(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch summary",
		})
	}

	return c.JSON(summary)
}
