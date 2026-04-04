// internal/handlers/handlers.go

package handlers

import (
	"finance-processing/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Transaction *TransactionHandler
	User        *UserHandler
	Dashboard   *DashboardHandler
	AuthHandler *AuthHandler
}

func NewHandlers(s *services.Services) *Handlers {
	return &Handlers{
		Transaction: NewTransactionHandler(s.Transaction),
		User:        NewUserHandler(s.User),
		Dashboard:   NewDashboardHandler(s.Dashboard),
		AuthHandler: NewAuthHandler(s.Auth),
	}
}

func handleServiceError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case strings.HasPrefix(msg, "forbidden"):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": msg})
	case strings.HasPrefix(msg, "not found"):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": msg})
	case strings.HasPrefix(msg, "invalid"):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	case strings.HasPrefix(msg, "already exists"):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": msg})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
