package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics/domain"
)

// AnalyticsRepository defines the data-access operations required by the handler.
type AnalyticsRepository interface {
	GetExpenseSummary(ctx context.Context, userID string, from, to time.Time, periodLabel string) (*domain.ExpenseSummary, error)
	GetIncomeSummary(ctx context.Context, userID string, from, to time.Time, periodLabel string) (*domain.IncomeSummary, error)
	GetDashboardSummary(ctx context.Context, userID string) (*domain.DashboardSummary, error)
	GetMonthlyExpenses(ctx context.Context, userID string, months int) ([]domain.MonthlySummary, error)
	GetMonthlyIncomes(ctx context.Context, userID string, months int) ([]domain.MonthlySummary, error)
	GetExpensesByCategory(ctx context.Context, userID string, from, to time.Time) ([]domain.CategorySummary, error)
	GetTransactionsForReport(ctx context.Context, userID string, from, to time.Time) ([]domain.ReportTransaction, error)
}

// AnalyticsHandler handles all analytics and dashboard HTTP requests.
type AnalyticsHandler struct {
	repo   AnalyticsRepository
	logger zerolog.Logger
}

// NewAnalyticsHandler creates a new AnalyticsHandler.
func NewAnalyticsHandler(repo AnalyticsRepository, logger zerolog.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo, logger: logger}
}

// parsePeriod converts the `period` query parameter (or `from`/`to` params) into a
// concrete time range and a human-readable label.
//
// Supported period values:
//   - this_month     → 1st of current month to now
//   - last_month     → 1st to last day of previous month
//   - this_year      → 1st Jan of current year to now
//   - last_30_days   → now - 30 days to now
//   - last_90_days   → now - 90 days to now
//
// If period is empty, `from` and `to` query params (RFC3339) are used.
func parsePeriod(c *gin.Context) (from, to time.Time, label string, err error) {
	now := time.Now().UTC()
	period := c.DefaultQuery("period", "")

	switch period {
	case "this_month":
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		to = now
		label = "this_month"

	case "last_month":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		from = firstOfThisMonth.AddDate(0, -1, 0)
		to = firstOfThisMonth.Add(-time.Second)
		label = "last_month"

	case "this_year":
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		to = now
		label = "this_year"

	case "last_30_days":
		from = now.AddDate(0, 0, -30)
		to = now
		label = "last_30_days"

	case "last_90_days":
		from = now.AddDate(0, 0, -90)
		to = now
		label = "last_90_days"

	default:
		// Fall back to explicit from/to query params.
		fromStr := c.Query("from")
		toStr := c.Query("to")
		if fromStr == "" || toStr == "" {
			// Default to current month when no parameters given.
			from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			to = now
			label = "this_month"
			return
		}

		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return
		}
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			return
		}
		label = "custom"
	}
	return
}

// GetExpenseSummary handles GET /analytics/expenses
// Query params: period, from, to
func (h *AnalyticsHandler) GetExpenseSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from, to, label, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
		return
	}

	summary, err := h.repo.GetExpenseSummary(c.Request.Context(), userID.(string), from, to, label)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get expense summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expense summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetIncomeSummary handles GET /analytics/incomes
// Query params: period, from, to
func (h *AnalyticsHandler) GetIncomeSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from, to, label, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
		return
	}

	summary, err := h.repo.GetIncomeSummary(c.Request.Context(), userID.(string), from, to, label)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get income summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get income summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetCategoryAnalysis handles GET /analytics/categories
// Query params: period, from, to, type (expenses|incomes — currently only expenses are categorized)
func (h *AnalyticsHandler) GetCategoryAnalysis(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from, to, _, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
		return
	}

	categories, err := h.repo.GetExpensesByCategory(c.Request.Context(), userID.(string), from, to)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get category analysis")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category analysis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"total": len(categories),
	})
}

// GetMonthlyTrends handles GET /analytics/monthly
// Query params: months (default 12)
func (h *AnalyticsHandler) GetMonthlyTrends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	months, err := strconv.Atoi(c.DefaultQuery("months", "12"))
	if err != nil || months < 1 || months > 60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid months parameter, must be between 1 and 60"})
		return
	}

	expenses, err := h.repo.GetMonthlyExpenses(c.Request.Context(), userID.(string), months)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get monthly expenses")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get monthly trends"})
		return
	}

	incomes, err := h.repo.GetMonthlyIncomes(c.Request.Context(), userID.(string), months)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get monthly incomes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get monthly trends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"expenses": expenses,
		"incomes":  incomes,
		"months":   months,
	})
}

// GetFinancialHealth handles GET /insights/financial-health
// Computes a financial health score from the current month's income vs. expense ratio.
func (h *AnalyticsHandler) GetFinancialHealth(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	summary, err := h.repo.GetDashboardSummary(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get financial health")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get financial health"})
		return
	}

	score := 100.0
	if summary.CurrentMonthIncomes > 0 {
		ratio := summary.CurrentMonthExpenses / summary.CurrentMonthIncomes
		switch {
		case ratio >= 1.0:
			score = 20.0
		case ratio >= 0.9:
			score = 40.0
		case ratio >= 0.7:
			score = 60.0
		case ratio >= 0.5:
			score = 80.0
		}
	} else if summary.CurrentMonthExpenses > 0 {
		score = 10.0
	}

	status := "excellent"
	switch {
	case score < 40:
		status = "critical"
	case score < 60:
		status = "poor"
	case score < 80:
		status = "fair"
	case score < 90:
		status = "good"
	}

	savingsRate := summary.SavingsRate
	if savingsRate < 0 {
		savingsRate = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"score":                  score,
		"status":                 status,
		"savings_rate":           savingsRate,
		"current_month_expenses": summary.CurrentMonthExpenses,
		"current_month_incomes":  summary.CurrentMonthIncomes,
		"current_month_balance":  summary.CurrentMonthBalance,
	})
}

// GetReport handles GET /reports
// Query params: start_date, end_date (YYYY-MM-DD).
// Returns total_income, total_expenses, transactions array, and category_summary.
func (h *AnalyticsHandler) GetReport(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	parseDate := func(s string, fallback time.Time) (time.Time, error) {
		if s == "" {
			return fallback, nil
		}
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			t, err = time.Parse(time.RFC3339, s)
		}
		return t, err
	}

	now := time.Now().UTC()
	from, err := parseDate(c.Query("start_date"), time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}
	toRaw, err := parseDate(c.Query("end_date"), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
		return
	}
	// Include the entire end day.
	to := time.Date(toRaw.Year(), toRaw.Month(), toRaw.Day(), 23, 59, 59, 0, time.UTC)

	ctx := c.Request.Context()
	uid := userID.(string)

	expenseSummary, err := h.repo.GetExpenseSummary(ctx, uid, from, to, "custom")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get expense summary for report")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	incomeSummary, err := h.repo.GetIncomeSummary(ctx, uid, from, to, "custom")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get income summary for report")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	transactions, err := h.repo.GetTransactionsForReport(ctx, uid, from, to)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get transactions for report")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	// Build category_summary using the field names the frontend expects.
	type categorySummaryItem struct {
		CategoryID   string  `json:"category_id"`
		CategoryName string  `json:"category_name"`
		TotalAmount  float64 `json:"total_amount"`
		Percentage   float64 `json:"percentage"`
	}
	categorySummary := make([]categorySummaryItem, len(expenseSummary.ByCategory))
	for i, cat := range expenseSummary.ByCategory {
		categorySummary[i] = categorySummaryItem{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			TotalAmount:  cat.Amount,
			Percentage:   cat.Percentage,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_income":     incomeSummary.TotalAmount,
		"total_expenses":   expenseSummary.TotalAmount,
		"transactions":     transactions,
		"category_summary": categorySummary,
	})
}

// GetDashboard handles GET /dashboard
func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	summary, err := h.repo.GetDashboardSummary(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get dashboard summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get dashboard"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
