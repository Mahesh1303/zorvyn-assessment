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
) error {

	if !policy.CanViewTransaction(actor) {
		return errors.New("forbidden")
	}
	return s.repo.Delete(ctx, id)
}
func (s *TransactionService) UpdateTransaction(
	ctx context.Context,
	actor policy.User,
	id string,
) error {

	if !policy.CanManageTransaction(actor) {
		return errors.New("forbidden")
	}

	return s.repo.Delete(ctx, id)
}

func (s *TransactionService) ListTransaction(
	ctx context.Context,
	actor policy.User,
	id string,
) error {

	if !policy.CanViewTransaction(actor) {
		return errors.New("forbidden")
	}

	return s.repo.Delete(ctx, id)
}
