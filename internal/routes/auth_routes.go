package routes

import (
	"finance-processing/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, h *handlers.Handlers) {
	public := app.Group("/auth")

	public.Post("/login", h.AuthHandler.Login)
}
