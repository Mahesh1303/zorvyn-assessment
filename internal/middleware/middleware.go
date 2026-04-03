package middleware

import (
	"finance-processing/internal/repository"

	"github.com/rs/zerolog"
)

type Middleware struct {
	userRepo *repository.UserRepository
	logger   zerolog.Logger
}

func NewMiddleware(userRepo *repository.UserRepository, logger zerolog.Logger) *Middleware {
	return &Middleware{
		userRepo: userRepo,
		logger:   logger,
	}
}
