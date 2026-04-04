package services

import (
	"context"
	"errors"

	auth "finance-processing/internal/lib/utils"
	"finance-processing/internal/models"
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

func (s *AuthService) RegisterAdmin(ctx context.Context, name, email, password string) (*models.User, error) {
	var count int64
	s.userRepo.CountAdmins(ctx, &count)
	if count > 0 {
		return nil, errors.New("forbidden: system already initialized")
	}

	hashed, err := auth.EncryptPassWord(password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashed,
		Role:     models.RoleAdmin,
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
