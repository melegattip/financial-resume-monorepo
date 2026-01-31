package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse estructura estándar para todas las respuestas de la API
type StandardResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta información adicional para paginación y filtros
type Meta struct {
	Total       int    `json:"total,omitempty"`
	Page        int    `json:"page,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
	TotalPages  int    `json:"total_pages,omitempty"`
	HasNextPage bool   `json:"has_next_page,omitempty"`
	FilteredBy  string `json:"filtered_by,omitempty"`
}

// TransactionResponse estructura estandarizada para transacciones
type TransactionResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Type        string  `json:"type"` // "income" | "expense"

	// Campos específicos de gastos
	AmountPaid    *float64 `json:"amount_paid,omitempty"`
	PendingAmount *float64 `json:"pending_amount,omitempty"`
	Paid          *bool    `json:"paid,omitempty"`
	DueDate       *string  `json:"due_date,omitempty"`

	// Campos calculados
	Percentage *float64 `json:"percentage,omitempty"`
}

// CategoryResponse estructura estandarizada para categorías
type CategoryResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// DashboardResponse estructura estandarizada para dashboard
type DashboardResponse struct {
	Period  PeriodInfo       `json:"period"`
	Metrics DashboardMetrics `json:"metrics"`
	Counts  DashboardCounts  `json:"counts"`
	Trends  DashboardTrends  `json:"trends"`
}

type PeriodInfo struct {
	Year    string `json:"year"`
	Month   string `json:"month"`
	Label   string `json:"label"`
	HasData bool   `json:"has_data"`
}

type DashboardMetrics struct {
	TotalIncome          float64 `json:"total_income"`
	TotalExpenses        float64 `json:"total_expenses"`
	Balance              float64 `json:"balance"`
	PendingExpenses      float64 `json:"pending_expenses"`
	PendingExpensesCount int     `json:"pending_expenses_count"`
}

type DashboardCounts struct {
	IncomeTransactions  int `json:"income_transactions"`
	ExpenseTransactions int `json:"expense_transactions"`
	CategoriesUsed      int `json:"categories_used"`
}

type DashboardTrends struct {
	ExpensePercentageOfIncome int `json:"expense_percentage_of_income"`
	MonthOverMonthChange      int `json:"month_over_month_change"`
}

// Funciones de respuesta estandarizadas
func SendSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, StandardResponse{
		Data:    data,
		Message: message,
	})
}

func SendSuccessWithMeta(c *gin.Context, statusCode int, data interface{}, message string, meta *Meta) {
	c.JSON(statusCode, StandardResponse{
		Data:    data,
		Message: message,
		Meta:    meta,
	})
}

func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, StandardResponse{
		Error: message,
	})
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, StandardResponse{
		Error: err.Error(),
	})
}

func Unauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, StandardResponse{
		Error: err.Error(),
	})
}

func NotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, StandardResponse{
		Error: err.Error(),
	})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, StandardResponse{
		Error: "Internal server error",
	})
}

func Conflict(c *gin.Context, err error) {
	c.JSON(http.StatusConflict, StandardResponse{
		Error: err.Error(),
	})
}

func TooManyRequests(c *gin.Context, err error) {
	c.JSON(http.StatusTooManyRequests, StandardResponse{
		Error: err.Error(),
	})
}

func HandleStandardError(c *gin.Context, err error) {
	switch err.(type) {
	case *BadRequestError:
		BadRequest(c, err)
	case *UnauthorizedError:
		Unauthorized(c, err)
	case *NotFoundError:
		NotFound(c, err)
	default:
		InternalServerError(c, err)
	}
}

// Tipos de error personalizados
type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
