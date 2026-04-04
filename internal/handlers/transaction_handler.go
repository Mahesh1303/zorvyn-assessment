// internal/handlers/transaction_handler.go
package handlers

import (
	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
	"finance-processing/internal/services"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionUpdateRequest struct {
	Amount      *float64 `json:"amount"`
	Category    *string  `json:"category"`
	Description *string  `json:"description"`
	Date        *string  `json:"date"`
}

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(s *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

// 🔹 Create Transaction
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req CreateTransactionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid date format (use YYYY-MM-DD)",
		})
	}

	tx := &models.Transaction{
		Amount:      req.Amount,
		Type:        models.RecordType(req.Type),
		Category:    req.Category,
		Description: req.Description,
		Date:        parsedDate,
	}

	actor := c.Locals("user").(policy.User)

	err = h.service.CreateTransaction(c.Context(), actor, tx)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "amount") {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
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
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "transaction not found"})
		}
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tx)
}

func (h *TransactionHandler) ListTransactions(c *fiber.Ctx) error {
	actor := c.Locals("user").(policy.User)

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid limit parameter"})
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid offset parameter"})
	}
	if offset < 0 {
		offset = 0
	}

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
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
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

	var req TransactionUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Amount == nil && req.Category == nil && req.Description == nil && req.Date == nil {
		return c.Status(400).JSON(fiber.Map{"error": "at least one field must be provided"})
	}

	updates := make(map[string]any)
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "amount must be positive"})
		}
		updates["amount"] = *req.Amount
	}
	if req.Category != nil && *req.Category != "" {
		updates["category"] = *req.Category
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Date != nil && *req.Date != "" {
		updates["date"] = *req.Date
	}

	tx, err := h.service.UpdateTransaction(c.Context(), actor, id, updates)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "transaction not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(tx)
}

func (h *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	actor := c.Locals("user").(policy.User)

	if err := h.service.DeleteTransaction(c.Context(), actor, id); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "transaction not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message": "deleted",
	})
}
