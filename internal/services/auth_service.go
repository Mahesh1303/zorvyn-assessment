package services

import (
	"context"
	"errors"

	auth "finance-processing/internal/lib/utils"
	"finance-processing/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwt      *auth.JWTManager
}

func NewAuthService(r *repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo: r,
		jwt:      jwt,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := auth.ComparePassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}
