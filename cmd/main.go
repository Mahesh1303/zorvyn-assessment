package main

import (
	"finance-processing/internal/config"
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

	// loading registries for db operations

	// loading services for the operations

	// loading handlers

	// loading middlewares

	// loading routes

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("config loaded successfully", cfg.App.Port)

	app.Listen(":8080")
}
