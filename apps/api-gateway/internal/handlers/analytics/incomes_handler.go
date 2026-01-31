package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
)

// IncomesHandler maneja las peticiones de analytics de ingresos
type IncomesHandler struct {
	service usecases.IncomesAnalyticsUseCase
}

// NewIncomesHandler crea una nueva instancia del handler
func NewIncomesHandler(service usecases.IncomesAnalyticsUseCase) *IncomesHandler {
	return &IncomesHandler{
		service: service,
	}
}

// GetIncomesSummary maneja la petición GET /api/v1/incomes/summary
func (h *IncomesHandler) GetIncomesSummary(c *gin.Context) {
	// Obtener userID del header
	userID := c.GetHeader("X-Caller-ID")
	if userID == "" {
		httpUtil.HandleError(c, errors.NewBadRequest("X-Caller-ID header es requerido"))
		return
	}

	// Obtener y validar parámetros de query
	params := usecases.IncomesSummaryParams{
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

	// Parsear criterios de ordenamiento
	if sortBy := c.Query("sort_by"); sortBy != "" {
		switch sortBy {
		case "date":
			params.Sorting.Field = usecases.SortByDate
		case "amount":
			params.Sorting.Field = usecases.SortByAmount
		case "category":
			params.Sorting.Field = usecases.SortByCategory
		default:
			httpUtil.HandleError(c, errors.NewBadRequest("Campo de ordenamiento inválido"))
			return
		}
	}

	if order := c.Query("order"); order != "" {
		switch order {
		case "asc":
			params.Sorting.Order = usecases.Ascending
		case "desc":
			params.Sorting.Order = usecases.Descending
		default:
			httpUtil.HandleError(c, errors.NewBadRequest("Orden inválido (debe ser asc o desc)"))
			return
		}
	}

	// Llamar al servicio
	response, err := h.service.GetIncomesSummary(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}
