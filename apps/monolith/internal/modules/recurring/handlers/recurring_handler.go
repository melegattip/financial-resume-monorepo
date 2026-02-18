package handlers

import (
	"fmt"
	"net/http"
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

// ExecuteResponse is the response body returned after a manual execution
type ExecuteResponse struct {
	RecurringID   string `json:"recurring_id"`
	TransactionID string `json:"transaction_id"`
	Type          string `json:"type"`
	NextDate      string `json:"next_date"`
	ExecutionCount int   `json:"execution_count"`
	IsActive      bool   `json:"is_active"`
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

// Create handles POST /recurring
func (h *RecurringHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	activeOnly := c.Query("active_only") == "true"

	var items []*domain.RecurringTransaction
	var err error

	if activeOnly {
		items, err = h.repo.ListActive(c.Request.Context(), userID.(string))
	} else {
		items, err = h.repo.List(c.Request.Context(), userID.(string))
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
		"data":  response,
		"total": len(response),
	})
}

// ListDue handles GET /recurring/due
// Returns all active recurring transactions that are due for execution (NextDate <= now)
func (h *RecurringHandler) ListDue(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), userID.(string), id)
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), userID.(string), id)
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), userID.(string), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get recurring transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recurring transaction"})
		return
	}

	if rt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring transaction not found"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID.(string), id); err != nil {
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), userID.(string), id)
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), userID.(string), id)
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	rt, err := h.repo.GetByID(c.Request.Context(), userID.(string), id)
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

	now := time.Now().UTC()
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

		// Advance the recurring transaction state
		rt.Execute()

		rt.UpdatedAt = time.Now().UTC()
		return tx.Table("recurring_transactions").
			Where("id = ? AND user_id = ? AND deleted_at IS NULL", rt.ID, rt.UserID).
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
