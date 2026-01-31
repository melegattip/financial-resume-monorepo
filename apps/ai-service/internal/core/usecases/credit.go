package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/financial-ai-service/internal/core/ports"
)

// CreditUseCase implementa el caso de uso para análisis crediticio
type CreditUseCase struct {
	openaiClient ports.OpenAIClient
	cacheClient  ports.CacheClient
}

// NewCreditUseCase crea un nuevo caso de uso de crédito
func NewCreditUseCase(openaiClient ports.OpenAIClient, cacheClient ports.CacheClient) ports.CreditAnalysisPort {
	return &CreditUseCase{
		openaiClient: openaiClient,
		cacheClient:  cacheClient,
	}
}

// GenerateImprovementPlan genera un plan personalizado de mejora crediticia
func (c *CreditUseCase) GenerateImprovementPlan(ctx context.Context, data ports.FinancialAnalysisData) (*ports.CreditPlan, error) {
	log.Printf("📊 Generating credit improvement plan for user: %s", data.UserID)

	// Intentar obtener del cache
	cacheKey := fmt.Sprintf("credit_plan:%s:%s", data.UserID, data.Period)
	if cached, err := c.cacheClient.Get(ctx, cacheKey); err == nil {
		log.Printf("📋 Cache hit for credit plan: %s", cacheKey)
		var plan ports.CreditPlan
		if err := json.Unmarshal(cached, &plan); err == nil {
			return &plan, nil
		}
	}

	// Construir prompts para IA
	systemPrompt := `Eres un asesor financiero experto especializado en análisis crediticio y mejora de score financiero.
	Tu trabajo es generar un plan detallado y realista para mejorar la situación crediticia del usuario.
	Debes proporcionar acciones específicas, priorizadas y con timeline claro.
	Responde ÚNICAMENTE con un JSON válido en el formato solicitado.`

	userPrompt := c.buildCreditAnalysisPrompt(data)

	// Llamar a OpenAI
	response, err := c.openaiClient.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ Error generating credit improvement plan: %v", err)
		return nil, fmt.Errorf("error generating credit improvement plan: %w", err)
	}

	// Parsear respuesta
	var plan ports.CreditPlan
	if err := json.Unmarshal([]byte(response), &plan); err != nil {
		log.Printf("❌ Error parsing credit plan response: %v", err)
		return nil, fmt.Errorf("error parsing credit plan response: %w", err)
	}

	// Agregar timestamp
	plan.GeneratedAt = time.Now()

	// Guardar en cache (24 horas)
	if cacheData, err := json.Marshal(plan); err == nil {
		c.cacheClient.Set(ctx, cacheKey, cacheData, 24*time.Hour)
	}

	log.Printf("✅ Credit improvement plan generated for user: %s (Current: %d, Target: %d)",
		data.UserID, plan.CurrentScore, plan.TargetScore)
	return &plan, nil
}

// CalculateCreditScore calcula el score crediticio basado en datos financieros
func (c *CreditUseCase) CalculateCreditScore(ctx context.Context, data ports.FinancialAnalysisData) (int, error) {
	log.Printf("🔢 Calculating credit score for user: %s", data.UserID)

	// Intentar obtener del cache
	cacheKey := fmt.Sprintf("credit_score:%s:%s", data.UserID, data.Period)
	if cached, err := c.cacheClient.Get(ctx, cacheKey); err == nil {
		log.Printf("📋 Cache hit for credit score: %s", cacheKey)
		var scoreData struct {
			Score int `json:"score"`
		}
		if err := json.Unmarshal(cached, &scoreData); err == nil {
			return scoreData.Score, nil
		}
	}

	// Construir prompts para IA
	systemPrompt := `Eres un experto en análisis crediticio que calcula scores financieros.
	Tu trabajo es evaluar la situación financiera y asignar un score entre 1-1000.
	Considera factores como ingresos, gastos, estabilidad y disciplina financiera.
	Responde ÚNICAMENTE con un JSON válido que contenga el score calculado.`

	userPrompt := c.buildScoreCalculationPrompt(data)

	// Llamar a OpenAI
	response, err := c.openaiClient.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ Error calculating credit score: %v", err)
		return 0, fmt.Errorf("error calculating credit score: %w", err)
	}

	// Parsear respuesta
	var scoreResponse struct {
		Score int `json:"score"`
	}
	if err := json.Unmarshal([]byte(response), &scoreResponse); err != nil {
		log.Printf("❌ Error parsing credit score response: %v", err)
		return 0, fmt.Errorf("error parsing credit score response: %w", err)
	}

	// Validar score
	if scoreResponse.Score < 1 || scoreResponse.Score > 1000 {
		log.Printf("⚠️ Invalid credit score: %d, using default calculation", scoreResponse.Score)
		scoreResponse.Score = c.calculateDefaultScore(data)
	}

	// Guardar en cache (12 horas)
	if cacheData, err := json.Marshal(scoreResponse); err == nil {
		c.cacheClient.Set(ctx, cacheKey, cacheData, 12*time.Hour)
	}

	log.Printf("✅ Credit score calculated for user: %s (Score: %d)", data.UserID, scoreResponse.Score)
	return scoreResponse.Score, nil
}

// buildCreditAnalysisPrompt construye el prompt para análisis crediticio
func (c *CreditUseCase) buildCreditAnalysisPrompt(data ports.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Genera un plan de mejora crediticia para el usuario:

Datos financieros:
- Ingresos totales: $%.2f
- Gastos totales: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.1f (0-1)
- Score financiero actual: %d
- Período: %s

Gastos por categoría:
%s

Responde en JSON:
{
  "current_score": 1-1000,
  "target_score": 1-1000,
  "timeline_months": 3-24,
  "actions": [
    {
      "title": "Título de la acción",
      "description": "Descripción detallada",
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
		c.formatExpensesByCategory(data.ExpensesByCategory),
	)
}

// buildScoreCalculationPrompt construye el prompt para cálculo de score
func (c *CreditUseCase) buildScoreCalculationPrompt(data ports.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Calcula el score crediticio para el usuario:

Datos financieros:
- Ingresos totales: $%.2f
- Gastos totales: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.1f (0-1)
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

// calculateDefaultScore calcula un score por defecto si la IA falla
func (c *CreditUseCase) calculateDefaultScore(data ports.FinancialAnalysisData) int {
	// Algoritmo simple de scoring
	baseScore := 500

	// Factor de ahorro (0-100 puntos)
	savingsBonus := int(data.SavingsRate * 100)
	if savingsBonus > 100 {
		savingsBonus = 100
	}

	// Factor de estabilidad (0-100 puntos)
	stabilityBonus := int(data.IncomeStability * 100)

	// Factor de capacidad de pago (0-150 puntos)
	if data.TotalIncome > 0 {
		ratio := (data.TotalIncome - data.TotalExpenses) / data.TotalIncome
		paymentCapacityBonus := int(ratio * 150)
		if paymentCapacityBonus < 0 {
			paymentCapacityBonus = 0
		}
		if paymentCapacityBonus > 150 {
			paymentCapacityBonus = 150
		}
		baseScore += paymentCapacityBonus
	}

	finalScore := baseScore + savingsBonus + stabilityBonus

	// Limitar entre 1-1000
	if finalScore < 1 {
		finalScore = 1
	}
	if finalScore > 1000 {
		finalScore = 1000
	}

	return finalScore
}

// formatExpensesByCategory formatea las categorías de gastos para el prompt
func (c *CreditUseCase) formatExpensesByCategory(expenses map[string]float64) string {
	if len(expenses) == 0 {
		return "No hay datos de gastos por categoría"
	}

	result := ""
	for category, amount := range expenses {
		result += fmt.Sprintf("- %s: $%.2f\n", category, amount)
	}
	return result
}
