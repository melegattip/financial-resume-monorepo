package insights

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
)

// generateInsights genera insights personalizados basados en el análisis
func (s *Service) generateInsights(data *AnalyzedFinancialData, healthScore int) []FinancialInsight {
	var insights []FinancialInsight

	// Insight sobre tasa de ahorro
	insights = append(insights, s.generateSavingsInsight(data.SavingsRate, data.TotalIncome))

	// Insight sobre categoría de mayor gasto
	if len(data.ExpenseCategories) > 0 {
		insights = append(insights, s.generateTopCategoryInsight(data.ExpenseCategories))
	}

	// Insight sobre estabilidad de ingresos
	insights = append(insights, s.generateIncomeStabilityInsight(data.IncomeStability))

	// Insight sobre gastos inusuales
	if len(data.SpendingPatterns.UnusualSpending) > 0 {
		insights = append(insights, s.generateUnusualSpendingInsight(data.SpendingPatterns.UnusualSpending))
	}

	// Insight sobre oportunidades de mejora
	insights = append(insights, s.generateOpportunityInsight(data, healthScore))

	return insights
}

// generateSavingsInsight genera insight sobre la tasa de ahorro
func (s *Service) generateSavingsInsight(savingsRate float64, totalIncome float64) FinancialInsight {
	savingsPercent := savingsRate * 100

	var title, description, icon string
	var impact InsightImpact
	var score int
	var action *RecommendedAction

	if savingsRate >= 0.20 {
		title = "¡Excelente capacidad de ahorro!"
		description = fmt.Sprintf("Estás ahorrando %.1f%% de tus ingresos. ¡Sigue así!", savingsPercent)
		icon = "💰"
		impact = InsightImpactGood
		score = 920
	} else if savingsRate >= 0.10 {
		title = "Buen ritmo de ahorro"
		description = fmt.Sprintf("Ahorras %.1f%% de tus ingresos. Considera aumentar gradualmente.", savingsPercent)
		icon = "📈"
		impact = InsightImpactMedium
		score = 750
		action = &RecommendedAction{
			Title:       "Aumentar ahorro mensual",
			Description: "Intenta ahorrar 5% adicional estableciendo una transferencia automática",
			Difficulty:  "Medio",
			XPReward:    150,
		}
	} else {
		title = "Oportunidad de mejorar el ahorro"
		description = fmt.Sprintf("Actualmente ahorras %.1f%%. El objetivo recomendado es al menos 10%%.", savingsPercent)
		icon = "🎯"
		impact = InsightImpactHigh
		score = 450
		action = &RecommendedAction{
			Title:       "Crear plan de ahorro",
			Description: "Establece un objetivo de ahorro automático del 10% de tus ingresos",
			Difficulty:  "Fácil",
			XPReward:    200,
		}
	}

	return FinancialInsight{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Type:        InsightTypeSaving,
		Impact:      impact,
		Score:       score,
		Icon:        icon,
		Action:      action,
		Metadata: map[string]interface{}{
			"savings_rate":   savingsRate,
			"savings_amount": totalIncome * savingsRate,
		},
	}
}

// generateTopCategoryInsight genera insight sobre la categoría de mayor gasto
func (s *Service) generateTopCategoryInsight(categories []CategoryAnalysis) FinancialInsight {
	// Ordenar por cantidad
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Amount > categories[j].Amount
	})

	topCategory := categories[0]

	var impact InsightImpact
	var score int
	var action *RecommendedAction

	if topCategory.Percentage > 40 {
		impact = InsightImpactHigh
		score = 400
		action = &RecommendedAction{
			Title:       "Revisar gastos en " + topCategory.CategoryName,
			Description: "Analiza si puedes reducir gastos en esta categoría que representa mucho porcentaje",
			Difficulty:  "Medio",
			XPReward:    180,
		}
	} else if topCategory.Percentage > 25 {
		impact = InsightImpactMedium
		score = 650
	} else {
		impact = InsightImpactGood
		score = 850
	}

	return FinancialInsight{
		ID:    uuid.New().String(),
		Title: fmt.Sprintf("Tu mayor gasto: %s", topCategory.CategoryName),
		Description: fmt.Sprintf("Gastas $%.2f (%.1f%%) en %s con %d transacciones.",
			topCategory.Amount, topCategory.Percentage, topCategory.CategoryName, topCategory.TransactionCount),
		Type:   InsightTypeSpending,
		Impact: impact,
		Score:  score,
		Icon:   "📊",
		Action: action,
		Metadata: map[string]interface{}{
			"category_id":   topCategory.CategoryID,
			"category_name": topCategory.CategoryName,
			"amount":        topCategory.Amount,
			"percentage":    topCategory.Percentage,
		},
	}
}

// generateIncomeStabilityInsight genera insight sobre estabilidad de ingresos
func (s *Service) generateIncomeStabilityInsight(stability IncomeStabilityAnalysis) FinancialInsight {
	var title, description, icon string
	var impact InsightImpact
	var score int

	if stability.IsStable {
		title = "Ingresos estables"
		description = fmt.Sprintf("Tus ingresos son consistentes con %.1f%% de ingresos recurrentes.", stability.RecurringIncomeRatio*100)
		icon = "✅"
		impact = InsightImpactGood
		score = 850
	} else {
		title = "Ingresos variables"
		description = "Tus ingresos varían significativamente. Considera diversificar fuentes de ingreso."
		icon = "⚠️"
		impact = InsightImpactMedium
		score = 600
	}

	return FinancialInsight{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Type:        InsightTypeIncome,
		Impact:      impact,
		Score:       score,
		Icon:        icon,
		Metadata: map[string]interface{}{
			"is_stable":              stability.IsStable,
			"average_monthly_income": stability.AverageMonthlyIncome,
			"income_variation":       stability.IncomeVariation,
			"recurring_income_ratio": stability.RecurringIncomeRatio,
		},
	}
}

// generateUnusualSpendingInsight genera insight sobre gastos inusuales
func (s *Service) generateUnusualSpendingInsight(unusualSpending []UnusualSpending) FinancialInsight {
	totalUnusual := 0.0
	for _, spending := range unusualSpending {
		totalUnusual += spending.Amount
	}

	return FinancialInsight{
		ID:          uuid.New().String(),
		Title:       "Gastos inusuales detectados",
		Description: fmt.Sprintf("Se detectaron %d gastos inusuales por un total de $%.2f.", len(unusualSpending), totalUnusual),
		Type:        InsightTypePattern,
		Impact:      InsightImpactMedium,
		Score:       500,
		Icon:        "🔍",
		Metadata: map[string]interface{}{
			"unusual_count": len(unusualSpending),
			"total_amount":  totalUnusual,
		},
	}
}

// generateOpportunityInsight genera insight sobre oportunidades de mejora
func (s *Service) generateOpportunityInsight(data *AnalyzedFinancialData, healthScore int) FinancialInsight {
	var title, description, icon string
	var action *RecommendedAction

	if healthScore >= 800 {
		title = "¡Excelente salud financiera!"
		description = "Mantén estos buenos hábitos y considera inversiones para hacer crecer tu dinero."
		icon = "🏆"
		action = &RecommendedAction{
			Title:       "Explorar inversiones",
			Description: "Investiga opciones de inversión para hacer crecer tus ahorros",
			Difficulty:  "Medio",
			XPReward:    300,
		}
	} else if healthScore >= 600 {
		title = "Buena base financiera"
		description = "Tienes una base sólida. Enfócate en optimizar gastos y aumentar ahorros."
		icon = "📈"
		action = &RecommendedAction{
			Title:       "Optimizar presupuesto",
			Description: "Revisa tus gastos mensuales y identifica áreas de mejora",
			Difficulty:  "Fácil",
			XPReward:    200,
		}
	} else {
		title = "Oportunidades de mejora"
		description = "Hay varias áreas donde puedes mejorar tu salud financiera. ¡Comencemos!"
		icon = "🎯"
		action = &RecommendedAction{
			Title:       "Plan de mejora financiera",
			Description: "Crea un plan paso a paso para mejorar tus finanzas",
			Difficulty:  "Fácil",
			XPReward:    250,
		}
	}

	return FinancialInsight{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Type:        InsightTypeOpportunity,
		Impact:      InsightImpactGood,
		Score:       healthScore,
		Icon:        icon,
		Action:      action,
		Metadata: map[string]interface{}{
			"health_score": healthScore,
		},
	}
}
