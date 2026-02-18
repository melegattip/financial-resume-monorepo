package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/domain"
)

// CreditService handles AI-powered credit analysis and improvement planning.
type CreditService struct {
	openai *OpenAIClient
}

// NewCreditService creates a new CreditService.
func NewCreditService(openai *OpenAIClient) *CreditService {
	return &CreditService{openai: openai}
}

// GenerateCreditPlan generates a personalised credit improvement plan.
func (s *CreditService) GenerateCreditPlan(ctx context.Context, data domain.FinancialAnalysisData) (*domain.CreditPlan, error) {
	systemPrompt := `Eres un asesor financiero experto especializado en análisis crediticio y mejora de score financiero.
Tu trabajo es generar un plan detallado y realista para mejorar la situación crediticia del usuario.
Debes proporcionar acciones específicas, priorizadas y con timeline claro.
Responde ÚNICAMENTE con un JSON válido en el formato solicitado.`

	userPrompt := s.buildCreditPlanPrompt(data)

	raw, err := s.openai.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("error generating credit improvement plan: %w", err)
	}

	raw = cleanJSONResponse(raw)

	var plan domain.CreditPlan
	if err := json.Unmarshal([]byte(raw), &plan); err != nil {
		// Return a safe default plan on parse failure.
		return &domain.CreditPlan{
			CurrentScore:   data.FinancialScore,
			TargetScore:    data.FinancialScore + 100,
			TimelineMonths: 12,
			Actions: []domain.CreditAction{
				{
					Title:       "Incrementar tasa de ahorro",
					Description: "Aumenta tu tasa de ahorro mensual al menos un 5% para mejorar tu score financiero.",
					Priority:    "alta",
					Timeline:    "1-3 meses",
					Impact:      30,
					Difficulty:  "media",
				},
			},
			KeyMetrics:  map[string]interface{}{"emergency_fund_months": 3},
			GeneratedAt: time.Now(),
		}, nil
	}

	plan.GeneratedAt = time.Now()
	return &plan, nil
}

// CalculateCreditScore calculates the user's credit score based on financial data.
func (s *CreditService) CalculateCreditScore(ctx context.Context, data domain.FinancialAnalysisData) (int, error) {
	systemPrompt := `Eres un experto en análisis crediticio que calcula scores financieros.
Tu trabajo es evaluar la situación financiera y asignar un score entre 1 y 1000.
Considera factores como ingresos, gastos, estabilidad y disciplina financiera.
Responde ÚNICAMENTE con un JSON válido que contenga el score calculado.`

	userPrompt := s.buildScoreCalculationPrompt(data)

	raw, err := s.openai.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		return 0, fmt.Errorf("error calculating credit score: %w", err)
	}

	raw = cleanJSONResponse(raw)

	var scoreResponse struct {
		Score int `json:"score"`
	}
	if err := json.Unmarshal([]byte(raw), &scoreResponse); err != nil {
		// Fall back to algorithmic calculation on parse failure.
		return s.calculateDefaultScore(data), nil
	}

	// Clamp to valid range 1-1000.
	if scoreResponse.Score < 1 || scoreResponse.Score > 1000 {
		return s.calculateDefaultScore(data), nil
	}

	return scoreResponse.Score, nil
}

// buildCreditPlanPrompt builds the structured prompt for credit improvement plan generation.
func (s *CreditService) buildCreditPlanPrompt(data domain.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Genera un plan de mejora crediticia para el usuario:

Datos financieros:
- Ingresos totales: $%.2f
- Gastos totales: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.2f (0-1)
- Score financiero actual: %d
- Período: %s

Gastos por categoría:
%s

Responde en JSON con este formato exacto:
{
  "current_score": 1-1000,
  "target_score": 1-1000,
  "timeline_months": 3-24,
  "actions": [
    {
      "title": "Título de la acción",
      "description": "Descripción detallada y accionable",
      "priority": "alta|media|baja",
      "timeline": "1-3 meses|3-6 meses|6-12 meses",
      "impact": 10-50,
      "difficulty": "fácil|media|difícil"
    }
  ],
  "key_metrics": {
    "savings_rate_improvement": 0.05,
    "debt_reduction_target": 0.15,
    "emergency_fund_months": 6
  }
}`,
		data.TotalIncome,
		data.TotalExpenses,
		data.SavingsRate*100,
		data.IncomeStability,
		data.FinancialScore,
		data.Period,
		formatExpensesByCategory(data.ExpensesByCategory),
	)
}

// buildScoreCalculationPrompt builds the structured prompt for credit score calculation.
func (s *CreditService) buildScoreCalculationPrompt(data domain.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Calcula el score crediticio para el usuario:

Datos financieros:
- Ingresos totales: $%.2f
- Gastos totales: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.2f (0-1)
- Score financiero base: %d

Factores a considerar:
- Capacidad de pago (ingresos vs gastos)
- Disciplina financiera (tasa de ahorro)
- Estabilidad de ingresos
- Diversificación de gastos

Responde en JSON:
{
  "score": 1-1000
}`,
		data.TotalIncome,
		data.TotalExpenses,
		data.SavingsRate*100,
		data.IncomeStability,
		data.FinancialScore,
	)
}

// calculateDefaultScore computes a deterministic fallback score when the AI response cannot be parsed.
func (s *CreditService) calculateDefaultScore(data domain.FinancialAnalysisData) int {
	baseScore := 500

	// Savings rate factor: up to +100 points.
	savingsBonus := int(data.SavingsRate * 100)
	if savingsBonus > 100 {
		savingsBonus = 100
	}
	if savingsBonus < 0 {
		savingsBonus = 0
	}

	// Income stability factor: up to +100 points.
	stabilityBonus := int(data.IncomeStability * 100)
	if stabilityBonus > 100 {
		stabilityBonus = 100
	}

	// Payment capacity factor: up to +150 points.
	paymentCapacityBonus := 0
	if data.TotalIncome > 0 {
		ratio := (data.TotalIncome - data.TotalExpenses) / data.TotalIncome
		paymentCapacityBonus = int(ratio * 150)
		if paymentCapacityBonus < 0 {
			paymentCapacityBonus = 0
		}
		if paymentCapacityBonus > 150 {
			paymentCapacityBonus = 150
		}
	}

	finalScore := baseScore + savingsBonus + stabilityBonus + paymentCapacityBonus

	if finalScore < 1 {
		finalScore = 1
	}
	if finalScore > 1000 {
		finalScore = 1000
	}

	return finalScore
}
