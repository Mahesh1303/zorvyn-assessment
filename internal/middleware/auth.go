package middleware

import (
	"finance-processing/internal/policy"
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

		claims, err := m.jwt.Verify(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		userID := claims.UserID.String()

		user, err := m.userRepo.GetByID(c.Context(), userID)
		if err != nil || !user.IsActive {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid or inactive user",
			})
		}

		// setting actor
		actor := policy.User{
			ID:   user.ID,
			Role: string(user.Role),
		}
		c.Locals("user", actor)

		c.Locals("user_id", actor.ID)
		c.Locals("role", actor.Role)

		return c.Next()
	}
}
