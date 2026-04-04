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
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null;index"                       json:"created_by"`
	UpdatedBy   *uuid.UUID     `gorm:"type:uuid"                                      json:"updated_by,omitempty"`
	Amount      float64        `gorm:"type:numeric(12,2);not null"                    json:"amount"`
	Type        RecordType     `gorm:"type:record_type;not null"                      json:"type"`
	Category    string         `gorm:"not null"                                       json:"category"`
	Description string         `json:"description,omitempty"`
	Date        time.Time      `gorm:"type:date;not null"                             json:"date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                                          json:"-"`
}

func (Transaction) TableName() string {
	return "financial_records"
}

type TransactionResponse struct {
	ID          uuid.UUID  `json:"id"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
	Amount      float64    `json:"amount"`
	Type        RecordType `json:"type"`
	Category    string     `json:"category"`
	Description string     `json:"description,omitempty"`
	Date        time.Time  `json:"date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (t *Transaction) ToResponse() TransactionResponse {
	return TransactionResponse{
		ID:          t.ID,
		CreatedBy:   t.CreatedBy,
		UpdatedBy:   t.UpdatedBy,
		Amount:      t.Amount,
		Type:        t.Type,
		Category:    t.Category,
		Description: t.Description,
		Date:        t.Date,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
