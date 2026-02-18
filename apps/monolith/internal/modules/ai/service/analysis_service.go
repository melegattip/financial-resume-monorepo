package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/domain"
)

// AnalysisService handles AI-powered financial health analysis.
type AnalysisService struct {
	openai *OpenAIClient
}

// NewAnalysisService creates a new AnalysisService.
func NewAnalysisService(openai *OpenAIClient) *AnalysisService {
	return &AnalysisService{openai: openai}
}

// AnalyzeFinancialHealth performs a full financial health analysis using AI.
func (s *AnalysisService) AnalyzeFinancialHealth(ctx context.Context, data domain.FinancialAnalysisData) (*domain.HealthAnalysis, error) {
	systemPrompt := `Eres un asesor financiero experto especializado en análisis de salud financiera.
Tu trabajo es evaluar la situación financiera del usuario y proporcionar un análisis detallado.
Debes generar insights específicos y accionables.
Responde ÚNICAMENTE con un JSON válido en el formato solicitado.`

	userPrompt := s.buildHealthAnalysisPrompt(data)

	raw, err := s.openai.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("error analyzing financial health: %w", err)
	}

	raw = cleanJSONResponse(raw)

	var analysis domain.HealthAnalysis
	if err := json.Unmarshal([]byte(raw), &analysis); err != nil {
		// Return a safe default on parse failure rather than propagating the error.
		return &domain.HealthAnalysis{
			Score:       500,
			Level:       "Regular",
			Message:     "Análisis completado. No se pudo parsear la respuesta detallada.",
			Insights:    []domain.AIInsight{},
			GeneratedAt: time.Now(),
		}, nil
	}

	analysis.GeneratedAt = time.Now()
	return &analysis, nil
}

// GenerateInsights generates personalised financial insights using AI.
func (s *AnalysisService) GenerateInsights(ctx context.Context, data domain.FinancialAnalysisData) ([]domain.AIInsight, error) {
	systemPrompt := `Eres un asesor financiero experto que genera insights personalizados.
Debes proporcionar recomendaciones específicas y accionables basadas en los datos financieros.
Responde ÚNICAMENTE con un JSON válido que contenga un array de insights.`

	userPrompt := s.buildInsightsPrompt(data)

	raw, err := s.openai.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("error generating insights: %w", err)
	}

	raw = cleanJSONResponse(raw)

	var insights []domain.AIInsight
	if err := json.Unmarshal([]byte(raw), &insights); err != nil {
		return []domain.AIInsight{}, nil
	}
	return insights, nil
}

// buildHealthAnalysisPrompt builds the structured prompt for financial health analysis.
func (s *AnalysisService) buildHealthAnalysisPrompt(data domain.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Analiza la salud financiera del usuario:

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
  "score": 0-1000,
  "level": "Excelente|Bueno|Regular|Malo",
  "message": "Resumen de la salud financiera",
  "insights": [
    {
      "title": "Título del insight",
      "description": "Descripción detallada",
      "impact": "Positivo|Negativo|Neutro",
      "score": 0-100,
      "action_type": "maintain|improve|optimize",
      "category": "savings|expenses|income|debt"
    }
  ]
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

// buildInsightsPrompt builds the structured prompt for insight generation.
func (s *AnalysisService) buildInsightsPrompt(data domain.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Genera insights personalizados para el usuario:

Datos financieros:
- Ingresos totales: $%.2f
- Gastos totales: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.2f (0-1)
- Score financiero: %d

Gastos por categoría:
%s

Responde con un array JSON de 3 a 5 insights:
[
  {
    "title": "Título del insight",
    "description": "Descripción detallada y accionable",
    "impact": "Positivo|Negativo|Neutro",
    "score": 0-100,
    "action_type": "maintain|improve|optimize|alert",
    "category": "savings|expenses|income|debt|investment"
  }
]`,
		data.TotalIncome,
		data.TotalExpenses,
		data.SavingsRate*100,
		data.IncomeStability,
		data.FinancialScore,
		formatExpensesByCategory(data.ExpensesByCategory),
	)
}

// cleanJSONResponse strips markdown code fences that some LLMs wrap their JSON in.
func cleanJSONResponse(raw string) string {
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}

// formatExpensesByCategory formats the expenses map as a bulleted list for prompt injection.
func formatExpensesByCategory(expenses map[string]float64) string {
	if len(expenses) == 0 {
		return "No hay datos de gastos por categoría"
	}

	var sb strings.Builder
	for category, amount := range expenses {
		sb.WriteString(fmt.Sprintf("- %s: $%.2f\n", category, amount))
	}
	return sb.String()
}
