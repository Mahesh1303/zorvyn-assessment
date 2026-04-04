package routes

import (
	"finance-processing/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router, h *handlers.Handlers) {

	users := api.Group("/users")

	users.Post("/", h.User.CreateUser)
	users.Get("/", h.User.ListUsers)
	users.Get("/:id", h.User.GetUser)
	users.Patch("/:id/role", h.User.ChangeRole)
	users.Patch("/:id/status", h.User.SetActive)
}
