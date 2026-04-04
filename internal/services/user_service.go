package services

import (
	"context"
	"errors"
	auth "finance-processing/internal/lib/utils"
	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
	"fmt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(ctx context.Context, actor policy.User, user *models.User) error {
	if !policy.CanCreateUser(actor) {
		return errors.New("forbidden")
	}
	fmt.Print("userPassword", user.Password)
	if user.Password == "" || user.Email == "" {
		return errors.New("invalid input")
	}
	hashed, err := auth.EncryptPassWord(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.repo.Create(ctx, user)
}

func (s *UserService) ChangeRole(ctx context.Context, actor policy.User, userID string, role string) error {

	if !policy.CanManageUsers(actor) {
		return errors.New("forbidden")
	}
	return s.repo.UpdateRole(ctx, userID, role)
}

func (s *UserService) SetActive(ctx context.Context, actor policy.User, userID string, active bool) error {

	if !policy.CanManageUsers(actor) {
		return errors.New("forbidden")
	}
	return s.repo.UpdateActive(ctx, userID, active)
}

func (s *UserService) ListAnalysts(ctx context.Context, actor policy.User) ([]models.User, error) {

	if !policy.CanManageUsers(actor) {
		return nil, errors.New("forbidden")
	}
	return s.repo.ListByRole(ctx, "analyst")
}

func (s *UserService) ListViewers(ctx context.Context, actor policy.User) ([]models.User, error) {
	if !policy.CanManageUsers(actor) {
		return nil, errors.New("forbidden")
	}
	return s.repo.ListByRole(ctx, "viewer")
}

func (s *UserService) ListUsers(ctx context.Context, actor policy.User) ([]models.User, error) {
	if !policy.CanManageUsers(actor) {
		return nil, errors.New("forbidden")
	}
	return s.repo.List(ctx)
}

func (s *UserService) GetUser(ctx context.Context, actor policy.User, userID string) (*models.User, error) {
	if !policy.CanManageUsers(actor) {
		return nil, errors.New("forbidden")
	}
	return s.repo.GetByID(ctx, userID)
}
