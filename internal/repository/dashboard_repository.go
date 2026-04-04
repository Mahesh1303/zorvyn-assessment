package repository

import (
	"context"

	"finance-processing/internal/models"

	"gorm.io/gorm"
)

type DashboardFilter struct {
	From       string   `query:"from"`
	To         string   `query:"to"`
	Type       string   `query:"type"`
	Categories []string `query:"category"`
}

type Summary struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	NetBalance    float64 `json:"net_balance"`
}

type CategoryTotal struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Total    float64 `json:"total"`
}

type MonthlyTrend struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type DashboardData struct {
	Summary    Summary                      `json:"summary"`
	Categories []CategoryTotal              `json:"categories"`
	Trends     []MonthlyTrend               `json:"trends"`
	Recent     []models.TransactionResponse `json:"recent"`
}

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) base(ctx context.Context, f DashboardFilter) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&models.Transaction{})

	if f.From != "" {
		q = q.Where("date >= ?", f.From)
	}
	if f.To != "" {
		q = q.Where("date <= ?", f.To)
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if len(f.Categories) > 0 {
		q = q.Where("category IN ?", f.Categories)
	}

	return q
}

func (r *DashboardRepository) GetSummary(ctx context.Context, f DashboardFilter) (*Summary, error) {
	var s Summary
	err := r.base(ctx, f).
		Select(`
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expenses
		`).
		Scan(&s).Error
	if err != nil {
		return nil, err
	}
	s.NetBalance = s.TotalIncome - s.TotalExpenses
	return &s, nil
}

func (r *DashboardRepository) GetCategoryTotals(ctx context.Context, f DashboardFilter) ([]CategoryTotal, error) {
	var totals []CategoryTotal
	err := r.base(ctx, f).
		Select("category, type, SUM(amount) AS total").
		Group("category, type").
		Order("total DESC").
		Scan(&totals).Error
	return totals, err
}

func (r *DashboardRepository) GetMonthlyTrends(ctx context.Context, f DashboardFilter) ([]MonthlyTrend, error) {
	var trends []MonthlyTrend
	err := r.base(ctx, f).
		Select(`
			TO_CHAR(date, 'YYYY-MM') AS month,
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense
		`).
		Group("TO_CHAR(date, 'YYYY-MM')").
		Order("month DESC").
		Limit(12).
		Scan(&trends).Error
	return trends, err
}

func (r *DashboardRepository) GetRecent(ctx context.Context, f DashboardFilter, limit int, offset int) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := r.base(ctx, f).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error
	return txs, err
}

func (r *DashboardRepository) GetDashboard(ctx context.Context, f DashboardFilter, limit int, offset int) (*DashboardData, error) {
	summary, err := r.GetSummary(ctx, f)
	if err != nil {
		return nil, err
	}

	categories, err := r.GetCategoryTotals(ctx, f)
	if err != nil {
		return nil, err
	}

	trends, err := r.GetMonthlyTrends(ctx, f)
	if err != nil {
		return nil, err
	}

	recentRaw, err := r.GetRecent(ctx, f, limit, offset)
	if err != nil {
		return nil, err
	}

	recent := make([]models.TransactionResponse, len(recentRaw))
	for i, t := range recentRaw {
		recent[i] = t.ToResponse()
	}

	return &DashboardData{
		Summary:    *summary,
		Categories: categories,
		Trends:     trends,
		Recent:     recent,
	}, nil
}