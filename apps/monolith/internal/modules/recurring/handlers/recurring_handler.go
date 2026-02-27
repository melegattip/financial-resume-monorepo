package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/ports"
)

// RecurringHandler handles HTTP requests for recurring transactions
type RecurringHandler struct {
	repo     ports.RecurringTransactionRepository
	db       *gorm.DB
	eventBus ports.EventBus
	logger   zerolog.Logger
}

// NewRecurringHandler creates a new recurring transaction handler
func NewRecurringHandler(
	repo ports.RecurringTransactionRepository,
	db *gorm.DB,
	eventBus ports.EventBus,
	logger zerolog.Logger,
) *RecurringHandler {
	return &RecurringHandler{
		repo:     repo,
		db:       db,
		eventBus: eventBus,
		logger:   logger,
	}
}

// --- Request / Response types ---

// CreateRecurringRequest is the request body for creating a recurring transaction
type CreateRecurringRequest struct {
	Amount        float64    `json:"amount" binding:"required,gt=0"`
	Description   string     `json:"description" binding:"required"`
	CategoryID    string     `json:"category_id"`
	Type          string     `json:"type" binding:"required"`
	Frequency     string     `json:"frequency" binding:"required"`
	NextDate      string     `json:"next_date" binding:"required"`
	AutoCreate    *bool      `json:"auto_create"`
	NotifyBefore  *int       `json:"notify_before"`
	EndDate       *string    `json:"end_date"`
	MaxExecutions *int       `json:"max_executions"`
}

// UpdateRecurringRequest is the request body for updating a recurring transaction
type UpdateRecurringRequest struct {
	Amount        float64  `json:"amount" binding:"required,gt=0"`
	Description   string   `json:"description" binding:"required"`
	CategoryID    string   `json:"category_id"`
	Type          string   `json:"type" binding:"required"`
	Frequency     string   `json:"frequency" binding:"required"`
	NextDate      string   `json:"next_date" binding:"required"`
	AutoCreate    *bool    `json:"auto_create"`
	NotifyBefore  *int     `json:"notify_before"`
	EndDate       *string  `json:"end_date"`
	MaxExecutions *int     `json:"max_executions"`
}

// RecurringResponse is the response format for a recurring transaction
type RecurringResponse struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description"`
	CategoryID     *string `json:"category_id"`
	Type           string  `json:"type"`
	Frequency      string  `json:"frequency"`
	NextDate       string  `json:"next_date"`
	LastExecuted   *string `json:"last_executed"`
	IsActive       bool    `json:"is_active"`
	AutoCreate     bool    `json:"auto_create"`
	NotifyBefore   int     `json:"notify_before"`
	EndDate        *string `json:"end_date"`
	ExecutionCount int     `json:"execution_count"`
	MaxExecutions  *int    `json:"max_executions"`
	DaysUntilNext  int     `json:"days_until_next"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ManualExecuteRequest allows the caller to specify a target execution date.
// If omitted, the execution date defaults to now.
type ManualExecuteRequest struct {
	ExecutionDate string `json:"execution_date"` // RFC3339; optional
}

// ExecuteResponse is the response body returned after a manual execution
type ExecuteResponse struct {
	RecurringID   string `json:"recurring_id"`
	TransactionID string `json:"transaction_id"`
	Type          string `json:"type"`
	NextDate      string `json:"next_date"`
	ExecutionCount int   `json:"execution_count"`
	IsActive      bool   `json:"is_active"`
}

// monthlyEquivalent converts a recurring transaction amount to its monthly equivalent.
func monthlyEquivalent(rt *domain.RecurringTransaction) float64 {
	switch rt.Frequency {
	case "daily":
		return rt.Amount * 30
	case "weekly":
		return rt.Amount * 4.33
	case "monthly":
		return rt.Amount
	case "yearly":
		return rt.Amount / 12
	default:
		return rt.Amount
	}
}

func toRecurringResponse(rt *domain.RecurringTransaction) RecurringResponse {
	resp := RecurringResponse{
		ID:             rt.ID,
		UserID:         rt.UserID,
		Amount:         rt.Amount,
		Description:    rt.Description,
		CategoryID:     rt.CategoryID,
		Type:           rt.Type,
		Frequency:      rt.Frequency,
		NextDate:       rt.NextDate.Format(time.RFC3339),
		IsActive:       rt.IsActive,
		AutoCreate:     rt.AutoCreate,
		NotifyBefore:   rt.NotifyBefore,
		ExecutionCount: rt.ExecutionCount,
		MaxExecutions:  rt.MaxExecutions,
		DaysUntilNext:  rt.GetDaysUntilNext(),
		CreatedAt:      rt.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      rt.UpdatedAt.Format(time.RFC3339),
	}

	if rt.LastExecuted != nil {
		s := rt.LastExecuted.Format(time.RFC3339)
		resp.LastExecuted = &s
	}

	if rt.EndDate != nil {
		s := rt.EndDate.Format(time.RFC3339)
		resp.EndDate = &s
	}

	return resp
}

// expenseModel is a minimal struct used to insert a row into the expenses table
// directly via GORM without importing the transactions module
type expenseModel struct {
	ID              string     `gorm:"column:id;primaryKey"`
	UserID          string     `gorm:"column:user_id"`
	TenantID        string     `gorm:"column:tenant_id"`
	CategoryID      string     `gorm:"column:category_id"`
	Amount          float64    `gorm:"column:amount"`
	Description     string     `gorm:"column:description"`
	TransactionDate time.Time  `gorm:"column:transaction_date"`
	PaymentMethod   string     `gorm:"column:payment_method"`
	Notes           string     `gorm:"column:notes"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
}

func (expenseModel) TableName() string { return "expenses" }

// incomeModel is a minimal struct used to insert a row into the incomes table
// directly via GORM without importing the transactions module
type incomeModel struct {
	ID           string     `gorm:"column:id;primaryKey"`
	UserID       string     `gorm:"column:user_id"`
	TenantID     string     `gorm:"column:tenant_id"`
	Amount       float64    `gorm:"column:amount"`
	Source       string     `gorm:"column:source"`
	Description  string     `gorm:"column:description"`
	ReceivedDate time.Time  `gorm:"column:received_date"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (incomeModel) TableName() string { return "incomes" }

// --- Handlers ---

// GetDashboard handles GET /recurring-transactions/dashboard
// Returns a summary of active recurring transactions for the current tenant.
func (h *RecurringHandler) GetDashboard(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	items, err := h.repo.ListActive(c.Request.Context(), tenantID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get recurring transactions dashboard")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transactions dashboard"})
		return
	}

	now := time.Now().UTC()
	var monthlyIncomeTotal, monthlyExpenseTotal float64
	var totalInactive int
	upcoming := make([]RecurringResponse, 0)

	// Count inactive recurring transactions too
	allItems, err := h.repo.List(c.Request.Context(), tenantID)
	if err != nil {
		allItems = items
	}
	totalInactive = len(allItems) - len(items)

	for _, rt := range items {
		monthly := monthlyEquivalent(rt)
		if rt.Type == "income" {
			monthlyIncomeTotal += monthly
		} else {
			monthlyExpenseTotal += monthly
		}
		if !rt.NextDate.After(now.AddDate(0, 0, 30)) {
			upcoming = append(upcoming, toRecurringResponse(rt))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary": gin.H{
				"total_active":          len(items),
				"total_inactive":        totalInactive,
				"monthly_income_total":  monthlyIncomeTotal,
				"monthly_expense_total": monthlyExpenseTotal,
			},
			"upcoming": upcoming,
		},
	})
}

// Create handles POST /recurring
func (h *RecurringHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := c.GetString("tenant_id")

	var req CreateRecurringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nextDate, err := time.Parse(time.RFC3339, req.NextDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid next_date format, expected RFC3339"})
		return
	}

	builder := domain.NewRecurringTransactionBuilder().
		SetUserID(userID.(string)).
		SetAmount(req.Amount).
		SetDescription(req.Description).
		SetCategoryID(req.CategoryID).
		SetType(req.Type).
		SetFrequency(req.Frequency).
		SetNextDate(nextDate)

	if req.AutoCreate != nil {
		builder.SetAutoCreate(*req.AutoCreate)
	}
	if req.NotifyBefore != nil {
		builder.SetNotifyBefore(*req.NotifyBefore)
	}
	if req.EndDate != nil {
		endDate, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected RFC3339"})
			return
		}
		builder.SetEndDate(&endDate)
	}
	if req.MaxExecutions != nil {
		builder.SetMaxExecutions(req.MaxExecutions)
	}

	rt, err := builder.Build()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rt.TenantID = tenantID

	if err := h.repo.Create(c.Request.Context(), rt); err != nil {
		h.logger.Error().Err(err).Msg("failed to create recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create recurring transaction"})
		return
	}

	event := domain.RecurringTransactionCreatedEvent{
		RecurringID: rt.ID,
		User:        rt.UserID,
		Type:        rt.Type,
		Amount:      rt.Amount,
		Frequency:   rt.Frequency,
		NextDate:    rt.NextDate,
		Timestamp:   time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish RecurringTransactionCreatedEvent")
	}

	c.JSON(http.StatusCreated, toRecurringResponse(rt))
}

// List handles GET /recurring
// Query params: type (income|expense), frequency (daily|weekly|monthly|yearly), active_only (true|false)
func (h *RecurringHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	activeOnly := c.Query("active_only") == "true"

	var items []*domain.RecurringTransaction
	var err error

	if activeOnly {
		items, err = h.repo.ListActive(c.Request.Context(), tenantID)
	} else {
		items, err = h.repo.List(c.Request.Context(), tenantID)
	}

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list recurring transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list recurring transactions"})
		return
	}

	// Optional client-side filters (type, frequency)
	typeFilter := c.Query("type")
	freqFilter := c.Query("frequency")

	filtered := make([]*domain.RecurringTransaction, 0, len(items))
	for _, rt := range items {
		if typeFilter != "" && rt.Type != typeFilter {
			continue
		}
		if freqFilter != "" && rt.Frequency != freqFilter {
			continue
		}
		filtered = append(filtered, rt)
	}

	response := make([]RecurringResponse, len(filtered))
	for i, rt := range filtered {
		response[i] = toRecurringResponse(rt)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"transactions": response,
		},
		"total": len(response),
	})
}

// ListDue handles GET /recurring/due
// Returns all active recurring transactions that are due for execution (NextDate <= now)
func (h *RecurringHandler) ListDue(c *gin.Context) {

	items, err := h.repo.ListDue(c.Request.Context(), time.Now().UTC())
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list due recurring transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list due recurring transactions"})
		return
	}

	response := make([]RecurringResponse, len(items))
	for i, rt := range items {
		response[i] = toRecurringResponse(rt)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"total": len(response),
	})
}

// GetByID handles GET /recurring/:id
func (h *RecurringHandler) GetByID(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transaction"})
		return
	}

	if rt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring transaction not found"})
		return
	}

	c.JSON(http.StatusOK, toRecurringResponse(rt))
}

// Update handles PUT /recurring/:id
func (h *RecurringHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transaction"})
		return
	}

	if rt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring transaction not found"})
		return
	}

	var req UpdateRecurringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nextDate, err := time.Parse(time.RFC3339, req.NextDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid next_date format, expected RFC3339"})
		return
	}

	// Apply updates to the domain entity
	rt.Amount = req.Amount
	rt.Description = req.Description
	rt.Type = req.Type
	rt.Frequency = req.Frequency
	rt.NextDate = nextDate
	rt.UpdatedAt = time.Now().UTC()

	if req.CategoryID != "" {
		rt.CategoryID = &req.CategoryID
	} else {
		rt.CategoryID = nil
	}

	if req.AutoCreate != nil {
		rt.AutoCreate = *req.AutoCreate
	}
	if req.NotifyBefore != nil {
		rt.NotifyBefore = *req.NotifyBefore
	}

	if req.EndDate != nil {
		endDate, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected RFC3339"})
			return
		}
		rt.EndDate = &endDate
	} else {
		rt.EndDate = nil
	}

	rt.MaxExecutions = req.MaxExecutions

	if err := rt.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), rt); err != nil {
		h.logger.Error().Err(err).Msg("failed to update recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update recurring transaction"})
		return
	}

	event := domain.RecurringTransactionUpdatedEvent{
		RecurringID: rt.ID,
		User:        rt.UserID,
		Timestamp:   time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish RecurringTransactionUpdatedEvent")
	}

	c.JSON(http.StatusOK, toRecurringResponse(rt))
}

// Delete handles DELETE /recurring/:id
func (h *RecurringHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transaction"})
		return
	}

	if rt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring transaction not found"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), tenantID, id); err != nil {
		h.logger.Error().Err(err).Msg("failed to delete recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete recurring transaction"})
		return
	}

	event := domain.RecurringTransactionDeletedEvent{
		RecurringID: rt.ID,
		User:        rt.UserID,
		Timestamp:   time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish RecurringTransactionDeletedEvent")
	}

	c.JSON(http.StatusNoContent, nil)
}

// Pause handles POST /recurring/:id/pause
func (h *RecurringHandler) Pause(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transaction"})
		return
	}

	if rt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring transaction not found"})
		return
	}

	if !rt.IsActive {
		c.JSON(http.StatusConflict, gin.H{"error": "recurring transaction is already paused"})
		return
	}

	rt.Pause()

	if err := h.repo.Update(c.Request.Context(), rt); err != nil {
		h.logger.Error().Err(err).Msg("failed to pause recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to pause recurring transaction"})
		return
	}

	event := domain.RecurringTransactionPausedEvent{
		RecurringID: rt.ID,
		User:        rt.UserID,
		Timestamp:   time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish RecurringTransactionPausedEvent")
	}

	c.JSON(http.StatusOK, toRecurringResponse(rt))
}

// Resume handles POST /recurring/:id/resume
func (h *RecurringHandler) Resume(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transaction"})
		return
	}

	if rt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring transaction not found"})
		return
	}

	if rt.IsActive {
		c.JSON(http.StatusConflict, gin.H{"error": "recurring transaction is already active"})
		return
	}

	rt.Resume()

	if err := h.repo.Update(c.Request.Context(), rt); err != nil {
		h.logger.Error().Err(err).Msg("failed to resume recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resume recurring transaction"})
		return
	}

	event := domain.RecurringTransactionResumedEvent{
		RecurringID: rt.ID,
		User:        rt.UserID,
		Timestamp:   time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish RecurringTransactionResumedEvent")
	}

	c.JSON(http.StatusOK, toRecurringResponse(rt))
}

// ManualExecute handles POST /recurring/:id/execute
// It creates the corresponding expense or income record and advances the recurring transaction
func (h *RecurringHandler) ManualExecute(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transaction"})
		return
	}

	if rt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring transaction not found"})
		return
	}

	if !rt.IsActive {
		c.JSON(http.StatusConflict, gin.H{"error": "cannot execute a paused recurring transaction"})
		return
	}

	// Parse optional execution_date from request body.
	// If not provided, fall back to now so existing callers are unaffected.
	var req ManualExecuteRequest
	_ = c.ShouldBindJSON(&req) // deliberately ignore binding errors — body is optional

	now := time.Now().UTC()
	transactionDate := now
	if req.ExecutionDate != "" {
		parsed, parseErr := time.Parse(time.RFC3339, req.ExecutionDate)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution_date format, expected RFC3339"})
			return
		}
		transactionDate = parsed
	}

	transactionID := uuid.New().String()

	// Execute atomically: create the expense/income AND update the recurring transaction
	// in a single database transaction to guarantee consistency.
	txErr := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// Create the actual expense or income record in the appropriate table
		switch rt.Type {
		case "expense":
			categoryID := ""
			if rt.CategoryID != nil {
				categoryID = *rt.CategoryID
			}
			record := &expenseModel{
				ID:              transactionID,
				UserID:          rt.UserID,
				TenantID:        rt.TenantID,
				CategoryID:      categoryID,
				Amount:          rt.Amount,
				Description:     rt.Description,
				TransactionDate: transactionDate,
				PaymentMethod:   "",
				Notes:           "Auto-created from recurring transaction: " + rt.ID,
				CreatedAt:       now,
				UpdatedAt:       now,
			}
			if err := tx.Create(record).Error; err != nil {
				return err
			}

		case "income":
			record := &incomeModel{
				ID:           transactionID,
				UserID:       rt.UserID,
				TenantID:     rt.TenantID,
				Amount:       rt.Amount,
				Source:       "recurring",
				Description:  rt.Description,
				ReceivedDate: transactionDate,
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if err := tx.Create(record).Error; err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown recurring transaction type: %s", rt.Type)
		}

		// When the caller specified a target date, advance NextDate from that date
		// so the next occurrence is calculated relative to the chosen period.
		if req.ExecutionDate != "" {
			rt.NextDate = transactionDate
		}

		// Advance the recurring transaction state
		rt.Execute()

		rt.UpdatedAt = time.Now().UTC()
		return tx.Table("recurring_transactions").
			Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", rt.ID, rt.TenantID).
			Updates(map[string]interface{}{
				"next_date":       rt.NextDate,
				"last_executed":   rt.LastExecuted,
				"execution_count": rt.ExecutionCount,
				"is_active":       rt.IsActive,
				"updated_at":      rt.UpdatedAt,
			}).Error
	})

	if txErr != nil {
		// Check if it was an unknown type error (should be 400, not 500)
		if txErr.Error() == "unknown recurring transaction type: "+rt.Type {
			c.JSON(http.StatusBadRequest, gin.H{"error": txErr.Error()})
			return
		}
		h.logger.Error().Err(txErr).Msg("failed to execute recurring transaction atomically")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to execute recurring transaction"})
		return
	}

	event := domain.RecurringTransactionExecutedEvent{
		RecurringID:   rt.ID,
		User:          rt.UserID,
		Type:          rt.Type,
		Amount:        rt.Amount,
		TransactionID: transactionID,
		Timestamp:     now,
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish RecurringTransactionExecutedEvent")
	}

	c.JSON(http.StatusOK, ExecuteResponse{
		RecurringID:    rt.ID,
		TransactionID:  transactionID,
		Type:           rt.Type,
		NextDate:       rt.NextDate.Format(time.RFC3339),
		ExecutionCount: rt.ExecutionCount,
		IsActive:       rt.IsActive,
	})
}

// GetProjection handles GET /recurring-transactions/projection?months=N
// Returns a cash flow projection for the next N months based on active recurring transactions.
func (h *RecurringHandler) GetProjection(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	months := 6
	if m := c.Query("months"); m != "" {
		if parsed, err := strconv.Atoi(m); err == nil && parsed > 0 && parsed <= 24 {
			months = parsed
		}
	}

	items, err := h.repo.ListActive(c.Request.Context(), tenantID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list active recurring transactions for projection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute projection"})
		return
	}

	var totalMonthlyIncome, totalMonthlyExpenses float64
	for _, rt := range items {
		m := monthlyEquivalent(rt)
		if rt.Type == "income" {
			totalMonthlyIncome += m
		} else {
			totalMonthlyExpenses += m
		}
	}

	now := time.Now().UTC()
	monthNames := []string{
		"enero", "febrero", "marzo", "abril", "mayo", "junio",
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
	}

	monthlyProjections := make([]gin.H, months)
	var cumulativeNet float64
	for i := 0; i < months; i++ {
		t := now.AddDate(0, i, 0)
		netAmount := totalMonthlyIncome - totalMonthlyExpenses
		cumulativeNet += netAmount
		monthlyProjections[i] = gin.H{
			"month":          t.Format("2006-01"),
			"month_display":  monthNames[t.Month()-1] + " " + strconv.Itoa(t.Year()),
			"income":         totalMonthlyIncome,
			"expenses":       totalMonthlyExpenses,
			"net_amount":     netAmount,
			"cumulative_net": cumulativeNet,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary": gin.H{
				"average_monthly_income":   totalMonthlyIncome,
				"average_monthly_expenses": totalMonthlyExpenses,
				"net_projected_amount":     totalMonthlyIncome - totalMonthlyExpenses,
				"months":                  months,
			},
			"monthly_projections": monthlyProjections,
		},
	})
}

// executeRecurringTransaction creates the expense/income record and advances the recurring transaction.
// Returns the created transaction ID and any error.
func (h *RecurringHandler) executeRecurringTransaction(ctx context.Context, rt *domain.RecurringTransaction) (string, error) {
	now := time.Now().UTC()
	transactionID := uuid.New().String()

	txErr := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch rt.Type {
		case "expense":
			categoryID := ""
			if rt.CategoryID != nil {
				categoryID = *rt.CategoryID
			}
			record := &expenseModel{
				ID:              transactionID,
				UserID:          rt.UserID,
				TenantID:        rt.TenantID,
				CategoryID:      categoryID,
				Amount:          rt.Amount,
				Description:     rt.Description,
				TransactionDate: now,
				PaymentMethod:   "",
				Notes:           "Auto-created from recurring transaction: " + rt.ID,
				CreatedAt:       now,
				UpdatedAt:       now,
			}
			if err := tx.Create(record).Error; err != nil {
				return err
			}
		case "income":
			record := &incomeModel{
				ID:           transactionID,
				UserID:       rt.UserID,
				TenantID:     rt.TenantID,
				Amount:       rt.Amount,
				Source:       "recurring",
				Description:  rt.Description,
				ReceivedDate: now,
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if err := tx.Create(record).Error; err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown recurring transaction type: %s", rt.Type)
		}

		rt.Execute()
		rt.UpdatedAt = time.Now().UTC()
		return tx.Table("recurring_transactions").
			Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", rt.ID, rt.TenantID).
			Updates(map[string]interface{}{
				"next_date":       rt.NextDate,
				"last_executed":   rt.LastExecuted,
				"execution_count": rt.ExecutionCount,
				"is_active":       rt.IsActive,
				"updated_at":      rt.UpdatedAt,
			}).Error
	})

	return transactionID, txErr
}

// ProcessPending handles POST /recurring-transactions/batch/process
// Processes all due recurring transactions for the authenticated tenant.
func (h *RecurringHandler) ProcessPending(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	dueItems, err := h.repo.ListDue(c.Request.Context(), time.Now().UTC())
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list due recurring transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list due transactions"})
		return
	}

	var successCount, failureCount int
	for _, rt := range dueItems {
		if rt.TenantID != tenantID || !rt.IsActive {
			continue
		}

		_, txErr := h.executeRecurringTransaction(c.Request.Context(), rt)
		if txErr != nil {
			h.logger.Error().Err(txErr).Str("id", rt.ID).Msg("failed to process pending recurring transaction")
			failureCount++
		} else {
			successCount++
			event := domain.RecurringTransactionExecutedEvent{
				RecurringID:   rt.ID,
				User:          rt.UserID,
				Type:          rt.Type,
				Amount:        rt.Amount,
				TransactionID: rt.ID,
				Timestamp:     time.Now().UTC(),
			}
			if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
				h.logger.Warn().Err(err).Msg("failed to publish RecurringTransactionExecutedEvent")
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"processed_count": successCount + failureCount,
		"success_count":   successCount,
		"failure_count":   failureCount,
	})
}

// SendNotifications handles POST /recurring-transactions/batch/notify
// Stub: notification service is not implemented yet.
func (h *RecurringHandler) SendNotifications(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sent_count": 0,
		"message":    "notification service not implemented",
	})
}
