package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/ports"
)

// BudgetHandler handles HTTP requests for budget operations.
type BudgetHandler struct {
	repo   ports.BudgetRepository
	logger zerolog.Logger
}

// NewBudgetHandler creates a new BudgetHandler.
func NewBudgetHandler(repo ports.BudgetRepository, logger zerolog.Logger) *BudgetHandler {
	return &BudgetHandler{
		repo:   repo,
		logger: logger,
	}
}

// --- Request / Response types ---

// CreateBudgetRequest is the request body for creating a budget.
type CreateBudgetRequest struct {
	CategoryID string              `json:"category_id" binding:"required"`
	Amount     float64             `json:"amount" binding:"required,gt=0"`
	Period     domain.BudgetPeriod `json:"period" binding:"required"`
	AlertAt    float64             `json:"alert_at"` // optional; defaults to 0.80
}

// UpdateBudgetRequest is the request body for updating a budget.
type UpdateBudgetRequest struct {
	Amount   *float64 `json:"amount,omitempty"`
	AlertAt  *float64 `json:"alert_at,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// BudgetResponse is the HTTP response format for a single budget.
type BudgetResponse struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	CategoryID       string  `json:"category_id"`
	Amount           float64 `json:"amount"`
	SpentAmount      float64 `json:"spent_amount"`
	RemainingAmount  float64 `json:"remaining_amount"`
	SpentPercentage  float64 `json:"spent_percentage"`
	Period           string  `json:"period"`
	PeriodStart      string  `json:"period_start"`
	PeriodEnd        string  `json:"period_end"`
	AlertAt          float64 `json:"alert_at"`
	AlertTriggered   bool    `json:"alert_triggered"`
	Status           string  `json:"status"`
	IsActive         bool    `json:"is_active"`
	IsInCurrentPeriod bool   `json:"is_in_current_period"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// BudgetStatusResponse is the response for GET /budgets/status.
type BudgetStatusResponse struct {
	Summary domain.BudgetSummary `json:"summary"`
	Budgets []BudgetResponse     `json:"budgets"`
}

func toBudgetResponse(b *domain.Budget) BudgetResponse {
	return BudgetResponse{
		ID:                b.ID,
		UserID:            b.UserID,
		CategoryID:        b.CategoryID,
		Amount:            b.Amount,
		SpentAmount:       b.SpentAmount,
		RemainingAmount:   b.GetRemainingAmount(),
		SpentPercentage:   b.GetSpentPercentage(),
		Period:            string(b.Period),
		PeriodStart:       b.PeriodStart.Format(time.RFC3339),
		PeriodEnd:         b.PeriodEnd.Format(time.RFC3339),
		AlertAt:           b.AlertAt,
		AlertTriggered:    b.IsAlertTriggered(),
		Status:            string(b.Status),
		IsActive:          b.IsActive,
		IsInCurrentPeriod: b.IsInCurrentPeriod(),
		CreatedAt:         b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         b.UpdatedAt.Format(time.RFC3339),
	}
}

func buildSummary(budgets []*domain.Budget) domain.BudgetSummary {
	summary := domain.BudgetSummary{
		TotalBudgets: len(budgets),
	}

	var totalUsage float64
	for _, b := range budgets {
		summary.TotalAllocated += b.Amount
		summary.TotalSpent += b.SpentAmount
		totalUsage += b.GetSpentPercentage()

		switch b.Status {
		case domain.BudgetStatusOnTrack:
			summary.OnTrackCount++
		case domain.BudgetStatusWarning:
			summary.WarningCount++
		case domain.BudgetStatusExceeded:
			summary.ExceededCount++
		}
	}

	if len(budgets) > 0 {
		summary.AverageUsage = totalUsage / float64(len(budgets))
	}

	return summary
}

// --- Handlers ---

// Create handles POST /budgets
func (h *BudgetHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := c.GetString("tenant_id")

	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate amount ceiling
	if req.Amount > 1e12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "budget amount is too large"})
		return
	}

	// Default alert_at to 0.80 if not provided or zero
	alertAt := req.AlertAt
	if alertAt == 0 {
		alertAt = 0.80
	}
	if alertAt < 0 || alertAt > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alert_at must be between 0 and 1"})
		return
	}

	// Build the domain entity
	budget, err := domain.NewBudgetBuilder().
		SetUserID(userID.(string)).
		SetCategoryID(req.CategoryID).
		SetAmount(req.Amount).
		SetPeriod(req.Period).
		SetAlertAt(alertAt).
		Build()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budget.TenantID = tenantID

	// Calculate current spent amount from existing expenses for the period
	spentAmount, err := h.repo.GetExpensesForPeriod(
		c.Request.Context(),
		tenantID,
		req.CategoryID,
		budget.PeriodStart,
		budget.PeriodEnd,
	)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get expenses for period")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate current spending"})
		return
	}
	budget.UpdateSpentAmount(spentAmount)

	// Persist
	if err := h.repo.Create(c.Request.Context(), budget); err != nil {
		h.logger.Error().Err(err).Msg("failed to create budget")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create budget"})
		return
	}

	c.JSON(http.StatusCreated, toBudgetResponse(budget))
}

// List handles GET /budgets
// Supported query params: period, category_id, status, active_only
func (h *BudgetHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	activeOnlyStr := c.Query("active_only")
	activeOnly := false
	if activeOnlyStr != "" {
		var err error
		activeOnly, err = strconv.ParseBool(activeOnlyStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid active_only parameter"})
			return
		}
	}

	var budgets []*domain.Budget
	var err error

	if activeOnly {
		budgets, err = h.repo.ListActive(c.Request.Context(), tenantID)
	} else {
		budgets, err = h.repo.List(c.Request.Context(), tenantID)
	}

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list budgets")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list budgets"})
		return
	}

	// Optional client-side filters
	period := domain.BudgetPeriod(c.Query("period"))
	categoryID := c.Query("category_id")
	status := domain.BudgetStatus(c.Query("status"))

	filtered := make([]*domain.Budget, 0, len(budgets))
	for _, b := range budgets {
		if period != "" && b.Period != period {
			continue
		}
		if categoryID != "" && b.CategoryID != categoryID {
			continue
		}
		if status != "" && b.Status != status {
			continue
		}
		filtered = append(filtered, b)
	}

	response := make([]BudgetResponse, len(filtered))
	for i, b := range filtered {
		response[i] = toBudgetResponse(b)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"total": len(response),
	})
}

// GetStatus handles GET /budgets/status
func (h *BudgetHandler) GetStatus(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	budgets, err := h.repo.List(c.Request.Context(), tenantID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get budget status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get budget status"})
		return
	}

	budgetResponses := make([]BudgetResponse, len(budgets))
	for i, b := range budgets {
		budgetResponses[i] = toBudgetResponse(b)
	}

	c.JSON(http.StatusOK, BudgetStatusResponse{
		Summary: buildSummary(budgets),
		Budgets: budgetResponses,
	})
}

// GetDashboard handles GET /budgets/dashboard — returns the same summary as GetStatus.
func (h *BudgetHandler) GetDashboard(c *gin.Context) {
	h.GetStatus(c)
}

// GetByID handles GET /budgets/:id
func (h *BudgetHandler) GetByID(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	budget, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get budget")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get budget"})
		return
	}

	if budget == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "budget not found"})
		return
	}

	c.JSON(http.StatusOK, toBudgetResponse(budget))
}

// Update handles PUT /budgets/:id
func (h *BudgetHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	budget, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get budget")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get budget"})
		return
	}

	if budget == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "budget not found"})
		return
	}

	var req UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and apply updates
	if req.Amount != nil {
		if *req.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
			return
		}
		if *req.Amount > 1e12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "budget amount is too large"})
			return
		}
		budget.Amount = *req.Amount
	}

	if req.AlertAt != nil {
		if *req.AlertAt < 0 || *req.AlertAt > 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "alert_at must be between 0 and 1"})
			return
		}
		budget.AlertAt = *req.AlertAt
	}

	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	}

	// Recalculate status after potential amount/alertAt change
	budget.UpdateSpentAmount(budget.SpentAmount)
	budget.UpdatedAt = time.Now()

	if err := h.repo.Update(c.Request.Context(), budget); err != nil {
		h.logger.Error().Err(err).Msg("failed to update budget")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update budget"})
		return
	}

	c.JSON(http.StatusOK, toBudgetResponse(budget))
}

// Delete handles DELETE /budgets/:id
func (h *BudgetHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	budget, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get budget")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get budget"})
		return
	}

	if budget == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "budget not found"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), tenantID, id); err != nil {
		h.logger.Error().Err(err).Msg("failed to delete budget")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete budget"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
