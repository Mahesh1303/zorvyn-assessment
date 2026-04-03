package services

// internal/services/transaction_service.go

import (
	"context"
	"errors"

	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(r *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: r}
}

func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	actor policy.User,
	tx *models.Transaction,
) error {

	if !policy.CanManageTransaction(actor) {
		return errors.New("forbidden")
	}

	if tx.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	tx.CreatedBy = actor.ID

	return s.repo.CreateTransaction(ctx, tx)
}

func (s *TransactionService) DeleteTransaction(
	ctx context.Context,
	actor policy.User,
	id string,
) error {

	if !policy.CanManageTransaction(actor) {
		return errors.New("forbidden")
	}

	return s.repo.Delete(ctx, id)
}

func (s *TransactionService) GetTransaction(
	ctx context.Context,
	actor policy.User,
	id string,
) (*models.Transaction, error) {

	if !policy.CanViewTransaction(actor) {
		return nil, errors.New("forbidden")
	}
	return s.repo.GetTransactionByID(ctx, id)
}

func (s *TransactionService) UpdateTransaction(
	ctx context.Context,
	actor policy.User,
	id string,
	updates map[string]any,
) (*models.Transaction, error) {

	if !policy.CanManageTransaction(actor) {
		return nil, errors.New("forbidden")
	}

	return s.repo.UpdateTransaction(ctx, id, updates)
}

func (s *TransactionService) ListTransaction(
	ctx context.Context,
	actor policy.User,
	filter repository.RecordFilter,
) ([]*models.Transaction, error) {

	if !policy.CanViewTransaction(actor) {
		return nil, errors.New("forbidden")
	}

	return s.repo.ListTransactions(ctx, filter)
}
