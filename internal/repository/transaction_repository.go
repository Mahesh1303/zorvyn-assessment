package repository

import (
	"context"
	"finance-processing/internal/models"

	"gorm.io/gorm"
)

type RecordFilter struct {
	Type     string
	Category string
	From     string
	To       string
	Limit    int
	Offset   int
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, rec *models.Transaction) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id string) (*models.Transaction, error) {
	var rec models.Transaction
	err := r.db.WithContext(ctx).First(&rec, "id = ?", id).Error
	return &rec, err
}

func (r *TransactionRepository) ListTransactions(ctx context.Context, filter RecordFilter) ([]*models.Transaction, error) {
	var records []*models.Transaction

	q := r.db.WithContext(ctx)
	if filter.Type != "" {
		q = q.Where("type = ?", filter.Type)
	}
	if filter.Category != "" {
		q = q.Where("category = ?", filter.Category)
	}
	if filter.From != "" {
		q = q.Where("date >= ?", filter.From)
	}
	if filter.To != "" {
		q = q.Where("date <= ?", filter.To)
	}
	err := q.Order("date desc").Find(&records).Error
	return records, err
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, id string, updates map[string]any) (*models.Transaction, error) {
	err := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return r.GetTransactionByID(ctx, id)
}

func (r *TransactionRepository) Delete(ctx context.Context, id string) error {
	// gorm.DeletedAt on the model makes this a soft delete automatically
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Transaction{}).Error
}
