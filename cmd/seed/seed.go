package main

import (
	"finance-processing/internal/config"
	"finance-processing/internal/database"
	seeder "finance-processing/internal/seed"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	if err := database.RunMigrations(cfg.DB.URL, logger); err != nil {
		logger.Fatal().Err(err).Msg("migrations failed")
	}

	db, err := database.Connect(cfg.DB.URL, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect db")
	}

	if err := seeder.RunSeed(db, logger); err != nil {
		logger.Fatal().Err(err).Msg("seed failed")
	}

	logger.Info().Msg("seed completed")
}
