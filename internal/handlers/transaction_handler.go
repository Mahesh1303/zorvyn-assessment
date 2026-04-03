// internal/handlers/transaction_handler.go
package handlers

import (
	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/services"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(s *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

// 🔹 Create Transaction
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var tx models.Transaction

	if err := c.BodyParser(&tx); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// getting user from middleware
	user := policy.User{
		Role: c.Locals("role").(string),
	}

	err := h.service.CreateTransaction(c.Context(), user, &tx)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "transaction created",
	})
}
