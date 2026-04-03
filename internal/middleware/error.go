package middleware

import "github.com/gofiber/fiber/v2"

func (m *Middleware) ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return nil
	}
}
