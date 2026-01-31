package analytics

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// ExpensesAnalyticsService implementa ExpensesAnalyticsUseCase
type ExpensesAnalyticsService struct {
	expenseRepo         baseRepo.ExpenseRepository
	incomeRepo          baseRepo.IncomeRepository
	categoryService     usecases.CategoryService
	periodCalculator    usecases.PeriodCalculator
	analyticsCalculator usecases.AnalyticsCalculator
}

// ExpensesAnalytics define la interfaz para analytics de gastos
type ExpensesAnalytics interface {
	GetExpensesSummary(ctx context.Context, request *ExpensesSummaryRequest) (*ExpensesSummaryResponse, error)
}

// NewExpensesAnalyticsService crea una nueva instancia del servicio
func NewExpensesAnalyticsService(
	expenseRepo baseRepo.ExpenseRepository,
	incomeRepo baseRepo.IncomeRepository,
	categoryService usecases.CategoryService,
	periodCalculator usecases.PeriodCalculator,
	analyticsCalculator usecases.AnalyticsCalculator,
) usecases.ExpensesAnalyticsUseCase {
	return &ExpensesAnalyticsService{
		expenseRepo:         expenseRepo,
		incomeRepo:          incomeRepo,
		categoryService:     categoryService,
		periodCalculator:    periodCalculator,
		analyticsCalculator: analyticsCalculator,
	}
}

// GetExpensesSummary implementa el caso de uso principal
func (s *ExpensesAnalyticsService) GetExpensesSummary(ctx context.Context, params usecases.ExpensesSummaryParams) (*usecases.ExpensesSummary, error) {
	// Validar parámetros
	if err := s.validateParams(params); err != nil {
		return nil, err
	}

	// Obtener datos base
	expenses, totalIncome, err := s.getBaseData(ctx, params.UserID, params.Period)
	if err != nil {
		return nil, err
	}

	// Aplicar ordenamiento
	sortedExpenses := s.applySorting(expenses, params.Sorting)

	// Aplicar paginación
	paginatedExpenses, totalCount := s.applyPagination(sortedExpenses, params.Pagination)

	// Convertir a items de respuesta
	expenseItems, err := s.convertToExpenseItems(paginatedExpenses, totalIncome)
	if err != nil {
		return nil, err
	}

	// Calcular resumen
	summary := s.calculateExpenseSummary(expenses, totalIncome)

	// Construir información de paginación
	pagination := s.buildPaginationInfo(totalCount, params.Pagination)

	return &usecases.ExpensesSummary{
		Expenses:   expenseItems,
		Summary:    summary,
		Pagination: pagination,
	}, nil
}

// validateParams valida los parámetros de entrada (SRP)
func (s *ExpensesAnalyticsService) validateParams(params usecases.ExpensesSummaryParams) error {
	if params.UserID == "" {
		return errors.NewBadRequest("El ID del usuario es requerido")
	}

	if params.Period.Year != nil {
		if *params.Period.Year < 1900 || *params.Period.Year > 2100 {
			return errors.NewBadRequest("Año inválido")
		}
	}

	if params.Period.Month != nil {
		if *params.Period.Month < 1 || *params.Period.Month > 12 {
			return errors.NewBadRequest("Mes inválido")
		}
	}

	if params.Pagination.Limit < 0 {
		return errors.NewBadRequest("El límite no puede ser negativo")
	}

	if params.Pagination.Offset < 0 {
		return errors.NewBadRequest("El offset no puede ser negativo")
	}

	return nil
}

// getBaseData obtiene gastos e ingresos base (SRP)
func (s *ExpensesAnalyticsService) getBaseData(ctx context.Context, userID string, period usecases.DatePeriod) ([]*domain.Expense, float64, error) {
	// Obtener gastos
	allExpenses, err := s.expenseRepo.List(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("error obteniendo gastos: %w", err)
	}

	// Filtrar gastos por período
	filteredExpenses := s.filterExpensesByPeriod(allExpenses, period)

	// Obtener ingresos para cálculo de porcentajes
	allIncomes, err := s.incomeRepo.List(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("error obteniendo ingresos: %w", err)
	}

	// Calcular total de ingresos del período
	totalIncome := s.calculateTotalIncomeForPeriod(allIncomes, period)

	return filteredExpenses, totalIncome, nil
}

// filterExpensesByPeriod filtra gastos por período (SRP)
func (s *ExpensesAnalyticsService) filterExpensesByPeriod(expenses []*domain.Expense, period usecases.DatePeriod) []*domain.Expense {
	if period.Year == nil && period.Month == nil {
		return expenses
	}

	var filtered []*domain.Expense
	for _, expense := range expenses {
		if s.matchesExpensePeriod(expense, period) {
			filtered = append(filtered, expense)
		}
	}

	return filtered
}

// matchesExpensePeriod verifica si un gasto coincide con el período (SRP)
func (s *ExpensesAnalyticsService) matchesExpensePeriod(expense *domain.Expense, period usecases.DatePeriod) bool {
	if period.Year != nil && expense.CreatedAt.Year() != *period.Year {
		return false
	}

	if period.Month != nil && int(expense.CreatedAt.Month()) != *period.Month {
		return false
	}

	return true
}

// calculateTotalIncomeForPeriod calcula ingresos totales del período (SRP)
func (s *ExpensesAnalyticsService) calculateTotalIncomeForPeriod(incomes []*domain.Income, period usecases.DatePeriod) float64 {
	var total float64
	for _, income := range incomes {
		if s.matchesIncomePeriod(income, period) {
			total += income.Amount
		}
	}
	return total
}

// matchesIncomePeriod verifica si un ingreso coincide con el período (SRP)
func (s *ExpensesAnalyticsService) matchesIncomePeriod(income *domain.Income, period usecases.DatePeriod) bool {
	if period.Year != nil && income.CreatedAt.Year() != *period.Year {
		return false
	}

	if period.Month != nil && int(income.CreatedAt.Month()) != *period.Month {
		return false
	}

	return true
}

// applySorting aplica criterios de ordenamiento (SRP)
func (s *ExpensesAnalyticsService) applySorting(expenses []*domain.Expense, sorting usecases.SortingCriteria) []*domain.Expense {
	sorted := make([]*domain.Expense, len(expenses))
	copy(sorted, expenses)

	switch sorting.Field {
	case usecases.SortByDate:
		sort.Slice(sorted, func(i, j int) bool {
			if sorting.Order == usecases.Ascending {
				return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
			}
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		})
	case usecases.SortByAmount:
		sort.Slice(sorted, func(i, j int) bool {
			if sorting.Order == usecases.Ascending {
				return sorted[i].Amount < sorted[j].Amount
			}
			return sorted[i].Amount > sorted[j].Amount
		})
	default:
		// Por defecto, ordenar por fecha descendente
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		})
	}

	return sorted
}

// applyPagination aplica paginación (SRP)
func (s *ExpensesAnalyticsService) applyPagination(expenses []*domain.Expense, pagination usecases.PaginationParams) ([]*domain.Expense, int) {
	totalCount := len(expenses)

	if pagination.Limit == 0 {
		pagination.Limit = 50 // Límite por defecto
	}

	start := pagination.Offset
	end := start + pagination.Limit

	if start >= totalCount {
		return []*domain.Expense{}, totalCount
	}

	if end > totalCount {
		end = totalCount
	}

	return expenses[start:end], totalCount
}

// convertToExpenseItems convierte gastos a items de respuesta (SRP)
func (s *ExpensesAnalyticsService) convertToExpenseItems(expenses []*domain.Expense, totalIncome float64) ([]usecases.ExpenseItem, error) {
	// Obtener nombres de categorías
	categoryIDs := s.extractCategoryIDs(expenses)
	categoryNames, err := s.categoryService.GetCategoryNames(categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo nombres de categorías: %w", err)
	}

	var items []usecases.ExpenseItem
	for _, expense := range expenses {
		item := s.buildExpenseItem(expense, totalIncome, categoryNames)
		items = append(items, item)
	}

	return items, nil
}

// extractCategoryIDs extrae IDs únicos de categorías (SRP)
func (s *ExpensesAnalyticsService) extractCategoryIDs(expenses []*domain.Expense) []string {
	categoryMap := make(map[string]bool)
	for _, expense := range expenses {
		if expense.CategoryID != nil && *expense.CategoryID != "" {
			categoryMap[*expense.CategoryID] = true
		}
	}

	var categoryIDs []string
	for id := range categoryMap {
		categoryIDs = append(categoryIDs, id)
	}

	return categoryIDs
}

// buildExpenseItem construye un item de gasto (SRP)
func (s *ExpensesAnalyticsService) buildExpenseItem(expense *domain.Expense, totalIncome float64, categoryNames map[string]string) usecases.ExpenseItem {
	categoryName := "Sin categoría"
	if expense.CategoryID != nil {
		if name, exists := categoryNames[*expense.CategoryID]; exists {
			categoryName = name
		}
	}

	percentageOfIncome := s.analyticsCalculator.CalculatePercentages(expense.Amount, totalIncome)

	// Calcular días hasta vencimiento
	var daysUntilDue *int
	if !expense.DueDate.IsZero() {
		days := int(time.Until(expense.DueDate).Hours() / 24)
		daysUntilDue = &days
	}

	return usecases.ExpenseItem{
		ID:                 expense.ID,
		Description:        expense.Description,
		Amount:             expense.Amount,
		AmountPaid:         expense.AmountPaid,
		PendingAmount:      expense.GetPendingAmount(),
		PercentageOfIncome: percentageOfIncome,
		CategoryID:         s.getCategoryID(expense),
		CategoryName:       categoryName,
		Paid:               expense.Paid,
		DueDate:            s.formatDueDate(expense.DueDate),
		CreatedAt:          expense.CreatedAt.Format("2006-01-02T15:04:05Z"),
		DaysUntilDue:       daysUntilDue,
	}
}

// getCategoryID obtiene el ID de categoría de manera segura (SRP)
func (s *ExpensesAnalyticsService) getCategoryID(expense *domain.Expense) string {
	if expense.CategoryID == nil {
		return ""
	}
	return *expense.CategoryID
}

// formatDueDate formatea la fecha de vencimiento (SRP)
func (s *ExpensesAnalyticsService) formatDueDate(dueDate time.Time) string {
	if dueDate.IsZero() {
		return ""
	}
	return dueDate.Format("2006-01-02")
}

// calculateExpenseSummary calcula el resumen de gastos (SRP)
func (s *ExpensesAnalyticsService) calculateExpenseSummary(expenses []*domain.Expense, totalIncome float64) usecases.ExpenseSummary {
	var totalAmount, paidAmount, pendingAmount float64

	for _, expense := range expenses {
		totalAmount += expense.Amount
		if expense.Paid {
			paidAmount += expense.Amount
		} else {
			pendingAmount += expense.GetPendingAmount()
		}
	}

	averageTransaction := s.analyticsCalculator.CalculateAverage(totalAmount, len(expenses))
	percentageOfTotalIncome := s.analyticsCalculator.CalculatePercentages(totalAmount, totalIncome)

	return usecases.ExpenseSummary{
		TotalAmount:             totalAmount,
		PaidAmount:              paidAmount,
		PendingAmount:           pendingAmount,
		AverageTransaction:      averageTransaction,
		PercentageOfTotalIncome: percentageOfTotalIncome,
	}
}

// buildPaginationInfo construye información de paginación (SRP)
func (s *ExpensesAnalyticsService) buildPaginationInfo(totalCount int, pagination usecases.PaginationParams) usecases.PaginationInfo {
	limit := pagination.Limit
	if limit == 0 {
		limit = 50
	}

	hasMore := pagination.Offset+limit < totalCount

	return usecases.PaginationInfo{
		Total:   totalCount,
		Limit:   limit,
		Offset:  pagination.Offset,
		HasMore: hasMore,
	}
}
