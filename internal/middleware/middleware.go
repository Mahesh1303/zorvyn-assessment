package middleware

import (
	"finance-processing/internal/config"
	"finance-processing/internal/repository"

	"github.com/rs/zerolog"
)

type Middleware struct {
	userRepo *repository.UserRepository
	logger   zerolog.Logger
}

func NewMiddleware(userRepo *repository.UserRepository, logger zerolog.Logger, cfg *config.AuthConfig) *Middleware {
	return &Middleware{
		userRepo: userRepo,
		logger:   logger,
	}
}
