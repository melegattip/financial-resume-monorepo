package analytics

import (
	"context"
	"fmt"
	"sort"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
)

// CategoriesAnalyticsService implementa CategoriesAnalyticsUseCase
type CategoriesAnalyticsService struct {
	expenseRepo         baseRepo.ExpenseRepository
	incomeRepo          baseRepo.IncomeRepository
	categoryService     usecases.CategoryService
	periodCalculator    usecases.PeriodCalculator
	analyticsCalculator usecases.AnalyticsCalculator
}

// NewCategoriesAnalyticsService crea una nueva instancia del servicio
func NewCategoriesAnalyticsService(
	expenseRepo baseRepo.ExpenseRepository,
	incomeRepo baseRepo.IncomeRepository,
	categoryService usecases.CategoryService,
	periodCalculator usecases.PeriodCalculator,
	analyticsCalculator usecases.AnalyticsCalculator,
) usecases.CategoriesAnalyticsUseCase {
	return &CategoriesAnalyticsService{
		expenseRepo:         expenseRepo,
		incomeRepo:          incomeRepo,
		categoryService:     categoryService,
		periodCalculator:    periodCalculator,
		analyticsCalculator: analyticsCalculator,
	}
}

// GetCategoriesAnalytics implementa el caso de uso principal
func (s *CategoriesAnalyticsService) GetCategoriesAnalytics(ctx context.Context, params usecases.CategoriesAnalyticsParams) (*usecases.CategoriesAnalytics, error) {
	// Validar parámetros
	if err := s.validateParams(params); err != nil {
		return nil, err
	}

	// Obtener datos base
	expenses, incomes, err := s.getBaseData(ctx, params.UserID, params.Period)
	if err != nil {
		return nil, err
	}

	// Calcular totales
	totalExpenses, totalIncome := s.calculateTotals(expenses, incomes)

	// Agrupar transacciones por categoría
	categoryGroups := s.groupTransactionsByCategory(expenses, incomes)

	// Convertir a items de respuesta
	categoryItems, err := s.convertToCategoryItems(params.UserID, categoryGroups, totalExpenses, totalIncome)
	if err != nil {
		return nil, err
	}

	// Ordenar por monto total descendente
	s.sortCategoriesByAmount(categoryItems)

	// Calcular resumen
	summary := s.calculateCategorySummary(categoryItems, totalExpenses)

	return &usecases.CategoriesAnalytics{
		Categories: categoryItems,
		Summary:    summary,
	}, nil
}

// validateParams valida los parámetros de entrada (SRP)
func (s *CategoriesAnalyticsService) validateParams(params usecases.CategoriesAnalyticsParams) error {
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

	return nil
}

// getBaseData obtiene gastos e ingresos base (SRP)
func (s *CategoriesAnalyticsService) getBaseData(ctx context.Context, userID string, period usecases.DatePeriod) ([]*domain.Expense, []*domain.Income, error) {
	// Obtener gastos
	allExpenses, err := s.expenseRepo.List(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("error obteniendo gastos: %w", err)
	}

	// Obtener ingresos
	allIncomes, err := s.incomeRepo.List(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("error obteniendo ingresos: %w", err)
	}

	// Filtrar por período
	filteredExpenses := s.filterExpensesByPeriod(allExpenses, period)
	filteredIncomes := s.filterIncomesByPeriod(allIncomes, period)

	return filteredExpenses, filteredIncomes, nil
}

// filterExpensesByPeriod filtra gastos por período (SRP)
func (s *CategoriesAnalyticsService) filterExpensesByPeriod(expenses []*domain.Expense, period usecases.DatePeriod) []*domain.Expense {
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

// filterIncomesByPeriod filtra ingresos por período (SRP)
func (s *CategoriesAnalyticsService) filterIncomesByPeriod(incomes []*domain.Income, period usecases.DatePeriod) []*domain.Income {
	if period.Year == nil && period.Month == nil {
		return incomes
	}

	var filtered []*domain.Income
	for _, income := range incomes {
		if s.matchesIncomePeriod(income, period) {
			filtered = append(filtered, income)
		}
	}

	return filtered
}

// matchesExpensePeriod verifica si un gasto coincide con el período (SRP)
func (s *CategoriesAnalyticsService) matchesExpensePeriod(expense *domain.Expense, period usecases.DatePeriod) bool {
	if period.Year != nil && expense.CreatedAt.Year() != *period.Year {
		return false
	}

	if period.Month != nil && int(expense.CreatedAt.Month()) != *period.Month {
		return false
	}

	return true
}

// matchesIncomePeriod verifica si un ingreso coincide con el período (SRP)
func (s *CategoriesAnalyticsService) matchesIncomePeriod(income *domain.Income, period usecases.DatePeriod) bool {
	if period.Year != nil && income.CreatedAt.Year() != *period.Year {
		return false
	}

	if period.Month != nil && int(income.CreatedAt.Month()) != *period.Month {
		return false
	}

	return true
}

// calculateTotals calcula totales de gastos e ingresos (SRP)
func (s *CategoriesAnalyticsService) calculateTotals(expenses []*domain.Expense, incomes []*domain.Income) (float64, float64) {
	var totalExpenses, totalIncome float64

	for _, expense := range expenses {
		totalExpenses += expense.Amount
	}

	for _, income := range incomes {
		totalIncome += income.Amount
	}

	return totalExpenses, totalIncome
}

// CategoryGroup representa un grupo de transacciones por categoría
type CategoryGroup struct {
	CategoryID       string
	ExpenseAmount    float64
	IncomeAmount     float64
	TransactionCount int
}

// groupTransactionsByCategory agrupa transacciones por categoría (SRP)
func (s *CategoriesAnalyticsService) groupTransactionsByCategory(expenses []*domain.Expense, incomes []*domain.Income) map[string]*CategoryGroup {
	groups := make(map[string]*CategoryGroup)

	// Procesar gastos
	for _, expense := range expenses {
		categoryID := s.getCategoryID(expense.CategoryID)
		if groups[categoryID] == nil {
			groups[categoryID] = &CategoryGroup{CategoryID: categoryID}
		}
		groups[categoryID].ExpenseAmount += expense.Amount
		groups[categoryID].TransactionCount++
	}

	// Procesar ingresos
	for _, income := range incomes {
		categoryID := s.getCategoryID(income.CategoryID)
		if groups[categoryID] == nil {
			groups[categoryID] = &CategoryGroup{CategoryID: categoryID}
		}
		groups[categoryID].IncomeAmount += income.Amount
		groups[categoryID].TransactionCount++
	}

	return groups
}

// getCategoryID obtiene el ID de categoría de manera segura (SRP)
func (s *CategoriesAnalyticsService) getCategoryID(categoryID *string) string {
	if categoryID == nil || *categoryID == "" {
		return "sin-categoria"
	}
	return *categoryID
}

// convertToCategoryItems convierte grupos a items de respuesta (SRP)
func (s *CategoriesAnalyticsService) convertToCategoryItems(userID string, groups map[string]*CategoryGroup, totalExpenses, totalIncome float64) ([]usecases.CategoryItem, error) {
	// Obtener nombres de categorías
	var categoryIDs []string
	for categoryID := range groups {
		if categoryID != "sin-categoria" {
			categoryIDs = append(categoryIDs, categoryID)
		}
	}

	// Usar el método que consulta la base de datos con userID
	categoryService := s.categoryService.(*services.CategoryServiceImpl)
	categoryNames, err := categoryService.GetCategoryNamesWithUserID(userID, categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo nombres de categorías: %w", err)
	}
	categoryNames["sin-categoria"] = "Sin categoría"

	// Construir items
	var items []usecases.CategoryItem
	for categoryID, group := range groups {
		item := s.buildCategoryItem(group, categoryNames[categoryID], totalExpenses, totalIncome)
		items = append(items, item)
	}

	return items, nil
}

// buildCategoryItem construye un item de categoría (SRP)
func (s *CategoriesAnalyticsService) buildCategoryItem(group *CategoryGroup, categoryName string, totalExpenses, totalIncome float64) usecases.CategoryItem {
	totalAmount := group.ExpenseAmount + group.IncomeAmount

	// Calcular porcentajes
	percentageOfExpenses := s.analyticsCalculator.CalculatePercentages(group.ExpenseAmount, totalExpenses)
	percentageOfIncome := s.analyticsCalculator.CalculatePercentages(totalAmount, totalIncome)

	// Calcular promedio por transacción
	averagePerTransaction := s.analyticsCalculator.CalculateAverage(totalAmount, group.TransactionCount)

	// Generar semilla de color
	colorSeed := s.analyticsCalculator.GenerateColorSeed(group.CategoryID)

	return usecases.CategoryItem{
		CategoryID:            group.CategoryID,
		CategoryName:          categoryName,
		TotalAmount:           totalAmount,
		PercentageOfExpenses:  percentageOfExpenses,
		PercentageOfIncome:    percentageOfIncome,
		TransactionCount:      group.TransactionCount,
		AveragePerTransaction: averagePerTransaction,
		ColorSeed:             colorSeed,
	}
}

// sortCategoriesByAmount ordena categorías por monto total descendente (SRP)
func (s *CategoriesAnalyticsService) sortCategoriesByAmount(items []usecases.CategoryItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].TotalAmount > items[j].TotalAmount
	})
}

// calculateCategorySummary calcula el resumen de categorías (SRP)
func (s *CategoriesAnalyticsService) calculateCategorySummary(items []usecases.CategoryItem, totalAmount float64) usecases.CategorySummary {
	if len(items) == 0 {
		return usecases.CategorySummary{
			TotalCategories:  0,
			LargestCategory:  "",
			SmallestCategory: "",
			TotalAmount:      0,
		}
	}

	// Items ya están ordenados por monto descendente
	largestCategory := items[0].CategoryName
	smallestCategory := items[len(items)-1].CategoryName

	return usecases.CategorySummary{
		TotalCategories:  len(items),
		LargestCategory:  largestCategory,
		SmallestCategory: smallestCategory,
		TotalAmount:      totalAmount,
	}
}
