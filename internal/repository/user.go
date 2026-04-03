// internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"

	"finance-processing/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {

	if user == nil {
		return errors.New("user is nil")
	}

	return r.db.WithContext(ctx).
		Table("users").
		Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Table("users").
		Where("email = ? AND deleted_at IS NULL", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Table("users").
		Where("deleted_at IS NULL").
		Find(&users).Error

	return users, err
}

func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID string, role string) error {
	return r.db.WithContext(ctx).
		Table("users").
		Where("id = ?", userID).
		Update("role", role).Error
}

func (r *UserRepository) UpdateActive(ctx context.Context, userID string, active bool) error {
	return r.db.WithContext(ctx).
		Table("users").
		Where("id = ?", userID).
		Update("is_active", active).Error
}
