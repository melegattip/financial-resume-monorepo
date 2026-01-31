package dashboard

import (
	"context"
	"fmt"
	"strconv"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// Service implementa DashboardUseCase siguiendo Clean Architecture
type Service struct {
	expenseRepo         baseRepo.ExpenseRepository
	incomeRepo          baseRepo.IncomeRepository
	periodCalculator    usecases.PeriodCalculator
	analyticsCalculator usecases.AnalyticsCalculator
}

// NewService crea una nueva instancia del servicio con dependency injection
func NewService(
	expenseRepo baseRepo.ExpenseRepository,
	incomeRepo baseRepo.IncomeRepository,
	periodCalculator usecases.PeriodCalculator,
	analyticsCalculator usecases.AnalyticsCalculator,
) usecases.DashboardUseCase {
	return &Service{
		expenseRepo:         expenseRepo,
		incomeRepo:          incomeRepo,
		periodCalculator:    periodCalculator,
		analyticsCalculator: analyticsCalculator,
	}
}

// GetDashboardOverview implementa el caso de uso principal
func (s *Service) GetDashboardOverview(ctx context.Context, params usecases.DashboardParams) (*usecases.DashboardOverview, error) {
	// Validar parámetros de entrada
	if err := s.validateParams(params); err != nil {
		return nil, err
	}

	// Obtener transacciones
	transactions, err := s.getTransactions(ctx, params.UserID, params.Period)
	if err != nil {
		return nil, err
	}

	// Calcular métricas usando el calculador
	metrics := s.periodCalculator.CalculateMetrics(transactions)

	// Calcular contadores
	counts := s.calculateCounts(transactions)

	// Calcular tendencias
	trends := s.calculateTrends(metrics)

	// Construir información del período
	period := s.buildPeriodInfo(params.Period, len(transactions) > 0)

	return &usecases.DashboardOverview{
		Period:  period,
		Metrics: metrics,
		Counts:  counts,
		Trends:  trends,
	}, nil
}

// validateParams valida los parámetros de entrada (SRP - Single Responsibility)
func (s *Service) validateParams(params usecases.DashboardParams) error {
	if params.UserID == "" {
		return errors.NewBadRequest("El ID del usuario es requerido")
	}

	if params.Period.Year != nil {
		if *params.Period.Year < 1900 || *params.Period.Year > 2100 {
			return errors.NewBadRequest("Año inválido. Debe estar entre 1900 y 2100")
		}
	}

	if params.Period.Month != nil {
		if *params.Period.Month < 1 || *params.Period.Month > 12 {
			return errors.NewBadRequest("Mes inválido. Debe estar entre 1 y 12")
		}
	}

	return nil
}

// getTransactions obtiene todas las transacciones del usuario (SRP)
func (s *Service) getTransactions(ctx context.Context, userID string, period usecases.DatePeriod) ([]usecases.Transaction, error) {
	// Obtener gastos
	expenses, err := s.expenseRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo gastos: %w", err)
	}

	// Obtener ingresos
	incomes, err := s.incomeRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ingresos: %w", err)
	}

	// Convertir a transacciones genéricas
	var transactions []usecases.Transaction
	for _, expense := range expenses {
		transactions = append(transactions, NewExpenseTransaction(expense))
	}
	for _, income := range incomes {
		transactions = append(transactions, NewIncomeTransaction(income))
	}

	// Filtrar por período usando el calculador
	return s.periodCalculator.FilterTransactionsByPeriod(transactions, period), nil
}

// calculateCounts calcula contadores de transacciones (SRP)
func (s *Service) calculateCounts(transactions []usecases.Transaction) usecases.TransactionCounts {
	var incomeCount, expenseCount int
	categories := make(map[string]bool)

	for _, transaction := range transactions {
		switch transaction.GetType() {
		case usecases.IncomeTransaction:
			incomeCount++
		case usecases.ExpenseTransaction:
			expenseCount++
		}

		if categoryID := transaction.GetCategoryID(); categoryID != "" {
			categories[categoryID] = true
		}
	}

	return usecases.TransactionCounts{
		IncomeTransactions:  incomeCount,
		ExpenseTransactions: expenseCount,
		CategoriesUsed:      len(categories),
	}
}

// calculateTrends calcula tendencias financieras (SRP)
func (s *Service) calculateTrends(metrics usecases.FinancialMetrics) usecases.FinancialTrends {
	expensePercentage := s.analyticsCalculator.CalculatePercentages(metrics.TotalExpenses, metrics.TotalIncome)

	return usecases.FinancialTrends{
		ExpensePercentageOfIncome: expensePercentage,
		MonthOverMonthChange:      0, // TODO: Implementar comparación con período anterior
	}
}

// buildPeriodInfo construye información del período (SRP)
func (s *Service) buildPeriodInfo(period usecases.DatePeriod, hasData bool) usecases.PeriodInfo {
	label := s.periodCalculator.FormatPeriodLabel(period)

	yearStr := ""
	if period.Year != nil {
		yearStr = strconv.Itoa(*period.Year)
	}

	monthStr := ""
	if period.Month != nil {
		monthStr = strconv.Itoa(*period.Month)
	}

	return usecases.PeriodInfo{
		Year:    yearStr,
		Month:   monthStr,
		Label:   label,
		HasData: hasData,
	}
}
