package insights

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/http/helpers"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/insights"
)

// Handler maneja las peticiones de insights financieros
type Handler struct {
	service            insights.InsightsUseCase
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de insights
func NewHandler(service insights.InsightsUseCase, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// MarkInsightAsUnderstoodRequest estructura para el request del endpoint
type MarkInsightAsUnderstoodRequest struct {
	InsightID   string `json:"insight_id" binding:"required"`
	InsightType string `json:"insight_type"`
	Description string `json:"description"`
}

// MarkInsightAsUnderstood maneja el botón "Marcar como revisado" del frontend
// @Summary Marcar insight como entendido
// @Description Registra que el usuario ha entendido/revisado un insight específico
// @Tags Insights
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param insight body MarkInsightAsUnderstoodRequest true "Datos del insight revisado"
// @Success 200 {object} map[string]interface{} "Insight marcado como entendido exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Datos de entrada inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Router /api/v1/insights/mark-understood [post]
func (h *Handler) MarkInsightAsUnderstood(c *gin.Context) {
	// Obtener userID del contexto JWT
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Parsear request body
	var request MarkInsightAsUnderstoodRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.HandleError(c, errors.NewBadRequest("Datos de entrada inválidos: "+err.Error()))
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		description := request.Description
		if description == "" {
			description = "Usuario entendió insight: " + request.InsightType
		}

		h.gamificationHelper.RecordInsightInteraction(
			userID,
			services.ActionUnderstandInsight,
			request.InsightID,
			description,
		)
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Insight marcado como entendido exitosamente",
		"data": gin.H{
			"insight_id":   request.InsightID,
			"insight_type": request.InsightType,
			"user_id":      userID,
		},
	})
}

// GetFinancialHealth maneja la petición GET /api/v1/insights/financial-health
// @Summary Obtener análisis de salud financiera
// @Description Obtiene un análisis completo de la salud financiera del usuario basado en sus transacciones reales
// @Tags Insights
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Año para filtrar (ejemplo: 2024)"
// @Param month query int false "Mes para filtrar (1-12)"
// @Success 200 {object} insights.FinancialHealthResponse "Análisis de salud financiera obtenido exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/insights/financial-health [get]
func (h *Handler) GetFinancialHealth(c *gin.Context) {
	// Obtener userID del contexto JWT usando helper centralizado
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Construir parámetros usando helpers centralizados
	params := insights.InsightsParams{
		UserID: userID,
		Period: insights.DatePeriod{},
	}

	// Parsear año si está presente
	year, err := helpers.ParseIntQuery(c, "year", false)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}
	params.Period.Year = year

	// Parsear mes si está presente
	month, err := helpers.ParseIntQuery(c, "month", false)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}
	params.Period.Month = month

	// Validar parámetros de fecha
	if err := helpers.ValidateYearMonth(year, month); err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Llamar al servicio
	response, err := h.service.GetFinancialHealth(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordInsightInteraction(
			userID,
			services.ActionViewInsight,
			"financial_health",
			"Usuario visualizó análisis de salud financiera",
		)
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}

// GetAIInsights maneja la petición GET /api/v1/ai/insights
// @Summary Obtener insights generados por IA
// @Description Obtiene insights financieros personalizados generados por inteligencia artificial
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Año para filtrar (ejemplo: 2024)"
// @Param month query int false "Mes para filtrar (1-12)"
// @Success 200 {object} insights.AIInsightsResponse "Insights de IA obtenidos exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/ai/insights [get]
func (h *Handler) GetAIInsights(c *gin.Context) {
	// Obtener userID del contexto JWT usando helper centralizado
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Construir parámetros usando helpers centralizados
	params := insights.InsightsParams{
		UserID: userID,
		Period: insights.DatePeriod{},
	}

	// Parsear año si está presente
	year, err := helpers.ParseIntQuery(c, "year", false)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}
	params.Period.Year = year

	// Parsear mes si está presente
	month, err := helpers.ParseIntQuery(c, "month", false)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}
	params.Period.Month = month

	// Validar parámetros de fecha
	if err := helpers.ValidateYearMonth(year, month); err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Llamar al servicio
	response, err := h.service.GetAIInsights(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordInsightInteraction(
			userID,
			services.ActionViewInsight,
			"ai_insights",
			"Usuario visualizó insights de IA",
		)
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}

// CanIBuy maneja la petición POST /api/v1/ai/can-i-buy
// @Summary Analizar si puedes permitirte una compra
// @Description Analiza tu situación financiera para determinar si puedes permitirte una compra específica
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param purchase body insights.CanIBuyParams true "Datos de la compra a analizar"
// @Success 200 {object} insights.CanIBuyResponse "Análisis de compra completado exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Datos de entrada inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/ai/can-i-buy [post]
func (h *Handler) CanIBuy(c *gin.Context) {
	// Obtener userID del contexto JWT usando helper centralizado
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Parsear request body
	var params insights.CanIBuyParams
	if err := c.ShouldBindJSON(&params); err != nil {
		httpUtil.HandleError(c, errors.NewBadRequest("Datos de entrada inválidos: "+err.Error()))
		return
	}

	// Asegurar que el userID coincida
	params.UserID = userID

	// Llamar al servicio
	response, err := h.service.CanIBuy(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordInsightInteraction(
			userID,
			services.ActionViewInsight,
			"can_i_buy",
			"Usuario consultó análisis de compra: "+params.ItemName,
		)
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}

// GetCreditImprovementPlan maneja la petición GET /api/v1/ai/credit-improvement-plan
// @Summary Obtener plan de mejora crediticia personalizado
// @Description Genera un plan detallado usando IA para mejorar tu perfil crediticio y salud financiera
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Año para filtrar (ejemplo: 2024)"
// @Param month query int false "Mes para filtrar (1-12)"
// @Success 200 {object} insights.CreditImprovementPlanResponse "Plan de mejora generado exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/ai/credit-improvement-plan [get]
func (h *Handler) GetCreditImprovementPlan(c *gin.Context) {
	// Obtener userID del contexto JWT usando helper centralizado
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Construir parámetros usando helpers centralizados
	params := insights.InsightsParams{
		UserID: userID,
		Period: insights.DatePeriod{},
	}

	// Parsear año si está presente
	year, err := helpers.ParseIntQuery(c, "year", false)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}
	params.Period.Year = year

	// Parsear mes si está presente
	month, err := helpers.ParseIntQuery(c, "month", false)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}
	params.Period.Month = month

	// Validar parámetros de fecha
	if err := helpers.ValidateYearMonth(year, month); err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Llamar al servicio
	response, err := h.service.GetCreditImprovementPlan(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordInsightInteraction(
			userID,
			services.ActionViewInsight,
			"credit_improvement_plan",
			"Usuario consultó plan de mejora crediticia",
		)
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}
