package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// DashboardHandler maneja las peticiones HTTP del dashboard
type DashboardHandler struct {
	dashboardUseCase usecases.DashboardUseCase
}

// NewDashboardHandler crea un nuevo handler del dashboard
func NewDashboardHandler(dashboardUseCase usecases.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUseCase: dashboardUseCase,
	}
}

// GetDashboard maneja la petición GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// Extraer parámetros de la petición
	params, err := h.buildDashboardParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ejecutar caso de uso
	result, err := h.dashboardUseCase.GetDashboardOverview(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convertir a formato de respuesta HTTP
	response := h.buildHTTPResponse(result)
	c.JSON(http.StatusOK, response)
}

// buildDashboardParams construye los parámetros del caso de uso a partir de la petición HTTP
func (h *DashboardHandler) buildDashboardParams(c *gin.Context) (usecases.DashboardParams, error) {
	params := usecases.DashboardParams{
		UserID: c.Query("user_id"),
		Period: usecases.DatePeriod{},
	}

	// Validar UserID
	if params.UserID == "" {
		params.UserID = c.GetString("user_id") // Desde middleware de autenticación
	}

	// Validar que el UserID esté presente
	if params.UserID == "" {
		return params, fmt.Errorf("user_id es requerido")
	}

	// Parsear año si está presente
	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return params, err
		}
		params.Period.Year = &year
	}

	// Parsear mes si está presente
	if monthStr := c.Query("month"); monthStr != "" {
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return params, err
		}
		params.Period.Month = &month
	}

	return params, nil
}

// buildHTTPResponse convierte la respuesta del caso de uso al formato HTTP
func (h *DashboardHandler) buildHTTPResponse(overview *usecases.DashboardOverview) DashboardResponse {
	return DashboardResponse{
		Period: PeriodInfo{
			Year:    overview.Period.Year,
			Month:   overview.Period.Month,
			Label:   overview.Period.Label,
			HasData: overview.Period.HasData,
		},
		Metrics: MetricsInfo{
			TotalIncome:          overview.Metrics.TotalIncome,
			TotalExpenses:        overview.Metrics.TotalExpenses,
			Balance:              overview.Metrics.Balance,
			PendingExpenses:      overview.Metrics.PendingExpenses,
			PendingExpensesCount: overview.Metrics.PendingExpensesCount,
		},
		Counts: CountsInfo{
			IncomeTransactions:  overview.Counts.IncomeTransactions,
			ExpenseTransactions: overview.Counts.ExpenseTransactions,
			CategoriesUsed:      overview.Counts.CategoriesUsed,
		},
		Trends: TrendsInfo{
			ExpensePercentageOfIncome: overview.Trends.ExpensePercentageOfIncome,
			MonthOverMonthChange:      overview.Trends.MonthOverMonthChange,
		},
	}
}

// DashboardResponse representa la respuesta HTTP del dashboard
type DashboardResponse struct {
	Period  PeriodInfo  `json:"period"`
	Metrics MetricsInfo `json:"metrics"`
	Counts  CountsInfo  `json:"counts"`
	Trends  TrendsInfo  `json:"trends"`
}

// PeriodInfo representa información del período en la respuesta HTTP
type PeriodInfo struct {
	Year    string `json:"year"`
	Month   string `json:"month"`
	Label   string `json:"label"`
	HasData bool   `json:"has_data"`
}

// MetricsInfo representa métricas financieras en la respuesta HTTP
type MetricsInfo struct {
	TotalIncome          float64 `json:"total_income"`
	TotalExpenses        float64 `json:"total_expenses"`
	Balance              float64 `json:"balance"`
	PendingExpenses      float64 `json:"pending_expenses"`
	PendingExpensesCount int     `json:"pending_expenses_count"`
}

// CountsInfo representa contadores en la respuesta HTTP
type CountsInfo struct {
	IncomeTransactions  int `json:"income_transactions"`
	ExpenseTransactions int `json:"expense_transactions"`
	CategoriesUsed      int `json:"categories_used"`
}

// TrendsInfo representa tendencias en la respuesta HTTP
type TrendsInfo struct {
	ExpensePercentageOfIncome float64 `json:"expense_percentage_of_income"`
	MonthOverMonthChange      float64 `json:"month_over_month_change"`
}
