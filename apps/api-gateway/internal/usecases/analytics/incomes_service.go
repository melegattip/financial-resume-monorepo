package analytics

import (
	"context"
	"fmt"
	"sort"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// IncomesAnalyticsService implementa IncomesAnalyticsUseCase
type IncomesAnalyticsService struct {
	incomeRepo          baseRepo.IncomeRepository
	categoryService     usecases.CategoryService
	periodCalculator    usecases.PeriodCalculator
	analyticsCalculator usecases.AnalyticsCalculator
}

// NewIncomesAnalyticsService crea una nueva instancia del servicio
func NewIncomesAnalyticsService(
	incomeRepo baseRepo.IncomeRepository,
	categoryService usecases.CategoryService,
	periodCalculator usecases.PeriodCalculator,
	analyticsCalculator usecases.AnalyticsCalculator,
) usecases.IncomesAnalyticsUseCase {
	return &IncomesAnalyticsService{
		incomeRepo:          incomeRepo,
		categoryService:     categoryService,
		periodCalculator:    periodCalculator,
		analyticsCalculator: analyticsCalculator,
	}
}

// GetIncomesSummary implementa el caso de uso principal
func (s *IncomesAnalyticsService) GetIncomesSummary(ctx context.Context, params usecases.IncomesSummaryParams) (*usecases.IncomesSummary, error) {
	// Validar parámetros
	if err := s.validateParams(params); err != nil {
		return nil, err
	}

	// Obtener datos base
	incomes, err := s.getBaseData(ctx, params.UserID, params.Period)
	if err != nil {
		return nil, err
	}

	// Aplicar ordenamiento
	sortedIncomes := s.applySorting(incomes, params.Sorting)

	// Convertir a items de respuesta
	incomeItems, err := s.convertToIncomeItems(sortedIncomes)
	if err != nil {
		return nil, err
	}

	// Calcular resumen
	summary := s.calculateIncomeSummary(incomes)

	return &usecases.IncomesSummary{
		Incomes: incomeItems,
		Summary: summary,
	}, nil
}

// validateParams valida los parámetros de entrada (SRP)
func (s *IncomesAnalyticsService) validateParams(params usecases.IncomesSummaryParams) error {
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

// getBaseData obtiene ingresos base (SRP)
func (s *IncomesAnalyticsService) getBaseData(ctx context.Context, userID string, period usecases.DatePeriod) ([]*domain.Income, error) {
	// Obtener ingresos
	allIncomes, err := s.incomeRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ingresos: %w", err)
	}

	// Filtrar por período
	filteredIncomes := s.filterIncomesByPeriod(allIncomes, period)

	return filteredIncomes, nil
}

// filterIncomesByPeriod filtra ingresos por período (SRP)
func (s *IncomesAnalyticsService) filterIncomesByPeriod(incomes []*domain.Income, period usecases.DatePeriod) []*domain.Income {
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

// matchesIncomePeriod verifica si un ingreso coincide con el período (SRP)
func (s *IncomesAnalyticsService) matchesIncomePeriod(income *domain.Income, period usecases.DatePeriod) bool {
	if period.Year != nil && income.CreatedAt.Year() != *period.Year {
		return false
	}

	if period.Month != nil && int(income.CreatedAt.Month()) != *period.Month {
		return false
	}

	return true
}

// applySorting aplica criterios de ordenamiento (SRP)
func (s *IncomesAnalyticsService) applySorting(incomes []*domain.Income, sorting usecases.SortingCriteria) []*domain.Income {
	sorted := make([]*domain.Income, len(incomes))
	copy(sorted, incomes)

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
	case usecases.SortByCategory:
		sort.Slice(sorted, func(i, j int) bool {
			categoryI := s.getCategoryID(sorted[i].CategoryID)
			categoryJ := s.getCategoryID(sorted[j].CategoryID)
			if sorting.Order == usecases.Ascending {
				return categoryI < categoryJ
			}
			return categoryI > categoryJ
		})
	default:
		// Por defecto, ordenar por fecha descendente
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		})
	}

	return sorted
}

// convertToIncomeItems convierte ingresos a items de respuesta (SRP)
func (s *IncomesAnalyticsService) convertToIncomeItems(incomes []*domain.Income) ([]usecases.IncomeItem, error) {
	// Obtener nombres de categorías
	categoryIDs := s.extractCategoryIDs(incomes)
	categoryNames, err := s.categoryService.GetCategoryNames(categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo nombres de categorías: %w", err)
	}

	var items []usecases.IncomeItem
	for _, income := range incomes {
		item := s.buildIncomeItem(income, categoryNames)
		items = append(items, item)
	}

	return items, nil
}

// extractCategoryIDs extrae IDs únicos de categorías (SRP)
func (s *IncomesAnalyticsService) extractCategoryIDs(incomes []*domain.Income) []string {
	categoryMap := make(map[string]bool)
	for _, income := range incomes {
		if income.CategoryID != nil && *income.CategoryID != "" {
			categoryMap[*income.CategoryID] = true
		}
	}

	var categoryIDs []string
	for id := range categoryMap {
		categoryIDs = append(categoryIDs, id)
	}

	return categoryIDs
}

// buildIncomeItem construye un item de ingreso (SRP)
func (s *IncomesAnalyticsService) buildIncomeItem(income *domain.Income, categoryNames map[string]string) usecases.IncomeItem {
	categoryName := "Sin categoría"
	if income.CategoryID != nil {
		if name, exists := categoryNames[*income.CategoryID]; exists {
			categoryName = name
		}
	}

	return usecases.IncomeItem{
		ID:           income.ID,
		Description:  income.Description,
		Amount:       income.Amount,
		CategoryID:   s.getCategoryID(income.CategoryID),
		CategoryName: categoryName,
		CreatedAt:    income.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// getCategoryID obtiene el ID de categoría de manera segura (SRP)
func (s *IncomesAnalyticsService) getCategoryID(categoryID *string) string {
	if categoryID == nil {
		return ""
	}
	return *categoryID
}

// calculateIncomeSummary calcula el resumen de ingresos (SRP)
func (s *IncomesAnalyticsService) calculateIncomeSummary(incomes []*domain.Income) usecases.IncomeSummary {
	var totalAmount float64
	count := len(incomes)

	for _, income := range incomes {
		totalAmount += income.Amount
	}

	averageTransaction := s.analyticsCalculator.CalculateAverage(totalAmount, count)

	return usecases.IncomeSummary{
		TotalAmount:        totalAmount,
		AverageTransaction: averageTransaction,
		TransactionCount:   count,
	}
}
