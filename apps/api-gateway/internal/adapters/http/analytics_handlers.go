package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
)

// AnalyticsHandlers maneja todas las peticiones de analytics
type AnalyticsHandlers struct {
	expensesService   usecases.ExpensesAnalyticsUseCase
	categoriesService usecases.CategoriesAnalyticsUseCase
	incomesService    usecases.IncomesAnalyticsUseCase
}

// NewAnalyticsHandlers crea una nueva instancia de los handlers de analytics
func NewAnalyticsHandlers(
	expensesService usecases.ExpensesAnalyticsUseCase,
	categoriesService usecases.CategoriesAnalyticsUseCase,
	incomesService usecases.IncomesAnalyticsUseCase,
) *AnalyticsHandlers {
	return &AnalyticsHandlers{
		expensesService:   expensesService,
		categoriesService: categoriesService,
		incomesService:    incomesService,
	}
}

// getUserIDFromContext extrae el user_id del contexto JWT
func (h *AnalyticsHandlers) getUserIDFromContext(c *gin.Context) (string, error) {
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

// GetExpensesSummary maneja la petición GET /api/v1/analytics/expenses
// @Summary Obtener resumen de gastos
// @Description Obtiene un resumen detallado de los gastos del usuario con análisis financiero
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Año para filtrar (ejemplo: 2024)"
// @Param month query int false "Mes para filtrar (1-12)"
// @Param sort_by query string false "Campo para ordenar (date, amount, category)" Enums(date, amount, category)
// @Param order query string false "Orden de clasificación (asc, desc)" Enums(asc, desc)
// @Param limit query int false "Límite de resultados (máximo 100, por defecto 50)"
// @Param offset query int false "Número de registros a omitir"
// @Success 200 {object} usecases.ExpensesSummaryResponse "Resumen de gastos obtenido exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/analytics/expenses [get]
func (h *AnalyticsHandlers) GetExpensesSummary(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	params := usecases.ExpensesSummaryParams{
		UserID: userID,
		Period: usecases.DatePeriod{},
		Sorting: usecases.SortingCriteria{
			Field: usecases.SortByDate,
			Order: usecases.Descending,
		},
		Pagination: usecases.PaginationParams{
			Limit:  50,
			Offset: 0,
		},
	}

	// Parsear filtros de período
	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			httpUtil.HandleError(c, errors.NewBadRequest("Año inválido"))
			return
		}
		params.Period.Year = &year
	}

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
			httpUtil.HandleError(c, errors.NewBadRequest("Orden inválido"))
			return
		}
	}

	// Parsear paginación
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			httpUtil.HandleError(c, errors.NewBadRequest("Límite inválido (debe ser entre 1 y 100)"))
			return
		}
		params.Pagination.Limit = limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			httpUtil.HandleError(c, errors.NewBadRequest("Offset inválido"))
			return
		}
		params.Pagination.Offset = offset
	}

	response, err := h.expensesService.GetExpensesSummary(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCategoriesAnalytics maneja la petición GET /api/v1/analytics/categories
// @Summary Obtener análisis de categorías
// @Description Obtiene un análisis detallado de gastos agrupados por categorías con métricas y porcentajes
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Año para filtrar (ejemplo: 2024)"
// @Param month query int false "Mes para filtrar (1-12)"
// @Success 200 {object} usecases.CategoriesAnalyticsResponse "Análisis de categorías obtenido exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/analytics/categories [get]
func (h *AnalyticsHandlers) GetCategoriesAnalytics(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	params := usecases.CategoriesAnalyticsParams{
		UserID: userID,
		Period: usecases.DatePeriod{},
	}

	// Parsear filtros de período
	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			httpUtil.HandleError(c, errors.NewBadRequest("Año inválido"))
			return
		}
		params.Period.Year = &year
	}

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

	response, err := h.categoriesService.GetCategoriesAnalytics(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetIncomesSummary maneja la petición GET /api/v1/analytics/incomes
// @Summary Obtener resumen de ingresos
// @Description Obtiene un resumen detallado de los ingresos del usuario con análisis financiero
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Año para filtrar (ejemplo: 2024)"
// @Param month query int false "Mes para filtrar (1-12)"
// @Param sort_by query string false "Campo para ordenar (date, amount, category)" Enums(date, amount, category)
// @Param order query string false "Orden de clasificación (asc, desc)" Enums(asc, desc)
// @Success 200 {object} usecases.IncomesSummaryResponse "Resumen de ingresos obtenido exitosamente"
// @Failure 400 {object} errors.ErrorResponse "Parámetros inválidos"
// @Failure 401 {object} errors.ErrorResponse "No autorizado"
// @Failure 500 {object} errors.ErrorResponse "Error interno del servidor"
// @Router /api/v1/analytics/incomes [get]
func (h *AnalyticsHandlers) GetIncomesSummary(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	params := usecases.IncomesSummaryParams{
		UserID: userID,
		Period: usecases.DatePeriod{},
		Sorting: usecases.SortingCriteria{
			Field: usecases.SortByDate,
			Order: usecases.Descending,
		},
	}

	// Parsear filtros de período
	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			httpUtil.HandleError(c, errors.NewBadRequest("Año inválido"))
			return
		}
		params.Period.Year = &year
	}

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
			httpUtil.HandleError(c, errors.NewBadRequest("Orden inválido"))
			return
		}
	}

	response, err := h.incomesService.GetIncomesSummary(c.Request.Context(), params)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
