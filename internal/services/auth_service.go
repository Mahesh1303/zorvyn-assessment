// internal/services/auth_service.go
package services

import (
	"context"
	"errors"

	"finance-processing/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(r *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: r}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// ⚠️ plain comparison for assignment (mention hashing in README)
	if user.Password != password {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", errors.New("user inactive")
	}

	// Generating Token
	return user.Name, err
}
