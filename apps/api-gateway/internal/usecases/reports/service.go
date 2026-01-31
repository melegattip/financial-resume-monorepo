package reports

import (
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/logs"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/logger"
)

type GenerateFinancialReport struct {
	repository ReportRepository
}

func NewGenerateFinancialReport(repository ReportRepository) *GenerateFinancialReport {
	return &GenerateFinancialReport{
		repository: repository,
	}
}

func (s *GenerateFinancialReport) Execute(startDate, endDate time.Time, userID string) (*FinancialReport, error) {
	transactions, err := s.repository.GetTransactions(startDate, endDate, userID)
	if err != nil {
		logger.Error(nil, err, logs.ErrorGeneratingReport.GetMessage(), logs.Tags{
			"start_date": startDate,
			"end_date":   endDate,
			"user_id":    userID,
		})
		return nil, err
	}

	// Obtener nombres de categorías
	categoryIDs := make([]string, 0)
	categoryIDSet := make(map[string]bool)
	for _, t := range transactions {
		if t.CategoryID != "" && !categoryIDSet[t.CategoryID] {
			categoryIDs = append(categoryIDs, t.CategoryID)
			categoryIDSet[t.CategoryID] = true
		}
	}

	categoryNames, err := s.repository.GetCategoryNames(categoryIDs)
	if err != nil {
		logger.Error(nil, err, logs.ErrorGeneratingReport.GetMessage(), logs.Tags{
			"user_id": userID,
		})
		return nil, err
	}

	// Calcular totales
	var totalIncome, totalExpenses float64
	categoryTotals := make(map[string]float64)

	for _, t := range transactions {
		if t.Type == "income" {
			totalIncome += t.Amount
		} else if t.Type == "expense" {
			totalExpenses += t.Amount
			if t.CategoryID != "" {
				categoryTotals[t.CategoryID] += t.Amount
			}
		}
	}

	// Calcular porcentaje de cada transacción respecto al total de ingresos
	for i := range transactions {
		if totalIncome > 0 {
			if transactions[i].Type == "expense" {
				transactions[i].Percentage = (transactions[i].Amount / totalIncome) * 100
			} else {
				// Los ingresos no necesitan porcentaje o pueden ser 100% del total
				transactions[i].Percentage = 0
			}
		} else {
			transactions[i].Percentage = 0
		}
	}

	// Crear resumen por categorías
	var summary []CategorySummary
	for categoryID, amount := range categoryTotals {
		categoryName := categoryNames[categoryID]
		if categoryName == "" {
			categoryName = "Sin categoría"
		}

		percent := 0.0
		if totalExpenses > 0 {
			percent = (amount / totalExpenses) * 100
		}

		summary = append(summary, CategorySummary{
			CategoryID:   categoryID,
			CategoryName: categoryName,
			TotalAmount:  amount,
			Percentage:   percent,
		})
	}

	builder := NewFinancialReportBuilder()
	report := builder.
		WithStartDate(startDate).
		WithEndDate(endDate).
		WithTransactions(transactions).
		WithCategorySummary(summary).
		Build()

	return report, nil
}
