package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/financial-ai-service/internal/core/ports"
)

// AnalysisUseCase implementa el caso de uso para análisis financiero
type AnalysisUseCase struct {
	openaiClient ports.OpenAIClient
	cacheClient  ports.CacheClient
}

// NewAnalysisUseCase crea un nuevo caso de uso de análisis
func NewAnalysisUseCase(openaiClient ports.OpenAIClient, cacheClient ports.CacheClient) ports.AIAnalysisPort {
	return &AnalysisUseCase{
		openaiClient: openaiClient,
		cacheClient:  cacheClient,
	}
}

// AnalyzeFinancialHealth analiza la salud financiera del usuario
func (a *AnalysisUseCase) AnalyzeFinancialHealth(ctx context.Context, data ports.FinancialAnalysisData) (*ports.HealthAnalysis, error) {
	log.Printf("🧠 Analyzing financial health for user: %s", data.UserID)

	// Intentar obtener del cache
	cacheKey := fmt.Sprintf("health_analysis:%s:%s", data.UserID, data.Period)
	if cached, err := a.cacheClient.Get(ctx, cacheKey); err == nil {
		log.Printf("📋 Cache hit for health analysis: %s", cacheKey)
		var analysis ports.HealthAnalysis
		if err := json.Unmarshal(cached, &analysis); err == nil {
			return &analysis, nil
		}
	}

	// Construir prompts para IA
	systemPrompt := `Eres un asesor financiero experto especializado en análisis de salud financiera. 
	Tu trabajo es evaluar la situación financiera del usuario y proporcionar un análisis detallado.
	Debes generar insights específicos y accionables.
	Responde ÚNICAMENTE con un JSON válido en el formato solicitado.`

	userPrompt := a.buildHealthAnalysisPrompt(data)

	// Llamar a OpenAI
	response, err := a.openaiClient.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ Error generating health analysis: %v", err)
		return nil, fmt.Errorf("error analyzing financial health: %w", err)
	}

	// Parsear respuesta
	var analysis ports.HealthAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		log.Printf("❌ Error parsing health analysis response: %v", err)
		return nil, fmt.Errorf("error parsing analysis response: %w", err)
	}

	// Agregar timestamp
	analysis.GeneratedAt = time.Now()

	// Guardar en cache (20 horas)
	if cacheData, err := json.Marshal(analysis); err == nil {
		a.cacheClient.Set(ctx, cacheKey, cacheData, 20*time.Hour)
	}

	log.Printf("✅ Health analysis completed for user: %s (Score: %d)", data.UserID, analysis.Score)
	return &analysis, nil
}

// GenerateInsights genera insights financieros personalizados
func (a *AnalysisUseCase) GenerateInsights(ctx context.Context, data ports.FinancialAnalysisData) ([]ports.AIInsight, error) {
	log.Printf("🧠 Generating insights for user: %s", data.UserID)

	// Intentar obtener del cache
	cacheKey := fmt.Sprintf("insights:%s:%s", data.UserID, data.Period)
	if cached, err := a.cacheClient.Get(ctx, cacheKey); err == nil {
		log.Printf("📋 Cache hit for insights: %s", cacheKey)
		var insights []ports.AIInsight
		if err := json.Unmarshal(cached, &insights); err == nil {
			return insights, nil
		}
	}

	// Construir prompts para IA
	systemPrompt := `Eres un asesor financiero experto que genera insights personalizados.
	Debes proporcionar recomendaciones específicas y accionables basadas en los datos financieros.
	Responde ÚNICAMENTE con un JSON válido que contenga un array de insights.`

	userPrompt := a.buildInsightsPrompt(data)

	// Llamar a OpenAI
	response, err := a.openaiClient.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ Error generating insights: %v", err)
		return nil, fmt.Errorf("error generating insights: %w", err)
	}

	// Parsear respuesta
	var insights []ports.AIInsight
	if err := json.Unmarshal([]byte(response), &insights); err != nil {
		log.Printf("❌ Error parsing insights response: %v", err)
		return nil, fmt.Errorf("error parsing insights response: %w", err)
	}

	// Guardar en cache (20 horas)
	if cacheData, err := json.Marshal(insights); err == nil {
		a.cacheClient.Set(ctx, cacheKey, cacheData, 20*time.Hour)
	}

	log.Printf("✅ Generated %d insights for user: %s", len(insights), data.UserID)
	return insights, nil
}

// buildHealthAnalysisPrompt construye el prompt para análisis de salud financiera
func (a *AnalysisUseCase) buildHealthAnalysisPrompt(data ports.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Analiza la salud financiera del usuario:

Datos financieros:
- Ingresos totales: $%.2f
- Gastos totales: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.1f
- Score financiero actual: %d
- Período: %s

Gastos por categoría:
%s

Responde en JSON:
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
		a.formatExpensesByCategory(data.ExpensesByCategory),
	)
}

// buildInsightsPrompt construye el prompt para generar insights
func (a *AnalysisUseCase) buildInsightsPrompt(data ports.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Genera insights personalizados para el usuario:

Datos financieros:
- Ingresos totales: $%.2f
- Gastos totales: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.1f
- Score financiero: %d

Gastos por categoría:
%s

Responde con un array JSON de insights:
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
		a.formatExpensesByCategory(data.ExpensesByCategory),
	)
}

// formatExpensesByCategory formatea las categorías de gastos para el prompt
func (a *AnalysisUseCase) formatExpensesByCategory(expenses map[string]float64) string {
	if len(expenses) == 0 {
		return "No hay datos de gastos por categoría"
	}

	result := ""
	for category, amount := range expenses {
		result += fmt.Sprintf("- %s: $%.2f\n", category, amount)
	}
	return result
}
