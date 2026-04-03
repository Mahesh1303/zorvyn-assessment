package services

import (
	"context"
	"errors"
	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(ctx context.Context, actor policy.User, user *models.User) error {

	if actor.Role != "admin" {
		return errors.New("only admin can create users")
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) ChangeRole(ctx context.Context, actor policy.User, userID string, role string) error {

	if actor.Role != "admin" {
		return errors.New("only admin can change roles")
	}

	return s.repo.UpdateRole(ctx, userID, role)
}

func (s *UserService) SetActive(ctx context.Context, actor policy.User, userID string, active bool) error {

	if actor.Role != "admin" {
		return errors.New("only admin can update user status")
	}

	return s.repo.UpdateActive(ctx, userID, active)
}

func (s *UserService) ListAnalyst(ctx context.Context, actor policy.User) string {
	return "Analyst chi list"
}

func (s *UserService) ListViewers(ctx context.Context, actor policy.User) string {
	return "ali Viewer chi list"
}
