package main

import (
	"finance-processing/internal/config"
	db "finance-processing/internal/database"
	"finance-processing/internal/handlers"
	"finance-processing/internal/middleware"
	"finance-processing/internal/repository"
	"finance-processing/internal/routes"
	"finance-processing/internal/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func main() {
	app := fiber.New()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Loading the configs
	cfg, err := config.LoadConfig(logger)

	// loading the databasee
	gormDB, err := db.Connect(cfg.DB.URL, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	if err := db.RunMigrations(cfg.DB.URL, logger); err != nil {
		logger.Fatal().Err(err).Msg("migrations failed")
	}

	// loading registries for db operations
	repos := repository.NewRepositories(gormDB)

	// loading services for the operations and passing the repositories to the services so that it can call them
	svcs := services.NewServices(repos)

	// loading handlers and passing services to it so that our handlers can call those
	// so ourr flow becomes routes-->middlewares-->handlers--> services-->repositories
	h := handlers.NewHandlers(svcs)

	// loading middlewares
	mw := middleware.NewMiddleware(repos.User, logger)

	// loading routes
	routes.RegisterRoutes(app, h, mw)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("config loaded successfully", cfg.App.Port)

	app.Listen(":8080")
}
