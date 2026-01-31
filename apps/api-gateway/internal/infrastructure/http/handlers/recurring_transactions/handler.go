package recurring_transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// Handler handles HTTP requests for recurring transactions
type Handler struct {
	useCase ports.RecurringTransactionUseCase
}

// NewHandler creates a new recurring transactions handler
func NewHandler(useCase ports.RecurringTransactionUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registers all recurring transaction routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	recurring := router.Group("/recurring-transactions")
	{
		// CRUD operations
		recurring.POST("", h.CreateRecurringTransaction)
		recurring.GET("", h.ListRecurringTransactions)
		recurring.GET("/:id", h.GetRecurringTransaction)
		recurring.PUT("/:id", h.UpdateRecurringTransaction)
		recurring.DELETE("/:id", h.DeleteRecurringTransaction)

		// Transaction control
		recurring.POST("/:id/pause", h.PauseRecurringTransaction)
		recurring.POST("/:id/resume", h.ResumeRecurringTransaction)
		recurring.POST("/:id/execute", h.ExecuteRecurringTransaction)

		// Analytics and dashboard
		recurring.GET("/dashboard", h.GetDashboard)
		recurring.GET("/projection", h.GetCashFlowProjection)

		// Admin operations (for batch processing)
		recurring.POST("/batch/process", h.ProcessPendingTransactions)
		recurring.POST("/batch/notify", h.SendPendingNotifications)
	}
}

// CreateRecurringTransaction creates a new recurring transaction
func (h *Handler) CreateRecurringTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var request ports.CreateRecurringTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Set user ID from token
	request.UserID = userID

	response, err := h.useCase.CreateRecurringTransaction(c.Request.Context(), &request)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Transacción recurrente creada exitosamente",
	})
}

// GetRecurringTransaction gets a recurring transaction by ID
func (h *Handler) GetRecurringTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	response, err := h.useCase.GetRecurringTransaction(c.Request.Context(), userID, transactionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ListRecurringTransactions lists recurring transactions with filters
func (h *Handler) ListRecurringTransactions(c *gin.Context) {
	userID := c.GetString("user_id")

	// Parse query parameters
	filters := ports.RecurringTransactionFilters{
		Type:       c.Query("type"),
		Frequency:  c.Query("frequency"),
		CategoryID: c.Query("category_id"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
	}

	// Parse is_active
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filters.IsActive = &isActive
		}
	}

	// Parse pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	response, err := h.useCase.ListRecurringTransactions(c.Request.Context(), userID, filters)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateRecurringTransaction updates a recurring transaction
func (h *Handler) UpdateRecurringTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	var request ports.UpdateRecurringTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	response, err := h.useCase.UpdateRecurringTransaction(c.Request.Context(), userID, transactionID, &request)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Transacción recurrente actualizada exitosamente",
	})
}

// DeleteRecurringTransaction deletes a recurring transaction
func (h *Handler) DeleteRecurringTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	err := h.useCase.DeleteRecurringTransaction(c.Request.Context(), userID, transactionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transacción recurrente eliminada exitosamente",
	})
}

// PauseRecurringTransaction pauses a recurring transaction
func (h *Handler) PauseRecurringTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	err := h.useCase.PauseRecurringTransaction(c.Request.Context(), userID, transactionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transacción recurrente pausada exitosamente",
	})
}

// ResumeRecurringTransaction resumes a recurring transaction
func (h *Handler) ResumeRecurringTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	err := h.useCase.ResumeRecurringTransaction(c.Request.Context(), userID, transactionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transacción recurrente reanudada exitosamente",
	})
}

// ExecuteRecurringTransaction manually executes a recurring transaction
func (h *Handler) ExecuteRecurringTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	result, err := h.useCase.ExecuteRecurringTransaction(c.Request.Context(), userID, transactionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var message string
	if result.Success {
		message = "Transacción ejecutada exitosamente"
	} else {
		message = "Error ejecutando transacción: " + result.Message
	}

	c.JSON(http.StatusOK, gin.H{
		"success": result.Success,
		"data":    result,
		"message": message,
	})
}

// GetDashboard gets the recurring transactions dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	userID := c.GetString("user_id")

	response, err := h.useCase.GetRecurringTransactionsDashboard(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetCashFlowProjection gets cash flow projection
func (h *Handler) GetCashFlowProjection(c *gin.Context) {
	userID := c.GetString("user_id")

	// Parse months parameter
	months := 6 // Default
	if monthsStr := c.Query("months"); monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m > 0 && m <= 24 {
			months = m
		}
	}

	response, err := h.useCase.GetCashFlowProjection(c.Request.Context(), userID, months)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ProcessPendingTransactions processes all pending recurring transactions (admin)
func (h *Handler) ProcessPendingTransactions(c *gin.Context) {
	// TODO: Add admin authentication check
	result, err := h.useCase.ProcessPendingTransactions(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Procesamiento de transacciones pendientes completado",
	})
}

// SendPendingNotifications sends pending notifications (admin)
func (h *Handler) SendPendingNotifications(c *gin.Context) {
	// TODO: Add admin authentication check
	result, err := h.useCase.SendPendingNotifications(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Envío de notificaciones pendientes completado",
	})
}

// handleError handles different types of errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *gin.Context, err error) {
	// You can implement more sophisticated error handling here
	// For now, we'll use a simple approach

	switch err.Error() {
	case "transacción recurrente no encontrada":
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Transacción recurrente no encontrada",
		})
	case "El ID del usuario es requerido":
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de usuario requerido",
		})
	default:
		// Check if it's a validation error (contains "Error validando")
		if contains(err.Error(), "Error validando") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// Default to internal server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error interno del servidor",
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && containsInMiddle(s, substr)
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
