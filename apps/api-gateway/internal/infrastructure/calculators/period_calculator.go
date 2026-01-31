package calculators

import (
	"fmt"

	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// PeriodCalculatorImpl implementa la interfaz PeriodCalculator
type PeriodCalculatorImpl struct{}

// NewPeriodCalculator crea una nueva instancia del calculador de períodos
func NewPeriodCalculator() usecases.PeriodCalculator {
	return &PeriodCalculatorImpl{}
}

// FilterTransactionsByPeriod filtra transacciones por período de fecha
func (p *PeriodCalculatorImpl) FilterTransactionsByPeriod(transactions []usecases.Transaction, period usecases.DatePeriod) []usecases.Transaction {
	if period.Year == nil && period.Month == nil {
		return transactions // Sin filtros
	}

	var filtered []usecases.Transaction
	for _, transaction := range transactions {
		if p.matchesPeriod(transaction, period) {
			filtered = append(filtered, transaction)
		}
	}

	return filtered
}

// CalculateMetrics calcula métricas financieras a partir de las transacciones
func (p *PeriodCalculatorImpl) CalculateMetrics(transactions []usecases.Transaction) usecases.FinancialMetrics {
	var totalIncome, totalExpenses, pendingExpenses float64
	var pendingCount int

	for _, transaction := range transactions {
		switch transaction.GetType() {
		case usecases.IncomeTransaction:
			totalIncome += transaction.GetAmount()
		case usecases.ExpenseTransaction:
			totalExpenses += transaction.GetAmount()
			if transaction.IsPending() {
				pendingExpenses += transaction.GetAmount()
				pendingCount++
			}
		}
	}

	return usecases.FinancialMetrics{
		TotalIncome:          totalIncome,
		TotalExpenses:        totalExpenses,
		Balance:              totalIncome - totalExpenses,
		PendingExpenses:      pendingExpenses,
		PendingExpensesCount: pendingCount,
	}
}

// FormatPeriodLabel formatea una etiqueta legible para el período
func (p *PeriodCalculatorImpl) FormatPeriodLabel(period usecases.DatePeriod) string {
	if period.Year == nil && period.Month == nil {
		return "Todos los períodos"
	}

	if period.Year != nil && period.Month == nil {
		return fmt.Sprintf("Año %d", *period.Year)
	}

	if period.Year != nil && period.Month != nil {
		monthNames := []string{
			"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
			"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
		}

		if *period.Month >= 1 && *period.Month <= 12 {
			return fmt.Sprintf("%s %d", monthNames[*period.Month], *period.Year)
		}
	}

	return "Período inválido"
}

// matchesPeriod verifica si una transacción coincide con el período especificado
func (p *PeriodCalculatorImpl) matchesPeriod(transaction usecases.Transaction, period usecases.DatePeriod) bool {
	date := transaction.GetCreatedAt()

	if period.Year != nil {
		if date.Year() != *period.Year {
			return false
		}
	}

	if period.Month != nil {
		if int(date.Month()) != *period.Month {
			return false
		}
	}

	return true
}
