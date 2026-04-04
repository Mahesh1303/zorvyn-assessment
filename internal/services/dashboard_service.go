package services

import (
	"context"
	"errors"

	"finance-processing/internal/policy"
	"finance-processing/internal/repository"
)

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