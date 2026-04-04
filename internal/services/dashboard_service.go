package services

import (
	"context"
	"errors"

	"finance-processing/internal/models"
	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
)

type DashboardAnalytics struct {
	Summary      *repository.Summary          `json:"summary"`
	Categories   []repository.CategoryTotal   `json:"categories"`
	Trends       []repository.MonthlyTrend    `json:"trends"`
	Recent       []models.TransactionResponse `json:"recent"`
	ExpenseRatio float64                      `json:"expense_ratio"`
}

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(r *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: r}
}

func (s *DashboardService) GetDashboard(
	ctx context.Context,
	actor policy.User,
	filter repository.DashboardFilter,
	limit int,
	offset int,
) (*repository.DashboardData, error) {
	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden: insufficient permissions")
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetDashboard(ctx, filter, limit, offset)
}

func (s *DashboardService) GetSummary(
	ctx context.Context,
	actor policy.User,
	filter repository.DashboardFilter,
) (*repository.Summary, error) {
	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden: insufficient permissions")
	}
	return s.repo.GetSummary(ctx, filter)
}

func (s *DashboardService) GetCategoryTotals(
	ctx context.Context,
	actor policy.User,
	filter repository.DashboardFilter,
) ([]repository.CategoryTotal, error) {
	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden: insufficient permissions")
	}
	return s.repo.GetCategoryTotals(ctx, filter)
}

func (s *DashboardService) GetMonthlyTrends(
	ctx context.Context,
	actor policy.User,
	filter repository.DashboardFilter,
) ([]repository.MonthlyTrend, error) {
	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden: insufficient permissions")
	}
	return s.repo.GetMonthlyTrends(ctx, filter)
}

func (s *DashboardService) GetRecent(
	ctx context.Context,
	actor policy.User,
	filter repository.DashboardFilter,
	limit int,
	offset int,
) ([]models.TransactionResponse, error) {
	if !policy.CanViewDashboard(actor) {
		return nil, errors.New("forbidden: insufficient permissions")
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	recent, err := s.repo.GetRecent(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]models.TransactionResponse, len(recent))
	for i, tx := range recent {
		response[i] = tx.ToResponse()
	}

	return response, nil
}

func (s *DashboardService) GetAnalytics(
	ctx context.Context,
	actor policy.User,
	filter repository.DashboardFilter,
	limit int,
	offset int,
) (*DashboardAnalytics, error) {
	if !policy.CanViewAnalytics(actor) {
		return nil, errors.New("forbidden: insufficient permissions")
	}

	data, err := s.repo.GetDashboard(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	expenseRatio := 0.0
	if data.Summary.TotalIncome > 0 {
		expenseRatio = data.Summary.TotalExpenses / data.Summary.TotalIncome
	}

	return &DashboardAnalytics{
		Summary:      &data.Summary,
		Categories:   data.Categories,
		Trends:       data.Trends,
		Recent:       data.Recent,
		ExpenseRatio: expenseRatio,
	}, nil
}
