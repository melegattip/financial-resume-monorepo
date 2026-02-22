package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics/domain"
)

// AnalyticsRepo executes aggregation queries over existing transaction tables.
type AnalyticsRepo struct {
	db *gorm.DB
}

// NewAnalyticsRepository creates a new AnalyticsRepo.
func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepo {
	return &AnalyticsRepo{db: db}
}

// categoryRow is used for scanning category-level aggregate results.
type categoryRow struct {
	CategoryID   string  `gorm:"column:category_id"`
	CategoryName string  `gorm:"column:category_name"`
	Amount       float64 `gorm:"column:amount"`
	Count        int     `gorm:"column:cnt"`
}

// monthlyRow is used for scanning month-level aggregate results.
type monthlyRow struct {
	Year   int     `gorm:"column:yr"`
	Month  int     `gorm:"column:mo"`
	Amount float64 `gorm:"column:amount"`
	Count  int     `gorm:"column:cnt"`
}

// GetExpenseSummary returns a summary of expenses for the given user within [from, to].
func (r *AnalyticsRepo) GetExpenseSummary(ctx context.Context, userID string, from, to time.Time, periodLabel string) (*domain.ExpenseSummary, error) {
	// --- Top-level totals ---
	type totalsRow struct {
		Total float64 `gorm:"column:total"`
		Count int     `gorm:"column:cnt"`
	}
	var totals totalsRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(amount), 0) AS total, COUNT(*) AS cnt
		FROM expenses
		WHERE user_id = ?
		  AND transaction_date >= ?
		  AND transaction_date <= ?
		  AND deleted_at IS NULL
	`, userID, from, to).Scan(&totals).Error
	if err != nil {
		return nil, err
	}

	// --- By category ---
	byCategory, err := r.GetExpensesByCategory(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	// --- By month ---
	var monthlyRows []monthlyRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transaction_date)::int  AS yr,
		       EXTRACT(MONTH FROM transaction_date)::int AS mo,
		       COALESCE(SUM(amount), 0)                  AS amount,
		       COUNT(*)                                   AS cnt
		FROM expenses
		WHERE user_id = ?
		  AND transaction_date >= ?
		  AND transaction_date <= ?
		  AND deleted_at IS NULL
		GROUP BY yr, mo
		ORDER BY yr ASC, mo ASC
	`, userID, from, to).Scan(&monthlyRows).Error
	if err != nil {
		return nil, err
	}

	byMonth := make([]domain.MonthlySummary, len(monthlyRows))
	for i, r := range monthlyRows {
		byMonth[i] = domain.MonthlySummary{Year: r.Year, Month: r.Month, Amount: r.Amount, Count: r.Count}
	}

	avg := 0.0
	if totals.Count > 0 {
		avg = totals.Total / float64(totals.Count)
	}

	return &domain.ExpenseSummary{
		TotalAmount:   totals.Total,
		Count:         totals.Count,
		AverageAmount: avg,
		ByCategory:    byCategory,
		ByMonth:       byMonth,
		Period:        periodLabel,
	}, nil
}

// GetIncomeSummary returns a summary of incomes for the given user within [from, to].
func (r *AnalyticsRepo) GetIncomeSummary(ctx context.Context, userID string, from, to time.Time, periodLabel string) (*domain.IncomeSummary, error) {
	// --- Top-level totals ---
	type totalsRow struct {
		Total float64 `gorm:"column:total"`
		Count int     `gorm:"column:cnt"`
	}
	var totals totalsRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(amount), 0) AS total, COUNT(*) AS cnt
		FROM incomes
		WHERE user_id = ?
		  AND received_date >= ?
		  AND received_date <= ?
		  AND deleted_at IS NULL
	`, userID, from, to).Scan(&totals).Error
	if err != nil {
		return nil, err
	}

	// --- By month ---
	var monthlyRows []monthlyRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM received_date)::int  AS yr,
		       EXTRACT(MONTH FROM received_date)::int AS mo,
		       COALESCE(SUM(amount), 0)               AS amount,
		       COUNT(*)                                AS cnt
		FROM incomes
		WHERE user_id = ?
		  AND received_date >= ?
		  AND received_date <= ?
		  AND deleted_at IS NULL
		GROUP BY yr, mo
		ORDER BY yr ASC, mo ASC
	`, userID, from, to).Scan(&monthlyRows).Error
	if err != nil {
		return nil, err
	}

	byMonth := make([]domain.MonthlySummary, len(monthlyRows))
	for i, r := range monthlyRows {
		byMonth[i] = domain.MonthlySummary{Year: r.Year, Month: r.Month, Amount: r.Amount, Count: r.Count}
	}

	avg := 0.0
	if totals.Count > 0 {
		avg = totals.Total / float64(totals.Count)
	}

	return &domain.IncomeSummary{
		TotalAmount:   totals.Total,
		Count:         totals.Count,
		AverageAmount: avg,
		ByMonth:       byMonth,
		Period:        periodLabel,
	}, nil
}

// GetDashboardSummary returns the aggregated dashboard data for the current month and all-time totals.
func (r *AnalyticsRepo) GetDashboardSummary(ctx context.Context, userID string) (*domain.DashboardSummary, error) {
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// --- Current-month totals ---
	type monthTotals struct {
		Expenses float64 `gorm:"column:expenses"`
		Incomes  float64 `gorm:"column:incomes"`
	}
	var mt monthTotals

	err := r.db.WithContext(ctx).Raw(`
		SELECT
		  COALESCE((SELECT SUM(amount) FROM expenses
		            WHERE user_id = ? AND transaction_date >= ? AND deleted_at IS NULL), 0) AS expenses,
		  COALESCE((SELECT SUM(amount) FROM incomes
		            WHERE user_id = ? AND received_date >= ? AND deleted_at IS NULL), 0) AS incomes
	`, userID, monthStart, userID, monthStart).Scan(&mt).Error
	if err != nil {
		return nil, err
	}

	// --- All-time totals ---
	type allTimeTotals struct {
		TotalExpenses float64 `gorm:"column:total_expenses"`
		TotalIncomes  float64 `gorm:"column:total_incomes"`
	}
	var att allTimeTotals
	err = r.db.WithContext(ctx).Raw(`
		SELECT
		  COALESCE((SELECT SUM(amount) FROM expenses WHERE user_id = ? AND deleted_at IS NULL), 0) AS total_expenses,
		  COALESCE((SELECT SUM(amount) FROM incomes  WHERE user_id = ? AND deleted_at IS NULL), 0) AS total_incomes
	`, userID, userID).Scan(&att).Error
	if err != nil {
		return nil, err
	}

	// --- Top categories (current month, top 5) ---
	topCategories, err := r.GetExpensesByCategory(ctx, userID, monthStart, now)
	if err != nil {
		return nil, err
	}
	if len(topCategories) > 5 {
		topCategories = topCategories[:5]
	}

	// --- Recent expenses (last 5) ---
	type recentExpenseRow struct {
		ID              string    `gorm:"column:id"`
		Amount          float64   `gorm:"column:amount"`
		Description     string    `gorm:"column:description"`
		CategoryName    string    `gorm:"column:category_name"`
		TransactionDate time.Time `gorm:"column:transaction_date"`
	}
	var recentExpRows []recentExpenseRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT e.id, e.amount, e.description, COALESCE(c.name, '') AS category_name, e.transaction_date
		FROM expenses e
		LEFT JOIN categories c ON c.id = e.category_id AND c.deleted_at IS NULL
		WHERE e.user_id = ? AND e.deleted_at IS NULL
		ORDER BY e.transaction_date DESC, e.created_at DESC
		LIMIT 5
	`, userID).Scan(&recentExpRows).Error
	if err != nil {
		return nil, err
	}

	recentExpenses := make([]domain.RecentItem, len(recentExpRows))
	for i, row := range recentExpRows {
		recentExpenses[i] = domain.RecentItem{
			ID:          row.ID,
			Amount:      row.Amount,
			Description: row.Description,
			Category:    row.CategoryName,
			Date:        row.TransactionDate,
		}
	}

	// --- Recent incomes (last 5) ---
	type recentIncomeRow struct {
		ID           string    `gorm:"column:id"`
		Amount       float64   `gorm:"column:amount"`
		Description  string    `gorm:"column:description"`
		ReceivedDate time.Time `gorm:"column:received_date"`
	}
	var recentIncRows []recentIncomeRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT id, amount, description, received_date
		FROM incomes
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY received_date DESC, created_at DESC
		LIMIT 5
	`, userID).Scan(&recentIncRows).Error
	if err != nil {
		return nil, err
	}

	recentIncomes := make([]domain.RecentItem, len(recentIncRows))
	for i, row := range recentIncRows {
		recentIncomes[i] = domain.RecentItem{
			ID:          row.ID,
			Amount:      row.Amount,
			Description: row.Description,
			Date:        row.ReceivedDate,
		}
	}

	// --- Savings rate ---
	savingsRate := 0.0
	if mt.Incomes > 0 {
		savingsRate = (mt.Incomes - mt.Expenses) / mt.Incomes * 100
	}

	return &domain.DashboardSummary{
		CurrentMonthExpenses: mt.Expenses,
		CurrentMonthIncomes:  mt.Incomes,
		CurrentMonthBalance:  mt.Incomes - mt.Expenses,
		TotalExpenses:        att.TotalExpenses,
		TotalIncomes:         att.TotalIncomes,
		SavingsRate:          savingsRate,
		TopCategories:        topCategories,
		RecentExpenses:       recentExpenses,
		RecentIncomes:        recentIncomes,
		UpdatedAt:            now,
	}, nil
}

// GetMonthlyExpenses returns expense totals grouped by month for the last N months.
func (r *AnalyticsRepo) GetMonthlyExpenses(ctx context.Context, userID string, months int) ([]domain.MonthlySummary, error) {
	from := time.Now().UTC().AddDate(0, -months, 0)

	var rows []monthlyRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM transaction_date)::int  AS yr,
		       EXTRACT(MONTH FROM transaction_date)::int AS mo,
		       COALESCE(SUM(amount), 0)                  AS amount,
		       COUNT(*)                                   AS cnt
		FROM expenses
		WHERE user_id = ?
		  AND transaction_date >= ?
		  AND deleted_at IS NULL
		GROUP BY yr, mo
		ORDER BY yr ASC, mo ASC
	`, userID, from).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.MonthlySummary, len(rows))
	for i, row := range rows {
		result[i] = domain.MonthlySummary{Year: row.Year, Month: row.Month, Amount: row.Amount, Count: row.Count}
	}
	return result, nil
}

// GetMonthlyIncomes returns income totals grouped by month for the last N months.
func (r *AnalyticsRepo) GetMonthlyIncomes(ctx context.Context, userID string, months int) ([]domain.MonthlySummary, error) {
	from := time.Now().UTC().AddDate(0, -months, 0)

	var rows []monthlyRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXTRACT(YEAR FROM received_date)::int  AS yr,
		       EXTRACT(MONTH FROM received_date)::int AS mo,
		       COALESCE(SUM(amount), 0)               AS amount,
		       COUNT(*)                                AS cnt
		FROM incomes
		WHERE user_id = ?
		  AND received_date >= ?
		  AND deleted_at IS NULL
		GROUP BY yr, mo
		ORDER BY yr ASC, mo ASC
	`, userID, from).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.MonthlySummary, len(rows))
	for i, row := range rows {
		result[i] = domain.MonthlySummary{Year: row.Year, Month: row.Month, Amount: row.Amount, Count: row.Count}
	}
	return result, nil
}

// GetTransactionsForReport returns a lightweight list of expenses and incomes for [from, to],
// containing only the fields needed for report generation (id, category_id, type, amount).
func (r *AnalyticsRepo) GetTransactionsForReport(ctx context.Context, userID string, from, to time.Time) ([]domain.ReportTransaction, error) {
	type row struct {
		ID         string  `gorm:"column:id"`
		CategoryID string  `gorm:"column:category_id"`
		TxType     string  `gorm:"column:tx_type"`
		Amount     float64 `gorm:"column:amount"`
	}

	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, COALESCE(category_id, '') AS category_id, 'expense' AS tx_type, amount
		FROM expenses
		WHERE user_id = ?
		  AND transaction_date >= ?
		  AND transaction_date <= ?
		  AND deleted_at IS NULL
		UNION ALL
		SELECT id, '' AS category_id, 'income' AS tx_type, amount
		FROM incomes
		WHERE user_id = ?
		  AND received_date >= ?
		  AND received_date <= ?
		  AND deleted_at IS NULL
	`, userID, from, to, userID, from, to).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.ReportTransaction, len(rows))
	for i, row := range rows {
		result[i] = domain.ReportTransaction{
			ID:         row.ID,
			CategoryID: row.CategoryID,
			Type:       row.TxType,
			Amount:     row.Amount,
		}
	}
	return result, nil
}

// GetExpensesByCategory returns expense totals grouped by category for [from, to].
// The results are sorted by amount descending and include percentage of total.
func (r *AnalyticsRepo) GetExpensesByCategory(ctx context.Context, userID string, from, to time.Time) ([]domain.CategorySummary, error) {
	var rows []categoryRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT e.category_id,
		       COALESCE(c.name, 'Unknown') AS category_name,
		       COALESCE(SUM(e.amount), 0)  AS amount,
		       COUNT(*)                    AS cnt
		FROM expenses e
		LEFT JOIN categories c ON c.id = e.category_id AND c.deleted_at IS NULL
		WHERE e.user_id = ?
		  AND e.transaction_date >= ?
		  AND e.transaction_date <= ?
		  AND e.deleted_at IS NULL
		GROUP BY e.category_id, c.name
		ORDER BY amount DESC
	`, userID, from, to).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Calculate total for percentage computation.
	total := 0.0
	for _, row := range rows {
		total += row.Amount
	}

	result := make([]domain.CategorySummary, len(rows))
	for i, row := range rows {
		pct := 0.0
		if total > 0 {
			pct = row.Amount / total * 100
		}
		result[i] = domain.CategorySummary{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			Amount:       row.Amount,
			Count:        row.Count,
			Percentage:   pct,
		}
	}
	return result, nil
}
