package middleware

import (
	auth "finance-processing/internal/lib/utils"
	"finance-processing/internal/repository"

	"github.com/rs/zerolog"
)

type Middleware struct {
	userRepo *repository.UserRepository
	jwt      *auth.JWTManager
	logger   zerolog.Logger
}

func NewMiddleware(userRepo *repository.UserRepository, jwt *auth.JWTManager, logger zerolog.Logger) *Middleware {
	return &Middleware{
		userRepo: userRepo,
		jwt:      jwt,
		logger:   logger,
	}
}
