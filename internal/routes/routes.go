package routes

import (
	"finance-processing/internal/handlers"
	"finance-processing/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RegisterRoutes(
	app *fiber.App,
	h *handlers.Handlers,
	mw *middleware.Middleware,
) {

	// default middleware provided by gofibeer
	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// globaal logging
	app.Use(mw.Logging())

	// golbal error handling
	app.Use(mw.ErrorHandler())

	public := app.Group("/")

	public.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Hello from server")
	})

	AuthRoutes(app, h)

	// Adding middlewares
	api := app.Group("/api", mw.Auth())
	TransactionRoutes(api, h)
	DashboardRoutes(api, h)
	UserRoutes(api, h)
}
