package repository

import (
	"context"
	"finance-processing/internal/models"

	"gorm.io/gorm"
)

type Summary struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	NetBalance    float64 `json:"net_balance"`
}

type CategoryTotal struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type MonthlyTrend struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetSummary(ctx context.Context) (*Summary, error) {
	var s Summary
	err := r.db.WithContext(ctx).
		Table("financial_records").
		Select(`
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expenses
		`).
		Where("deleted_at IS NULL").
		Scan(&s).Error
	s.NetBalance = s.TotalIncome - s.TotalExpenses
	return &s, err
}

func (r *DashboardRepository) GetCategoryTotals(ctx context.Context) ([]*CategoryTotal, error) {
	var totals []*CategoryTotal
	err := r.db.WithContext(ctx).
		Table("financial_records").
		Select("category, SUM(amount) AS total").
		Where("deleted_at IS NULL").
		Group("category").
		Order("total DESC").
		Scan(&totals).Error
	return totals, err
}

func (r *DashboardRepository) GetMonthlyTrends(ctx context.Context) ([]*MonthlyTrend, error) {
	var trends []*MonthlyTrend
	err := r.db.WithContext(ctx).
		Table("financial_records").
		Select(`
			TO_CHAR(date, 'YYYY-MM') AS month,
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense
		`).
		Where("deleted_at IS NULL").
		Group("TO_CHAR(date, 'YYYY-MM')").
		Order("month DESC").
		Limit(12).
		Scan(&trends).Error
	return trends, err
}

func (r *DashboardRepository) GetRecent(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Transaction, error) {

	var txs []models.Transaction

	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error

	return txs, err
}
