package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		m.logger.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", time.Since(start))

		return err
	}
}
