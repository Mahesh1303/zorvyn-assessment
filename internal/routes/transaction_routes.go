package routes

import (
	"finance-processing/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(api fiber.Router, h *handlers.Handlers) {

	transactions := api.Group("/transactions")

	transactions.Post("/", h.Transaction.CreateTransaction)
	// Paginated list: GET /api/transactions?limit=10&offset=0
	transactions.Get("/", h.Transaction.ListTransactions)
	transactions.Get("/:id", h.Transaction.GetTransaction)
	transactions.Put("/:id", h.Transaction.UpdateTransaction)
	transactions.Delete("/:id", h.Transaction.DeleteTransaction)
}
