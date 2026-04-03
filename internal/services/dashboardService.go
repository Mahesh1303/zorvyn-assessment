// internal/services/dashboard_service.go
package services

import (
	"context"

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

func (s *DashboardService) GetSummary(ctx context.Context, userID string) (*Summary, error) {
	// 🔥 for now assume token = userID
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
