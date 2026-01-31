package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/http/helpers"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
)

// Handler maneja las peticiones del dashboard
type Handler struct {
	service            usecases.DashboardUseCase
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de dashboard
func NewHandler(service usecases.DashboardUseCase, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// GetDashboard maneja la petición GET /api/v1/dashboard
// @Summary Obtener resumen del dashboard
// @Description Obtiene un resumen completo del dashboard financiero del usuario incluyendo balance, gastos e ingresos
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Año para filtrar (ejemplo: 2024)"
// @Param month query int false "Mes para filtrar (1-12)"
// @Success 200 {object} usecases.DashboardResponse "Resumen del dashboard obtenido exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/dashboard [get]
func (h *Handler) GetDashboard(c *gin.Context) {
	// Obtener userID del contexto JWT usando helper centralizado
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Construir parámetros usando helpers centralizados
	params := usecases.DashboardParams{
		UserID: userID,
		Period: usecases.DatePeriod{},
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
	response, err := h.service.GetDashboardOverview(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordDashboardView(userID)
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}
