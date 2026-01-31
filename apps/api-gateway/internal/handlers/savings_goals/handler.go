package savings_goals

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
)

// Handler handles HTTP requests for savings goals (Adapter pattern)
type Handler struct {
	useCase ports.SavingsGoalUseCase
}

// NewHandler creates a new savings goals handler (Factory pattern)
func NewHandler(useCase ports.SavingsGoalUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// getUserIDFromContext extrae el user_id del contexto JWT
func (h *Handler) getUserIDFromContext(c *gin.Context) (string, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", errors.NewUnauthorizedRequest("Usuario no autenticado")
	}

	// Convertir a string según el tipo
	switch v := userIDInterface.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case string:
		return v, nil
	default:
		return "", errors.NewBadRequest("Formato de user_id inválido")
	}
}

// CreateGoal handles POST /api/v1/savings-goals
func (h *Handler) CreateGoal(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	var request ports.CreateSavingsGoalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, errors.NewBadRequest("Invalid request body"))
		return
	}

	request.UserID = userID

	// Debug: log the incoming payload to help diagnose 500s
	log.Printf("[SavingsGoals] CreateGoal request: user_id=%s name=%q target_amount=%.2f category=%q priority=%q target_date=%v is_auto_save=%v auto_save_amount=%.2f auto_save_frequency=%q",
		request.UserID, request.Name, request.TargetAmount, request.Category, request.Priority, request.TargetDate, request.IsAutoSave, request.AutoSaveAmount, request.AutoSaveFrequency,
	)

	response, err := h.useCase.CreateGoal(c.Request.Context(), request)
	if err != nil {
		// Map validation errors to 400 for clearer feedback in frontend
		if strings.Contains(strings.ToLower(err.Error()), "validation failed") {
			httpUtil.BadRequest(c, errors.NewBadRequest(err.Error()))
			return
		}
		if strings.Contains(strings.ToLower(err.Error()), "failed to build goal") {
			httpUtil.BadRequest(c, errors.NewBadRequest(err.Error()))
			return
		}
		// Default: internal error
		log.Printf("[SavingsGoals] CreateGoal error: %v", err)
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusCreated, response, "Savings goal created successfully")
}

// GetGoal handles GET /api/v1/savings-goals/:id
func (h *Handler) GetGoal(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	request := ports.GetSavingsGoalRequest{
		UserID: userID,
		GoalID: goalID,
	}

	response, err := h.useCase.GetGoal(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Savings goal retrieved successfully")
}

// ListGoals handles GET /api/v1/savings-goals
func (h *Handler) ListGoals(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	request := ports.ListSavingsGoalsRequest{
		UserID: userID,
	}

	// Parse query parameters for filtering
	if status := c.Query("status"); status != "" {
		request.Status = domain.SavingsGoalStatus(status)
	}

	if category := c.Query("category"); category != "" {
		request.Category = domain.SavingsGoalCategory(category)
	}

	if priority := c.Query("priority"); priority != "" {
		request.Priority = domain.SavingsGoalPriority(priority)
	}

	response, err := h.useCase.ListGoals(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Savings goals retrieved successfully")
}

// UpdateGoal handles PUT /api/v1/savings-goals/:id
func (h *Handler) UpdateGoal(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	var request ports.UpdateSavingsGoalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, errors.NewBadRequest("Invalid request body"))
		return
	}

	request.UserID = userID
	request.GoalID = goalID

	// Debug: log incoming update payload for diagnostics
	log.Printf("[SavingsGoals] UpdateGoal request: user_id=%s goal_id=%s name=%v target_amount=%v category=%v priority=%v target_date=%v is_auto_save=%v auto_save_amount=%v auto_save_frequency=%v image_url_set=%v",
		request.UserID, request.GoalID, request.Name, request.TargetAmount, request.Category, request.Priority, request.TargetDate, request.IsAutoSave, request.AutoSaveAmount, request.AutoSaveFrequency, request.ImageURL != nil)

	response, err := h.useCase.UpdateGoal(c.Request.Context(), request)
	if err != nil {
		// Map validation errors to 400 for clearer feedback
		if strings.Contains(strings.ToLower(err.Error()), "validation failed") ||
			strings.Contains(strings.ToLower(err.Error()), "cannot ") ||
			strings.Contains(strings.ToLower(err.Error()), "must be") {
			httpUtil.BadRequest(c, errors.NewBadRequest(err.Error()))
			return
		}
		log.Printf("[SavingsGoals] UpdateGoal error: %v", err)
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Savings goal updated successfully")
}

// DeleteGoal handles DELETE /api/v1/savings-goals/:id
func (h *Handler) DeleteGoal(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	request := ports.DeleteSavingsGoalRequest{
		UserID: userID,
		GoalID: goalID,
	}

	err = h.useCase.DeleteGoal(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, gin.H{}, "Savings goal deleted successfully")
}

// AddSavings handles POST /api/v1/savings-goals/:id/add-savings
func (h *Handler) AddSavings(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	var request ports.AddSavingsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, errors.NewBadRequest("Invalid request body"))
		return
	}

	request.UserID = userID
	request.GoalID = goalID

	response, err := h.useCase.AddSavings(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Savings added successfully")
}

// WithdrawSavings handles POST /api/v1/savings-goals/:id/withdraw-savings
func (h *Handler) WithdrawSavings(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	var request ports.WithdrawSavingsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, errors.NewBadRequest("Invalid request body"))
		return
	}

	request.UserID = userID
	request.GoalID = goalID

	response, err := h.useCase.WithdrawSavings(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Savings withdrawn successfully")
}

// PauseGoal handles POST /api/v1/savings-goals/:id/pause
func (h *Handler) PauseGoal(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	request := ports.PauseGoalRequest{
		UserID: userID,
		GoalID: goalID,
	}

	err = h.useCase.PauseGoal(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, gin.H{}, "Savings goal paused successfully")
}

// ResumeGoal handles POST /api/v1/savings-goals/:id/resume
func (h *Handler) ResumeGoal(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	request := ports.ResumeGoalRequest{
		UserID: userID,
		GoalID: goalID,
	}

	err = h.useCase.ResumeGoal(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, gin.H{}, "Savings goal resumed successfully")
}

// CancelGoal handles POST /api/v1/savings-goals/:id/cancel
func (h *Handler) CancelGoal(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	request := ports.CancelGoalRequest{
		UserID: userID,
		GoalID: goalID,
	}

	err = h.useCase.CancelGoal(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, gin.H{}, "Savings goal cancelled successfully")
}

// GetGoalSummary handles GET /api/v1/savings-goals/summary
func (h *Handler) GetGoalSummary(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	request := ports.GetGoalSummaryRequest{
		UserID: userID,
	}

	// Parse optional category filter
	if category := c.Query("category"); category != "" {
		request.Category = domain.SavingsGoalCategory(category)
	}

	response, err := h.useCase.GetGoalSummary(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Goal summary retrieved successfully")
}

// GetSavingsGoalsDashboard handles GET /api/v1/savings-goals/dashboard
func (h *Handler) GetSavingsGoalsDashboard(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Get goal summary which contains dashboard data
	request := ports.GetGoalSummaryRequest{
		UserID: userID,
	}

	response, err := h.useCase.GetGoalSummary(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "dashboard de metas de ahorro obtenido exitosamente")
}

// GetGoalTransactions handles GET /api/v1/savings-goals/:id/transactions
func (h *Handler) GetGoalTransactions(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	goalID := c.Param("id")
	if goalID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Goal ID is required"))
		return
	}

	request := ports.GetGoalTransactionsRequest{
		UserID: userID,
		GoalID: goalID,
	}

	// Parse pagination parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			request.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			request.Offset = offset
		}
	}

	response, err := h.useCase.GetGoalTransactions(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Goal transactions retrieved successfully")
}

// RegisterRoutes registers all savings goal routes (Router pattern)
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	savingsGoals := router.Group("/savings-goals")
	{
		// Main CRUD operations
		savingsGoals.POST("", h.CreateGoal)                        // POST /api/v1/savings-goals
		savingsGoals.GET("", h.ListGoals)                          // GET /api/v1/savings-goals
		savingsGoals.GET("/dashboard", h.GetSavingsGoalsDashboard) // GET /api/v1/savings-goals/dashboard
		savingsGoals.GET("/summary", h.GetGoalSummary)             // GET /api/v1/savings-goals/summary
		savingsGoals.GET("/:id", h.GetGoal)                        // GET /api/v1/savings-goals/:id
		savingsGoals.PUT("/:id", h.UpdateGoal)                     // PUT /api/v1/savings-goals/:id
		savingsGoals.DELETE("/:id", h.DeleteGoal)                  // DELETE /api/v1/savings-goals/:id

		// Savings operations
		savingsGoals.POST("/:id/add-savings", h.AddSavings)           // POST /api/v1/savings-goals/:id/add-savings
		savingsGoals.POST("/:id/withdraw-savings", h.WithdrawSavings) // POST /api/v1/savings-goals/:id/withdraw-savings

		// Goal management operations
		savingsGoals.POST("/:id/pause", h.PauseGoal)   // POST /api/v1/savings-goals/:id/pause
		savingsGoals.POST("/:id/resume", h.ResumeGoal) // POST /api/v1/savings-goals/:id/resume
		savingsGoals.POST("/:id/cancel", h.CancelGoal) // POST /api/v1/savings-goals/:id/cancel

		// Analytics and reporting
		savingsGoals.GET("/:id/transactions", h.GetGoalTransactions) // GET /api/v1/savings-goals/:id/transactions
	}
}
