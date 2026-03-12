package handlers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics/domain"
)

// AnalyticsRepository defines the data-access operations required by the handler.
type AnalyticsRepository interface {
	GetExpenseSummary(ctx context.Context, tenantID string, from, to time.Time, periodLabel string) (*domain.ExpenseSummary, error)
	GetIncomeSummary(ctx context.Context, tenantID string, from, to time.Time, periodLabel string) (*domain.IncomeSummary, error)
	GetDashboardSummary(ctx context.Context, tenantID string) (*domain.DashboardSummary, error)
	GetMonthlyExpenses(ctx context.Context, tenantID string, months int) ([]domain.MonthlySummary, error)
	GetMonthlyIncomes(ctx context.Context, tenantID string, months int) ([]domain.MonthlySummary, error)
	GetExpensesByCategory(ctx context.Context, tenantID string, from, to time.Time) ([]domain.CategorySummary, error)
	GetTransactionsForReport(ctx context.Context, tenantID string, from, to time.Time) ([]domain.ReportTransaction, error)
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
	tenantID := c.GetString("tenant_id")

	from, to, label, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
		return
	}

	summary, err := h.repo.GetExpenseSummary(c.Request.Context(), tenantID, from, to, label)
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
	tenantID := c.GetString("tenant_id")

	from, to, label, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
		return
	}

	summary, err := h.repo.GetIncomeSummary(c.Request.Context(), tenantID, from, to, label)
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
	tenantID := c.GetString("tenant_id")

	from, to, _, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
		return
	}

	categories, err := h.repo.GetExpensesByCategory(c.Request.Context(), tenantID, from, to)
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
	tenantID := c.GetString("tenant_id")

	months, err := strconv.Atoi(c.DefaultQuery("months", "12"))
	if err != nil || months < 1 || months > 60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid months parameter, must be between 1 and 60"})
		return
	}

	expenses, err := h.repo.GetMonthlyExpenses(c.Request.Context(), tenantID, months)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get monthly expenses")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get monthly trends"})
		return
	}

	incomes, err := h.repo.GetMonthlyIncomes(c.Request.Context(), tenantID, months)
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

// productiveKeywords identifies category names that represent investment, savings, or
// capital allocation — these are positive financial behaviours and should not penalise
// the consumption ratio used in the health score.
var productiveKeywords = []string{
	"invers", "ahorro", "seguro", "educac", "retiro", "pension", "fondo",
	"activo", "propiedad", "inmueble", "capital", "emerg", "patrimonio",
	"cripto", "bitcoin", "etf", "accion", "bono", "plazo fijo",
}

// isProductiveCategory returns true when the category name contains a keyword that
// indicates investment, savings, insurance, or asset accumulation.
func isProductiveCategory(name string) bool {
	lower := strings.ToLower(name)
	for _, kw := range productiveKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// parseIntQuery reads an optional integer query parameter. Returns 0 when absent or invalid.
func parseIntQuery(c *gin.Context, key string) int {
	val := c.Query(key)
	if val == "" {
		return 0
	}
	n, err := strconv.Atoi(val)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

// clamp100 constrains a float64 to [0, 100].
func clamp100(v float64) float64 {
	return math.Max(0, math.Min(100, v))
}

// GetFinancialHealth handles GET /insights/financial-health
//
// Computes a multi-dimensional financial health score (0-1000):
//   - Cash-flow dimension   (40%): consumption vs income ratio (continuous, not buckets)
//   - Planning dimension    (30%): budgets, savings goals, recurring setups
//   - Consistency dimension (20%): streak + days active
//   - Engagement dimension  (10%): AI usage + analytics views
//
// Optional behavioral query params (sourced from the BehaviorProfile endpoint):
// streak, days_active, budgets_created, budget_compliance, savings_goals,
// savings_deposits, recurring_setups, ai_applied.
func (h *AnalyticsHandler) GetFinancialHealth(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	ctx := c.Request.Context()

	summary, err := h.repo.GetDashboardSummary(ctx, tenantID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get financial health")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get financial health"})
		return
	}

	// --- Split expenses into consumption vs productive ---
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	categories, _ := h.repo.GetExpensesByCategory(ctx, tenantID, monthStart, now)

	productiveExpenses := 0.0
	for _, cat := range categories {
		if isProductiveCategory(cat.CategoryName) {
			productiveExpenses += cat.Amount
		}
	}
	consumptionExpenses := summary.CurrentMonthExpenses - productiveExpenses
	if consumptionExpenses < 0 {
		consumptionExpenses = 0
	}

	income := summary.CurrentMonthIncomes

	// --- Read optional behavioral params ---
	streak          := parseIntQuery(c, "streak")
	daysActive      := parseIntQuery(c, "days_active")
	budgetsCreated  := parseIntQuery(c, "budgets_created")
	compliance      := parseIntQuery(c, "budget_compliance")
	savingsGoals    := parseIntQuery(c, "savings_goals")
	savingsDeposits := parseIntQuery(c, "savings_deposits")
	recurringSetups := parseIntQuery(c, "recurring_setups")
	aiApplied       := parseIntQuery(c, "ai_applied")

	// ================================================================
	// Dimension 1 — Cash Flow (0-100)
	// Continuous formula instead of 5 hard buckets.
	// ================================================================
	cashFlowScore := 100.0
	if income > 0 {
		ratio := consumptionExpenses / income
		cashFlowScore = clamp100(100 - ratio*100)
		// Bonus for high productive savings rate.
		trueSavingsRate := (summary.CurrentMonthBalance + productiveExpenses) / income * 100
		if trueSavingsRate > 20 {
			cashFlowScore = clamp100(cashFlowScore + 10)
		}
	} else if summary.CurrentMonthExpenses > 0 {
		cashFlowScore = 5 // no income but has expenses = critical
	}

	// ================================================================
	// Dimension 2 — Planning (0-100)
	// Rewards intentional financial planning actions.
	// ================================================================
	planningScore := 0.0
	if budgetsCreated > 0 {
		planningScore += 25
	}
	planningScore += clamp100(float64(compliance) * 10)
	if savingsGoals > 0 {
		planningScore += 20
	}
	planningScore += clamp100(float64(savingsDeposits) * 5)
	planningScore = clamp100(planningScore)

	// ================================================================
	// Dimension 3 — Consistency (0-100)
	// Rewards regular usage and long tenure.
	// ================================================================
	streakFactor  := math.Min(float64(streak)/30.0, 1.0) * 60
	tenureFactor  := math.Min(float64(daysActive)/90.0, 1.0) * 40
	consistencyScore := streakFactor + tenureFactor

	// ================================================================
	// Dimension 4 — Engagement (0-100)
	// Rewards use of AI and analytical tools.
	// ================================================================
	engagementScore := clamp100(float64(aiApplied)*20 + float64(recurringSetups)*15)

	// ================================================================
	// Weighted composite score (0-100) → display score (0-1000)
	// ================================================================
	finalScore := cashFlowScore*0.40 + planningScore*0.30 + consistencyScore*0.20 + engagementScore*0.10
	displayScore := finalScore * 10 // 0-1000 scale

	status := "excellent"
	switch {
	case finalScore < 30:
		status = "critical"
	case finalScore < 50:
		status = "poor"
	case finalScore < 70:
		status = "fair"
	case finalScore < 85:
		status = "good"
	}

	// True savings rate for informational display.
	trueSavingsRate := 0.0
	if income > 0 {
		trueSavings := summary.CurrentMonthBalance + productiveExpenses
		trueSavingsRate = trueSavings / income * 100
		if trueSavingsRate < 0 {
			trueSavingsRate = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"score":                  displayScore,
		"status":                 status,
		"savings_rate":           trueSavingsRate,
		"current_month_expenses": summary.CurrentMonthExpenses,
		"current_month_incomes":  summary.CurrentMonthIncomes,
		"current_month_balance":  summary.CurrentMonthBalance,
		"productive_expenses":    productiveExpenses,
		"consumption_expenses":   consumptionExpenses,
		// Dimension breakdown for frontend display.
		"cash_flow_score":    math.Round(cashFlowScore),
		"planning_score":     math.Round(planningScore),
		"consistency_score":  math.Round(consistencyScore),
		"engagement_score":   math.Round(engagementScore),
	})
}

// GetReport handles GET /reports
// Query params: start_date, end_date (YYYY-MM-DD).
// Returns total_income, total_expenses, transactions array, and category_summary.
func (h *AnalyticsHandler) GetReport(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

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

	expenseSummary, err := h.repo.GetExpenseSummary(ctx, tenantID, from, to, "custom")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get expense summary for report")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	incomeSummary, err := h.repo.GetIncomeSummary(ctx, tenantID, from, to, "custom")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get income summary for report")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	transactions, err := h.repo.GetTransactionsForReport(ctx, tenantID, from, to)
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
	tenantID := c.GetString("tenant_id")

	summary, err := h.repo.GetDashboardSummary(c.Request.Context(), tenantID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get dashboard summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get dashboard"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
