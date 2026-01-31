package insights

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// calculateSavingsScore calcula el score basado en la tasa de ahorro
func (s *Service) calculateSavingsScore(savingsRate float64) float64 {
	if savingsRate >= 0.30 { // 30% o más de ahorro
		return 1000
	} else if savingsRate >= 0.20 { // 20-30% de ahorro
		return 800 + (savingsRate-0.20)*2000 // 800-1000
	} else if savingsRate >= 0.10 { // 10-20% de ahorro
		return 600 + (savingsRate-0.10)*2000 // 600-800
	} else if savingsRate >= 0.05 { // 5-10% de ahorro
		return 400 + (savingsRate-0.05)*4000 // 400-600
	} else if savingsRate >= 0 { // 0-5% de ahorro
		return savingsRate * 8000 // 0-400
	} else { // Gastando más de lo que gana
		return math.Max(0, 200+savingsRate*1000) // Penalización
	}
}

// calculateIncomeStabilityScore calcula el score basado en estabilidad de ingresos
func (s *Service) calculateIncomeStabilityScore(stability IncomeStabilityAnalysis) float64 {
	if stability.IsStable {
		return 800 + stability.RecurringIncomeRatio*200 // 800-1000
	}
	return 400 + stability.RecurringIncomeRatio*400 // 400-800
}

// calculateDiversificationScore calcula el score basado en diversificación de gastos
func (s *Service) calculateDiversificationScore(categories []CategoryAnalysis) float64 {
	if len(categories) == 0 {
		return 500 // Neutral si no hay datos
	}

	// Calcular concentración (índice Herfindahl)
	var herfindahl float64
	for _, cat := range categories {
		share := cat.Percentage / 100
		herfindahl += share * share
	}

	// Convertir a score (menos concentración = mejor score)
	// Herfindahl va de 0 (perfectamente diversificado) a 1 (todo en una categoría)
	diversificationScore := (1 - herfindahl) * 1000
	return math.Max(200, diversificationScore) // Mínimo 200 puntos
}

// calculateSpendingControlScore calcula el score basado en control de gastos
func (s *Service) calculateSpendingControlScore(patterns SpendingPatternAnalysis) float64 {
	// Penalizar gastos inusuales
	unusualPenalty := float64(len(patterns.UnusualSpending)) * 50
	baseScore := 800.0

	// Bonificar si el gasto promedio es consistente
	if patterns.AverageTransactionAmount > 0 {
		variation := (patterns.LargestExpense - patterns.SmallestExpense) / patterns.AverageTransactionAmount
		if variation < 2 { // Gastos consistentes
			baseScore += 100
		}
	}

	return math.Max(300, baseScore-unusualPenalty)
}

// analyzeCategoriesSpending analiza el gasto por categorías
func (s *Service) analyzeCategoriesSpending(ctx context.Context, transactions []usecases.Transaction, userID string) ([]CategoryAnalysis, error) {
	// Obtener todas las categorías del usuario
	categories, err := s.categoryRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo categorías: %w", err)
	}

	// Crear mapa de categorías para lookup
	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	// Agrupar gastos por categoría
	categorySpending := make(map[string]*CategoryAnalysis)
	totalExpenses := 0.0

	for _, transaction := range transactions {
		if transaction.GetType() == usecases.ExpenseTransaction {
			categoryID := transaction.GetCategoryID()
			amount := transaction.GetAmount()
			totalExpenses += amount

			if analysis, exists := categorySpending[categoryID]; exists {
				analysis.Amount += amount
				analysis.TransactionCount++
			} else {
				categoryName := categoryMap[categoryID]
				if categoryName == "" {
					categoryName = "Sin categoría"
				}
				categorySpending[categoryID] = &CategoryAnalysis{
					CategoryID:       categoryID,
					CategoryName:     categoryName,
					Amount:           amount,
					TransactionCount: 1,
				}
			}
		}
	}

	// Convertir a slice y calcular porcentajes
	var result []CategoryAnalysis
	for _, analysis := range categorySpending {
		if totalExpenses > 0 {
			analysis.Percentage = (analysis.Amount / totalExpenses) * 100
		}
		analysis.AverageAmount = analysis.Amount / float64(analysis.TransactionCount)
		analysis.IsRecurring = analysis.TransactionCount > 1 // Simplificado
		result = append(result, *analysis)
	}

	return result, nil
}

// analyzeSpendingPatterns analiza patrones de gasto
func (s *Service) analyzeSpendingPatterns(transactions []usecases.Transaction) SpendingPatternAnalysis {
	var expenses []float64
	var totalExpenses float64

	for _, transaction := range transactions {
		if transaction.GetType() == usecases.ExpenseTransaction {
			amount := transaction.GetAmount()
			expenses = append(expenses, amount)
			totalExpenses += amount
		}
	}

	if len(expenses) == 0 {
		return SpendingPatternAnalysis{}
	}

	sort.Float64s(expenses)

	return SpendingPatternAnalysis{
		AverageTransactionAmount: totalExpenses / float64(len(expenses)),
		LargestExpense:           expenses[len(expenses)-1],
		SmallestExpense:          expenses[0],
		FrequentCategories:       []string{},          // TODO: implementar
		UnusualSpending:          []UnusualSpending{}, // TODO: implementar
		DailyAverageSpending:     totalExpenses / 30,  // Simplificado
	}
}

// analyzeIncomeStability analiza la estabilidad de ingresos
func (s *Service) analyzeIncomeStability(transactions []usecases.Transaction) IncomeStabilityAnalysis {
	var incomes []float64
	var totalIncome float64

	for _, transaction := range transactions {
		if transaction.GetType() == usecases.IncomeTransaction {
			amount := transaction.GetAmount()
			incomes = append(incomes, amount)
			totalIncome += amount
		}
	}

	if len(incomes) == 0 {
		return IncomeStabilityAnalysis{
			IsStable:             false,
			AverageMonthlyIncome: 0,
			IncomeVariation:      0,
			RecurringIncomeRatio: 0,
		}
	}

	avgIncome := totalIncome / float64(len(incomes))

	// Calcular variación
	var variance float64
	for _, income := range incomes {
		variance += math.Pow(income-avgIncome, 2)
	}
	variance /= float64(len(incomes))
	variation := math.Sqrt(variance) / avgIncome

	return IncomeStabilityAnalysis{
		IsStable:             variation < 0.2, // Menos del 20% de variación
		AverageMonthlyIncome: avgIncome,
		IncomeVariation:      variation,
		RecurringIncomeRatio: 0.8, // Mock por ahora
	}
}

// analyzeBudgetCompliance analiza el cumplimiento del presupuesto
func (s *Service) analyzeBudgetCompliance(transactions []usecases.Transaction) BudgetComplianceAnalysis {
	// Por ahora mock - en el futuro integrar con sistema de presupuestos
	return BudgetComplianceAnalysis{
		HasBudget:           false,
		BudgetCompliance:    0,
		OverspentCategories: []string{},
	}
}

// calculateSavingsRate calcula la tasa de ahorro
func (s *Service) calculateSavingsRate(totalIncome, totalExpenses float64) float64 {
	if totalIncome <= 0 {
		return 0
	}
	return (totalIncome - totalExpenses) / totalIncome
}

// buildPeriodInfo construye información del período
func (s *Service) buildPeriodInfo(period DatePeriod) PeriodInfo {
	var year, month, label string
	daysInPeriod := 30 // Default

	if period.Year != nil {
		year = strconv.Itoa(*period.Year)
	}
	if period.Month != nil {
		month = strconv.Itoa(*period.Month)
	}

	if period.Year != nil && period.Month != nil {
		label = fmt.Sprintf("%s %d", getMonthName(*period.Month), *period.Year)
		daysInPeriod = getDaysInMonth(*period.Year, *period.Month)
	} else if period.Year != nil {
		label = strconv.Itoa(*period.Year)
		daysInPeriod = 365
	} else {
		label = "Período actual"
	}

	return PeriodInfo{
		Year:         year,
		Month:        month,
		Label:        label,
		DaysInPeriod: daysInPeriod,
	}
}

// Funciones auxiliares para fechas
func getMonthName(month int) string {
	months := []string{
		"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}
	if month >= 1 && month <= 12 {
		return months[month]
	}
	return "Mes"
}

func getDaysInMonth(year, month int) int {
	// Simplificado - en producción usar time.Date
	switch month {
	case 2:
		if isLeapYear(year) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
