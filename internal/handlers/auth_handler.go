// internal/handlers/auth_handler.go
package handlers

import (
	"finance-processing/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {

	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body req

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	token, err := h.service.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
