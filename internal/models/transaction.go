package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecordType string

const (
	RecordIncome  RecordType = "income"
	RecordExpense RecordType = "expense"
)

type Transaction struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CreatedBy uuid.UUID  `gorm:"type:uuid;not null;index" json:"created_by"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`

	Amount      float64    `gorm:"type:numeric(12,2);not null;check:amount > 0" json:"amount"`
	Type        RecordType `gorm:"type:record_type;not null"                    json:"type"`
	Category    string     `gorm:"not null;index"                               json:"category"`
	Description string     `json:"description"`
	Date        time.Time  `gorm:"type:date;not null;index"                     json:"date"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
