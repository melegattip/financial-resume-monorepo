package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
)

// CategoriesHandler maneja las peticiones de analytics de categorías
type CategoriesHandler struct {
	service usecases.CategoriesAnalyticsUseCase
}

// NewCategoriesHandler crea una nueva instancia del handler
func NewCategoriesHandler(service usecases.CategoriesAnalyticsUseCase) *CategoriesHandler {
	return &CategoriesHandler{
		service: service,
	}
}

// GetCategoriesAnalytics maneja la petición GET /api/v1/categories/analytics
func (h *CategoriesHandler) GetCategoriesAnalytics(c *gin.Context) {
	// Obtener userID del header
	userID := c.GetHeader("X-Caller-ID")
	if userID == "" {
		httpUtil.HandleError(c, errors.NewBadRequest("X-Caller-ID header es requerido"))
		return
	}

	// Obtener y validar parámetros de query
	params := usecases.CategoriesAnalyticsParams{
		UserID: userID,
	}

	// Parsear año si está presente
	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			httpUtil.HandleError(c, errors.NewBadRequest("Año inválido"))
			return
		}
		params.Period.Year = &year
	}

	// Parsear mes si está presente
	if monthStr := c.Query("month"); monthStr != "" {
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			httpUtil.HandleError(c, errors.NewBadRequest("Mes inválido"))
			return
		}
		if month < 1 || month > 12 {
			httpUtil.HandleError(c, errors.NewBadRequest("Mes debe estar entre 1 y 12"))
			return
		}
		params.Period.Month = &month
	}

	// Llamar al servicio
	response, err := h.service.GetCategoriesAnalytics(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}
