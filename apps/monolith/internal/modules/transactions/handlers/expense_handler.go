package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/ports"
)

type ExpenseHandler struct {
	repo     ports.ExpenseRepository
	eventBus ports.EventBus
	logger   zerolog.Logger
}

func NewExpenseHandler(repo ports.ExpenseRepository, eventBus ports.EventBus, logger zerolog.Logger) *ExpenseHandler {
	return &ExpenseHandler{
		repo:     repo,
		eventBus: eventBus,
		logger:   logger,
	}
}

// CreateExpenseRequest is the request body for creating an expense.
// transaction_date is optional (defaults to now); due_date is accepted as an alias.
// category_id is optional (defaults to empty string when not provided).
type CreateExpenseRequest struct {
	CategoryID      string  `json:"category_id"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	Description     string  `json:"description" binding:"required"`
	TransactionDate string  `json:"transaction_date"`
	DueDate         string  `json:"due_date"`         // alias for transaction_date
	PaymentMethod   string  `json:"payment_method"`
	Notes           string  `json:"notes"`
	Paid            bool    `json:"paid"`
}

// UpdateExpenseRequest is the request body for updating an expense.
// All fields are optional; omitted fields keep their existing values.
// payment_amount: delta amount being paid now (used for partial payments).
// amount_paid: cumulative amount paid so far (used for total/overpayment payments).
type UpdateExpenseRequest struct {
	CategoryID      string   `json:"category_id"`
	Amount          *float64 `json:"amount"`
	Description     string   `json:"description"`
	TransactionDate string   `json:"transaction_date"`
	DueDate         string   `json:"due_date"` // alias for transaction_date
	PaymentMethod   string   `json:"payment_method"`
	Notes           string   `json:"notes"`
	Paid            *bool    `json:"paid"`
	AmountPaid      *float64 `json:"amount_paid"`
	PendingAmount   *float64 `json:"pending_amount"`
	PaymentAmount   *float64 `json:"payment_amount"` // delta: adds to existing amount_paid
}

// ExpenseResponse is the response format for an expense
type ExpenseResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	CategoryID      string  `json:"category_id"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
	PaymentMethod   string  `json:"payment_method"`
	Notes           string  `json:"notes"`
	Paid            bool    `json:"paid"`
	AmountPaid      float64 `json:"amount_paid"`
	PendingAmount   float64 `json:"pending_amount"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func toExpenseResponse(e *domain.Expense) ExpenseResponse {
	return ExpenseResponse{
		ID:              e.ID,
		UserID:          e.UserID,
		CategoryID:      e.CategoryID,
		Amount:          e.Amount,
		Description:     e.Description,
		TransactionDate: e.TransactionDate.Format(time.RFC3339),
		PaymentMethod:   e.PaymentMethod,
		Notes:           e.Notes,
		Paid:            e.Paid,
		AmountPaid:      e.AmountPaid,
		PendingAmount:   e.PendingAmount,
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       e.UpdatedAt.Format(time.RFC3339),
	}
}

// Create handles POST /api/v1/expenses
func (h *ExpenseHandler) Create(c *gin.Context) {
	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id and tenant_id from JWT (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := c.GetString("tenant_id")

	// Resolve transaction date: prefer transaction_date, fall back to due_date, then now.
	rawDate := req.TransactionDate
	if rawDate == "" {
		rawDate = req.DueDate
	}
	var transactionDate time.Time
	if rawDate == "" {
		transactionDate = time.Now().UTC()
	} else {
		var parseErr error
		transactionDate, parseErr = time.Parse(time.RFC3339, rawDate)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
			return
		}
	}

	// Create expense domain entity
	expense, err := domain.NewExpense(
		userID.(string),
		req.CategoryID,
		req.Amount,
		req.Description,
		transactionDate,
		req.PaymentMethod,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.Notes = req.Notes
	expense.TenantID = tenantID
	if req.Paid {
		expense.ApplyPayment(true, expense.Amount)
	}

	// Save to repository
	if err := h.repo.Create(c.Request.Context(), expense); err != nil {
		h.logger.Error().Err(err).Msg("failed to create expense")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create expense"})
		return
	}

	// Publish event
	event := domain.ExpenseCreatedEvent{
		ExpenseID:       expense.ID,
		User:            expense.UserID,
		CategoryID:      expense.CategoryID,
		Amount:          expense.Amount,
		Description:     expense.Description,
		TransactionDate: expense.TransactionDate,
		Timestamp:       time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish ExpenseCreatedEvent")
	}

	c.JSON(http.StatusCreated, toExpenseResponse(expense))
}

// List handles GET /api/v1/expenses
func (h *ExpenseHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	// Pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	expenses, err := h.repo.FindByTenantID(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list expenses")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list expenses"})
		return
	}

	response := make([]ExpenseResponse, len(expenses))
	for i, e := range expenses {
		response[i] = toExpenseResponse(e)
	}

	c.JSON(http.StatusOK, gin.H{
		"expenses": response,
		"total":    len(response),
		"limit":    limit,
		"offset":   offset,
	})
}

// GetByID handles GET /api/v1/expenses/:id
func (h *ExpenseHandler) GetByID(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	expense, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get expense")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expense"})
		return
	}

	if expense == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
		return
	}

	// Verify tenant membership
	if expense.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.JSON(http.StatusOK, toExpenseResponse(expense))
}

// Update handles PUT /api/v1/expenses/:id
func (h *ExpenseHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	expense, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get expense")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expense"})
		return
	}

	if expense == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
		return
	}

	// Verify tenant membership
	if expense.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- Handle paid-only update (e.g. toggle paid, partial payment) ---
	// If the request only contains payment fields (no amount/description), skip Update() and
	// go straight to ApplyPayment so we don't require amount/description in the payload.
	hasCoreFields := req.Amount != nil || req.Description != "" || req.CategoryID != "" || req.TransactionDate != "" || req.DueDate != ""
	hasPaymentFields := req.Paid != nil || req.AmountPaid != nil || req.PaymentAmount != nil

	if hasCoreFields {
		rawDate := req.TransactionDate
		if rawDate == "" {
			rawDate = req.DueDate
		}
		var transactionDate time.Time
		if rawDate == "" {
			transactionDate = expense.TransactionDate
		} else {
			var parseErr error
			transactionDate, parseErr = time.Parse(time.RFC3339, rawDate)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
				return
			}
		}

		categoryID := req.CategoryID
		if categoryID == "" {
			categoryID = expense.CategoryID
		}

		amount := expense.Amount
		if req.Amount != nil {
			amount = *req.Amount
		}
		description := expense.Description
		if req.Description != "" {
			description = req.Description
		}
		paymentMethod := expense.PaymentMethod
		if req.PaymentMethod != "" {
			paymentMethod = req.PaymentMethod
		}
		notes := expense.Notes
		if req.Notes != "" {
			notes = req.Notes
		}

		if err := expense.Update(categoryID, amount, description, transactionDate, paymentMethod, notes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// --- Apply payment state changes ---
	if hasPaymentFields {
		if req.Paid != nil && !*req.Paid {
			// Explicitly marking as unpaid
			expense.ApplyPayment(false, 0)
		} else if req.PaymentAmount != nil {
			// Partial payment delta: add to existing amount_paid
			newAmountPaid := expense.AmountPaid + *req.PaymentAmount
			expense.ApplyPayment(true, newAmountPaid)
		} else if req.AmountPaid != nil {
			// Absolute cumulative amount paid provided
			paid := true
			if req.Paid != nil {
				paid = *req.Paid
			}
			expense.ApplyPayment(paid, *req.AmountPaid)
		} else if req.Paid != nil && *req.Paid {
			// Mark fully paid without explicit amount — assume full amount
			expense.ApplyPayment(true, expense.Amount)
		}
	}

	if err := h.repo.Update(c.Request.Context(), expense); err != nil {
		h.logger.Error().Err(err).Msg("failed to update expense")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update expense"})
		return
	}

	// Publish event
	event := domain.ExpenseUpdatedEvent{
		ExpenseID: expense.ID,
		User:      expense.UserID,
		Timestamp: time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish ExpenseUpdatedEvent")
	}

	c.JSON(http.StatusOK, toExpenseResponse(expense))
}

// Delete handles DELETE /api/v1/expenses/:id
func (h *ExpenseHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	expense, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get expense")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expense"})
		return
	}

	if expense == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
		return
	}

	// Verify tenant membership
	if expense.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Soft delete
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error().Err(err).Msg("failed to delete expense")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete expense"})
		return
	}

	// Publish event
	event := domain.ExpenseDeletedEvent{
		ExpenseID: expense.ID,
		User:      expense.UserID,
		Amount:    expense.Amount,
		Timestamp: time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish ExpenseDeletedEvent")
	}

	c.JSON(http.StatusNoContent, nil)
}
