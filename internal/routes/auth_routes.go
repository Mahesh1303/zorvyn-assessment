package routes

import (
	"finance-processing/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, h *handlers.Handlers) {
	auth := app.Group("/auth")

	auth.Post("/login", h.AuthHandler.LoginUser)
	auth.Post("/register-admin", h.AuthHandler.RegisterAdmin)
}
