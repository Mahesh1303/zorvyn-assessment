package handlers

import (
	"strings"

	"finance-processing/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) LoginUser(c *fiber.Ctx) error {

	var body LoginRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body"})
	}

	body.Email = strings.ToLower(strings.TrimSpace(body.Email))

	if body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "email and password are required",
		})
	}
	if len(body.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{
			"error": "password must be at least 6 characters",
		})
	}

	token, err := h.service.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": token,
	})
}

func (h *AuthHandler) RegisterAdmin(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name, email and password are required",
		})
	}

	user, err := h.service.RegisterAdmin(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "admin registered successfully",
		"data":    user.ToResponse(),
	})
}
