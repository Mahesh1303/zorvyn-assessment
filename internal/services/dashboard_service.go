// internal/services/dashboard_service.go
package services

import (
	"context"
	"errors"

	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(r *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: r}
}

type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Net          float64 `json:"net"`
}

func (s *DashboardService) GetSummary(
	ctx context.Context,
	actor policy.User,
) (*Summary, error) {

	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden")
	}

	txs, err := s.repo.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &Summary{
		TotalIncome:  txs.TotalIncome,
		TotalExpense: txs.TotalExpenses,
		Net:          txs.NetBalance,
	}, nil
}

func (s *DashboardService) GetRecent(
	ctx context.Context,
	actor policy.User,
	limit int,
	offset int,
) ([]models.Transaction, error) {

	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetRecent(ctx, limit, offset)
}

func (s *DashboardService) GetCategoryTotals(
	ctx context.Context,
	actor policy.User,
) ([]*repository.CategoryTotal, error) {

	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetCategoryTotals(ctx)
}

func (s *DashboardService) GetMonthlyTrends(
	ctx context.Context,
	actor policy.User,
) ([]*repository.MonthlyTrend, error) {

	if !policy.CanViewAnalytics(actor) {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetMonthlyTrends(ctx)
}
