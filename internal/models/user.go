package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleAnalyst UserRole = "analyst"
	RoleViewer  UserRole = "viewer"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Name     string `gorm:"not null"             json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null"             json:"-"`

	Role     UserRole `gorm:"type:user_role;default:viewer;not null" json:"role"`
	IsActive bool     `gorm:"default:true;not null"                  json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
