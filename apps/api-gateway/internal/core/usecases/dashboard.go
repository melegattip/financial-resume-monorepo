package usecases

import (
	"context"
	"time"
)

// DashboardUseCase define el caso de uso principal del dashboard
type DashboardUseCase interface {
	GetDashboardOverview(ctx context.Context, params DashboardParams) (*DashboardOverview, error)
}

// DashboardParams representa los parámetros de entrada
type DashboardParams struct {
	UserID string
	Period DatePeriod
}

// DatePeriod representa un período de fechas
type DatePeriod struct {
	Year  *int
	Month *int
}

// DashboardOverview representa la vista general del dashboard
type DashboardOverview struct {
	Period  PeriodInfo
	Metrics FinancialMetrics
	Counts  TransactionCounts
	Trends  FinancialTrends
}

// PeriodInfo contiene información del período consultado
type PeriodInfo struct {
	Year    string
	Month   string
	Label   string
	HasData bool
}

// FinancialMetrics contiene las métricas financieras principales
type FinancialMetrics struct {
	TotalIncome          float64
	TotalExpenses        float64
	Balance              float64
	PendingExpenses      float64
	PendingExpensesCount int
}

// TransactionCounts contiene contadores de transacciones
type TransactionCounts struct {
	IncomeTransactions  int
	ExpenseTransactions int
	CategoriesUsed      int
}

// FinancialTrends contiene información de tendencias
type FinancialTrends struct {
	ExpensePercentageOfIncome float64
	MonthOverMonthChange      float64
}

// PeriodCalculator define la interfaz para cálculos de período
type PeriodCalculator interface {
	FilterTransactionsByPeriod(transactions []Transaction, period DatePeriod) []Transaction
	CalculateMetrics(transactions []Transaction) FinancialMetrics
	FormatPeriodLabel(period DatePeriod) string
}

// Transaction representa una transacción genérica
type Transaction interface {
	GetID() string
	GetUserID() string
	GetAmount() float64
	GetCreatedAt() time.Time
	GetType() TransactionType
	GetCategoryID() string
	IsPending() bool
}

// TransactionType representa el tipo de transacción
type TransactionType string

const (
	IncomeTransaction  TransactionType = "income"
	ExpenseTransaction TransactionType = "expense"
)
