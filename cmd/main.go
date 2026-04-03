package main

import (
	"finance-processing/internal/config"
	"finance-processing/internal/database"
	"finance-processing/internal/handlers"
	"finance-processing/internal/lib/utils"
	"finance-processing/internal/middleware"
	"finance-processing/internal/repository"
	"finance-processing/internal/routes"
	"finance-processing/internal/services"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func main() {
	app := fiber.New()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	gormDB, err := database.Connect(cfg.DB.URL, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect database")
	}

	if err := database.RunMigrations(cfg.DB.URL, logger); err != nil {
		logger.Fatal().Err(err).Msg("migration failed")
	}

	repos := repository.NewRepositories(gormDB)
	jwtManager := utils.NewJWTManager(cfg.Auth.JWTSecret)

	svcs := services.NewServices(repos, jwtManager)
	h := handlers.NewHandlers(svcs)
	mw := middleware.NewMiddleware(repos.User, jwtManager, logger)

	routes.RegisterRoutes(app, h, mw)

	logger.Info().Msgf("server running on port %d", cfg.App.Port)

	if err := app.Listen(":" + strconv.Itoa(cfg.App.Port)); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
