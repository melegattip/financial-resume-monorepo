package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
)

// ExpensesHandler maneja las peticiones de analytics de gastos
type ExpensesHandler struct {
	service usecases.ExpensesAnalyticsUseCase
}

// NewExpensesHandler crea una nueva instancia del handler
func NewExpensesHandler(service usecases.ExpensesAnalyticsUseCase) *ExpensesHandler {
	return &ExpensesHandler{
		service: service,
	}
}

// GetExpensesSummary maneja la petición GET /api/v1/expenses/summary
func (h *ExpensesHandler) GetExpensesSummary(c *gin.Context) {
	// Obtener userID del header
	userID := c.GetHeader("X-Caller-ID")
	if userID == "" {
		httpUtil.HandleError(c, errors.NewBadRequest("X-Caller-ID header es requerido"))
		return
	}

	// Obtener y validar parámetros de query
	params := usecases.ExpensesSummaryParams{
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

	// Parsear paginación
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			httpUtil.HandleError(c, errors.NewBadRequest("Límite inválido"))
			return
		}
		if limit < 1 || limit > 1000 {
			httpUtil.HandleError(c, errors.NewBadRequest("Límite inválido (debe ser entre 1 y 1000)"))
			return
		}
		params.Pagination.Limit = limit
	} else {
		params.Pagination.Limit = 50 // default
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			httpUtil.HandleError(c, errors.NewBadRequest("Offset inválido"))
			return
		}
		if offset < 0 {
			httpUtil.HandleError(c, errors.NewBadRequest("Offset inválido (debe ser >= 0)"))
			return
		}
		params.Pagination.Offset = offset
	}

	// Llamar al servicio
	response, err := h.service.GetExpensesSummary(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// Devolver respuesta
	c.JSON(http.StatusOK, response)
}
