// internal/handlers/transaction_handler.go
package handlers

import (
	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
	"finance-processing/internal/services"
	"strconv"

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

	actor := c.Locals("user").(policy.User)

	err := h.service.CreateTransaction(c.Context(), actor, &tx)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "transaction created",
	})
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	actor := c.Locals("user").(policy.User)

	tx, err := h.service.GetTransaction(c.Context(), actor, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tx)
}

func (h *TransactionHandler) ListTransactions(c *fiber.Ctx) error {
	actor := c.Locals("user").(policy.User)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	filter := repository.RecordFilter{
		Type:     c.Query("type"),
		Category: c.Query("category"),
		From:     c.Query("from"),
		To:       c.Query("to"),
		Limit:    limit,
		Offset:   offset,
	}

	data, err := h.service.ListTransaction(c.Context(), actor, filter)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *TransactionHandler) UpdateTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	actor := c.Locals("user").(policy.User)

	var updates map[string]any
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	tx, err := h.service.UpdateTransaction(c.Context(), actor, id, updates)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tx)
}

func (h *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	actor := c.Locals("user").(policy.User)

	if err := h.service.DeleteTransaction(c.Context(), actor, id); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "deleted",
	})
}
