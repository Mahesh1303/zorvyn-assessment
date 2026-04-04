package main

import (
	"finance-processing/internal/server"
	"os"

	"github.com/rs/zerolog"
)

func main() {

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	srv, err := server.New(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize server")
	}

	if err := srv.Start(); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
