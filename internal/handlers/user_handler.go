package handlers

import (
	"errors"
	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
	"finance-processing/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
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

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     models.UserRole(req.Role),
	}

	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}

	if err := h.service.CreateUser(c.Context(), actor, user); err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "user already exists"})
		}
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "invalid") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created successfully"})
}

func (h *UserHandler) ChangeRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user id is required"})
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if body.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "role is required"})
	}

	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}

	if err := h.service.ChangeRole(c.Context(), actor, userID, body.Role); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "role updated"})
}

func (h *UserHandler) SetActive(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user id is required"})
	}

	var body struct {
		Active bool `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}

	if err := h.service.SetActive(c.Context(), actor, userID, body.Active); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user status updated"})
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}

	users, err := h.service.ListUsers(c.Context(), actor)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	response := make([]models.UserResponse, len(users))
	for i, u := range users {
		response[i] = u.ToResponse()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": response})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(policy.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user context"})
	}
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user id is required"})
	}

	user, err := h.service.GetUser(c.Context(), actor, userID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(user.ToResponse())
}
