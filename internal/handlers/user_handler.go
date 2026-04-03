// internal/handlers/dashboard_handler.go
package handlers

import (
	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	actor := policy.User{
		Role: c.Locals("role").(string),
	}

	// 🔹 Calling  Creating service which is in the UserService
	if err := h.service.CreateUser(c.Context(), actor, &user); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "user created successfully",
	})
}

// 🔹 Get Summary
func (h *UserHandler) ChangeRole(c *fiber.Ctx) error {

	type req struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}

	var body req

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	actor := policy.User{
		Role: c.Locals("role").(string),
	}

	err := h.service.ChangeRole(c.Context(), actor, body.UserID, body.Role)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "role updated"})
}

func (h *UserHandler) SetActive(c *fiber.Ctx) error {

	type req struct {
		UserID string `json:"user_id"`
		Active bool   `json:"active"`
	}

	var body req

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	actor := policy.User{
		Role: c.Locals("role").(string),
	}

	err := h.service.SetActive(c.Context(), actor, body.UserID, body.Active)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "user status updated"})
}
