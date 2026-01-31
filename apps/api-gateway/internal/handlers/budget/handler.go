package budget

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
)

// Handler handles HTTP requests for budget operations (Adapter pattern)
type Handler struct {
	budgetUseCase ports.BudgetUseCase
}

// NewHandler creates a new budget handler
func NewHandler(budgetUseCase ports.BudgetUseCase) *Handler {
	return &Handler{
		budgetUseCase: budgetUseCase,
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

// CreateBudget handles POST /api/v1/budgets
func (h *Handler) CreateBudget(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	var request ports.CreateBudgetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, errors.NewBadRequest("datos inválidos: "+err.Error()))
		return
	}

	// Set user ID from context
	request.UserID = userID

	// Validate request
	if err := h.validateCreateBudgetRequest(request); err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	response, err := h.budgetUseCase.CreateBudget(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusCreated, response, "presupuesto creado exitosamente")
}

// GetBudget handles GET /api/v1/budgets/:id
func (h *Handler) GetBudget(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	budgetID := c.Param("id")
	if budgetID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("ID de presupuesto requerido"))
		return
	}

	request := ports.GetBudgetRequest{
		UserID:   userID,
		BudgetID: budgetID,
	}

	response, err := h.budgetUseCase.GetBudget(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "presupuesto obtenido exitosamente")
}

// ListBudgets handles GET /api/v1/budgets
func (h *Handler) ListBudgets(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Parse query parameters
	request := ports.ListBudgetsRequest{
		UserID:     userID,
		Period:     domain.BudgetPeriod(c.Query("period")),
		CategoryID: c.Query("category_id"),
		Status:     domain.BudgetStatus(c.Query("status")),
	}

	// Parse active_only parameter
	if activeOnlyStr := c.Query("active_only"); activeOnlyStr != "" {
		activeOnly, err := strconv.ParseBool(activeOnlyStr)
		if err != nil {
			httpUtil.BadRequest(c, errors.NewBadRequest("parámetro active_only inválido"))
			return
		}
		request.ActiveOnly = activeOnly
	}

	response, err := h.budgetUseCase.ListBudgets(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "presupuestos obtenidos exitosamente")
}

// UpdateBudget handles PUT /api/v1/budgets/:id
func (h *Handler) UpdateBudget(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	budgetID := c.Param("id")
	if budgetID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("ID de presupuesto requerido"))
		return
	}

	var requestBody struct {
		Amount   *float64 `json:"amount,omitempty"`
		AlertAt  *float64 `json:"alert_at,omitempty"`
		IsActive *bool    `json:"is_active,omitempty"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		httpUtil.BadRequest(c, errors.NewBadRequest("datos inválidos: "+err.Error()))
		return
	}

	request := ports.UpdateBudgetRequest{
		UserID:   userID,
		BudgetID: budgetID,
		Amount:   requestBody.Amount,
		AlertAt:  requestBody.AlertAt,
		IsActive: requestBody.IsActive,
	}

	// Validate request
	if err := h.validateUpdateBudgetRequest(request); err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	response, err := h.budgetUseCase.UpdateBudget(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "presupuesto actualizado exitosamente")
}

// DeleteBudget handles DELETE /api/v1/budgets/:id
func (h *Handler) DeleteBudget(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	budgetID := c.Param("id")
	if budgetID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("ID de presupuesto requerido"))
		return
	}

	request := ports.DeleteBudgetRequest{
		UserID:   userID,
		BudgetID: budgetID,
	}

	err = h.budgetUseCase.DeleteBudget(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, nil, "presupuesto eliminado exitosamente")
}

// GetBudgetStatus handles GET /api/v1/budgets/status
func (h *Handler) GetBudgetStatus(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	request := ports.GetBudgetStatusRequest{
		UserID:     userID,
		CategoryID: c.Query("category_id"),
		Period:     domain.BudgetPeriod(c.Query("period")),
	}

	response, err := h.budgetUseCase.GetBudgetStatus(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "estado de presupuestos obtenido exitosamente")
}

// RefreshBudgetAmounts handles POST /api/v1/budgets/refresh
func (h *Handler) RefreshBudgetAmounts(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	err = h.budgetUseCase.RefreshBudgetAmounts(c.Request.Context(), userID)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, nil, "presupuestos actualizados exitosamente")
}

// GetBudgetDashboard handles GET /api/v1/budgets/dashboard
func (h *Handler) GetBudgetDashboard(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Get budget status which contains dashboard data
	request := ports.GetBudgetStatusRequest{
		UserID: userID,
	}

	response, err := h.budgetUseCase.GetBudgetStatus(c.Request.Context(), request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "dashboard de presupuestos obtenido exitosamente")
}

// Validation methods (Single Responsibility Principle)

func (h *Handler) validateCreateBudgetRequest(request ports.CreateBudgetRequest) error {
	if request.UserID == "" {
		return errors.NewBadRequest("ID de usuario requerido")
	}

	if request.CategoryID == "" {
		return errors.NewBadRequest("ID de categoría requerido")
	}

	if request.Amount <= 0 {
		return errors.NewBadRequest("el monto del presupuesto debe ser mayor a 0")
	}

	if request.Amount > 1e12 { // 1 trillion
		return errors.NewBadRequest("el monto del presupuesto es demasiado grande")
	}

	validPeriods := map[domain.BudgetPeriod]bool{
		domain.BudgetPeriodMonthly: true,
		domain.BudgetPeriodWeekly:  true,
		domain.BudgetPeriodYearly:  true,
	}

	if !validPeriods[request.Period] {
		return errors.NewBadRequest("período de presupuesto inválido. Debe ser 'monthly', 'weekly' o 'yearly'")
	}

	if request.AlertAt < 0 || request.AlertAt > 1 {
		return errors.NewBadRequest("umbral de alerta debe estar entre 0 y 1")
	}

	return nil
}

func (h *Handler) validateUpdateBudgetRequest(request ports.UpdateBudgetRequest) error {
	if request.UserID == "" {
		return errors.NewBadRequest("ID de usuario requerido")
	}

	if request.BudgetID == "" {
		return errors.NewBadRequest("ID de presupuesto requerido")
	}

	if request.Amount != nil {
		if *request.Amount <= 0 {
			return errors.NewBadRequest("el monto del presupuesto debe ser mayor a 0")
		}
		if *request.Amount > 1e12 {
			return errors.NewBadRequest("el monto del presupuesto es demasiado grande")
		}
	}

	if request.AlertAt != nil {
		if *request.AlertAt < 0 || *request.AlertAt > 1 {
			return errors.NewBadRequest("umbral de alerta debe estar entre 0 y 1")
		}
	}

	return nil
}

// RegisterRoutes registers budget routes with the router (Router pattern)
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	budgets := router.Group("/budgets")
	{
		budgets.POST("", h.CreateBudget)                 // POST /api/v1/budgets
		budgets.GET("", h.ListBudgets)                   // GET /api/v1/budgets
		budgets.GET("/dashboard", h.GetBudgetDashboard)  // GET /api/v1/budgets/dashboard
		budgets.GET("/status", h.GetBudgetStatus)        // GET /api/v1/budgets/status
		budgets.POST("/refresh", h.RefreshBudgetAmounts) // POST /api/v1/budgets/refresh
		budgets.GET("/:id", h.GetBudget)                 // GET /api/v1/budgets/:id
		budgets.PUT("/:id", h.UpdateBudget)              // PUT /api/v1/budgets/:id
		budgets.DELETE("/:id", h.DeleteBudget)           // DELETE /api/v1/budgets/:id
	}
}
