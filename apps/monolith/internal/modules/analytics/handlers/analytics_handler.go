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
