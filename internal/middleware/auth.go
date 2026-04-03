package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
			})
		}
		// get the payload using jwt
		userID := token

		user, err := m.userRepo.GetByID(c.Context(), userID)
		if err != nil || !user.IsActive {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid or inactive user",
			})
		}
		uid := user.ID
		c.Locals("user_id", uid)
		c.Locals("role", user.Role)

		return c.Next()
	}
}
