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

type IncomeHandler struct {
	repo     ports.IncomeRepository
	eventBus ports.EventBus
	logger   zerolog.Logger
}

func NewIncomeHandler(repo ports.IncomeRepository, eventBus ports.EventBus, logger zerolog.Logger) *IncomeHandler {
	return &IncomeHandler{
		repo:     repo,
		eventBus: eventBus,
		logger:   logger,
	}
}

// CreateIncomeRequest is the request body for creating an income.
// source and received_date are optional; sensible defaults are applied when absent.
type CreateIncomeRequest struct {
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Source       string  `json:"source"`
	Description  string  `json:"description" binding:"required"`
	ReceivedDate string  `json:"received_date"`
	CategoryID   string  `json:"category_id"` // accepted but ignored (incomes have no category)
}

// UpdateIncomeRequest is the request body for updating an income.
type UpdateIncomeRequest struct {
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Source       string  `json:"source"`
	Description  string  `json:"description" binding:"required"`
	ReceivedDate string  `json:"received_date"`
	CategoryID   string  `json:"category_id"` // accepted but ignored
}

// IncomeResponse is the response format for an income
type IncomeResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	Amount       float64 `json:"amount"`
	Source       string  `json:"source"`
	Description  string  `json:"description"`
	ReceivedDate string  `json:"received_date"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func toIncomeResponse(i *domain.Income) IncomeResponse {
	return IncomeResponse{
		ID:           i.ID,
		UserID:       i.UserID,
		Amount:       i.Amount,
		Source:       i.Source,
		Description:  i.Description,
		ReceivedDate: i.ReceivedDate.Format(time.RFC3339),
		CreatedAt:    i.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    i.UpdatedAt.Format(time.RFC3339),
	}
}

// Create handles POST /api/v1/incomes
func (h *IncomeHandler) Create(c *gin.Context) {
	var req CreateIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := c.GetString("tenant_id")

	// Default source to description when not provided.
	source := req.Source
	if source == "" {
		source = req.Description
	}

	// Default received_date to now when not provided.
	var receivedDate time.Time
	if req.ReceivedDate == "" {
		receivedDate = time.Now().UTC()
	} else {
		var err error
		receivedDate, err = time.Parse(time.RFC3339, req.ReceivedDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid received_date format, expected RFC3339"})
			return
		}
	}

	income, err := domain.NewIncome(userID.(string), req.Amount, source, req.Description, receivedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	income.TenantID = tenantID

	if err := h.repo.Create(c.Request.Context(), income); err != nil {
		h.logger.Error().Err(err).Msg("failed to create income")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create income"})
		return
	}

	event := domain.IncomeCreatedEvent{
		IncomeID:     income.ID,
		User:         income.UserID,
		Amount:       income.Amount,
		Source:       income.Source,
		Description:  income.Description,
		ReceivedDate: income.ReceivedDate,
		Timestamp:    time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish IncomeCreatedEvent")
	}

	c.JSON(http.StatusCreated, toIncomeResponse(income))
}

// List handles GET /api/v1/incomes
func (h *IncomeHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	incomes, err := h.repo.FindByTenantID(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list incomes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list incomes"})
		return
	}

	response := make([]IncomeResponse, len(incomes))
	for i, income := range incomes {
		response[i] = toIncomeResponse(income)
	}

	c.JSON(http.StatusOK, gin.H{
		"incomes": response,
		"total":   len(response),
		"limit":   limit,
		"offset":  offset,
	})
}

// GetByID handles GET /api/v1/incomes/:id
func (h *IncomeHandler) GetByID(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	income, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get income")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get income"})
		return
	}

	if income == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "income not found"})
		return
	}

	if income.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.JSON(http.StatusOK, toIncomeResponse(income))
}

// Update handles PUT /api/v1/incomes/:id
func (h *IncomeHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	income, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get income")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get income"})
		return
	}

	if income == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "income not found"})
		return
	}

	if income.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req UpdateIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source := req.Source
	if source == "" {
		source = req.Description
	}

	var receivedDate time.Time
	if req.ReceivedDate == "" {
		receivedDate = income.ReceivedDate // keep existing date
	} else {
		var parseErr error
		receivedDate, parseErr = time.Parse(time.RFC3339, req.ReceivedDate)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid received_date format, expected RFC3339"})
			return
		}
	}

	if err := income.Update(req.Amount, source, req.Description, receivedDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), income); err != nil {
		h.logger.Error().Err(err).Msg("failed to update income")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update income"})
		return
	}

	event := domain.IncomeUpdatedEvent{
		IncomeID:  income.ID,
		User:      income.UserID,
		Timestamp: time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish IncomeUpdatedEvent")
	}

	c.JSON(http.StatusOK, toIncomeResponse(income))
}

// Delete handles DELETE /api/v1/incomes/:id
func (h *IncomeHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	income, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get income")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get income"})
		return
	}

	if income == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "income not found"})
		return
	}

	if income.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error().Err(err).Msg("failed to delete income")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete income"})
		return
	}

	event := domain.IncomeDeletedEvent{
		IncomeID:  income.ID,
		User:      income.UserID,
		Amount:    income.Amount,
		Timestamp: time.Now().UTC(),
	}
	if err := h.eventBus.Publish(c.Request.Context(), event); err != nil {
		h.logger.Warn().Err(err).Msg("failed to publish IncomeDeletedEvent")
	}

	c.JSON(http.StatusNoContent, nil)
}
