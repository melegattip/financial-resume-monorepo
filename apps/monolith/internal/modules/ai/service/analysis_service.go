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
	systemPrompt := `Eres un asesor financiero senior especializado en el mercado latinoamericano, con experiencia en patrimonio, inversiones y planificación financiera personal.

Tu trabajo es hacer un análisis CONTEXTUAL e INTELIGENTE, NO mecánico. El ratio ingresos/egresos bruto es solo un punto de partida; debes interpretar la NATURALEZA de cada egreso.

PRINCIPIOS FUNDAMENTALES:
1. Egresos en inversión, ahorro, activos, seguros y educación son CONSTRUCCIÓN DE PATRIMONIO — son señal de salud financiera superior, no de gasto.
2. Un usuario que destina el 40% de sus ingresos a inversión/activos y el 50% a consumo está en MEJOR posición que uno que gasta el 80% en consumo puro.
3. El ratio "gasto/ingreso" solo es relevante para gastos de CONSUMO. Los egresos productivos deben analizarse por separado como "inversión de capital".
4. Considera el contexto latinoamericano: compra de propiedades, dolarización de ahorros, plazo fijo, fondos comunes, son movimientos normales y positivos.
5. Sé HONESTO si la situación es preocupante, pero siempre con contexto y propuestas concretas.

Responde ÚNICAMENTE con un JSON válido en el formato solicitado. Sin texto adicional.`

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
	systemPrompt := `Eres un asesor financiero senior especializado en el mercado latinoamericano, con experiencia en patrimonio, inversiones y planificación financiera personal.

Generas insights financieros INTELIGENTES Y CONTEXTUALES, no evaluaciones mecánicas de ratio.

PRINCIPIOS FUNDAMENTALES:
1. DISTINGUE siempre entre egresos de CONSUMO (alimentación, transporte, ocio, servicios) y egresos PRODUCTIVOS (inversión, ahorro, activos, seguros, educación, fondos, propiedades, cripto, plazo fijo).
2. Los egresos productivos son señal de MADUREZ FINANCIERA. Si el usuario invierte/ahorra agresivamente, es un LOGRO a destacar explícitamente.
3. Evalúa el CONSUMO neto (total_expenses menos egresos productivos) al analizar el ratio de gasto.
4. Si el usuario tiene alto ratio total pero porción significativa en inversión/activos, el insight debe celebrarlo y sugerir optimización adicional, no alarmarse.
5. Contexto latinoamericano: compra de propiedades, dolarización, plazo fijo, fondos comunes de inversión, son estrategias válidas y positivas.
6. Usa números EXACTOS en los insights. Sé concreto: "invertiste $X en Y" es mejor que "tienes gastos en inversión".

Responde ÚNICAMENTE con un JSON array válido. Sin texto adicional.`

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
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`Analiza la salud financiera del usuario con datos del período: %s

## FLUJO DE EFECTIVO
- Ingresos totales: $%.2f
- Egresos totales: $%.2f
- Superávit/Déficit: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.2f/1.0
- Score financiero actual: %d/1000

## EGRESOS POR CATEGORÍA
%s`,
		data.Period,
		data.TotalIncome,
		data.TotalExpenses,
		data.TotalIncome-data.TotalExpenses,
		data.SavingsRate*100,
		data.IncomeStability,
		data.FinancialScore,
		formatExpensesByCategory(data.ExpensesByCategory),
	))

	if len(data.SavingsGoals) > 0 {
		sb.WriteString("\n## METAS DE AHORRO ACTIVAS\n")
		sb.WriteString(formatSavingsGoals(data.SavingsGoals))
	}

	if data.BudgetsSummary != nil && data.BudgetsSummary.TotalBudgets > 0 {
		sb.WriteString("\n## CUMPLIMIENTO DE PRESUPUESTOS\n")
		sb.WriteString(formatBudgetsSummary(data.BudgetsSummary))
	}

	sb.WriteString(`

## CRITERIO CLAVE — Categorías que son MOVIMIENTOS POSITIVOS:
Inversión, Ahorro, Fondo de emergencia, Seguros, Educación, Retiro, Pensión, Activos.
Si el usuario tiene egresos en estas categorías, es señal de BUENA salud financiera.

Responde en JSON con este formato exacto:
{
  "score": 0-1000,
  "level": "Excelente|Bueno|Regular|Mejorable",
  "message": "Resumen ejecutivo de la situación financiera (2-3 oraciones)",
  "insights": [
    {
      "title": "Título del insight (máx 60 caracteres)",
      "description": "Análisis con datos específicos y acción concreta",
      "impact": "Positivo|Negativo|Neutro",
      "score": 0-100,
      "action_type": "maintain|improve|optimize|alert|invest",
      "category": "savings|expenses|income|debt|investment|budget|goals"
    }
  ]
}`)

	return sb.String()
}

// buildInsightsPrompt builds the structured prompt for insight generation.
func (s *AnalysisService) buildInsightsPrompt(data domain.FinancialAnalysisData) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`Genera un análisis financiero profundo y personalizado.

## FLUJO DE EFECTIVO — PERÍODO: %s
- Ingresos totales: $%.2f
- Egresos totales: $%.2f
- Superávit/Déficit: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.2f/1.0
- Score financiero: %d/1000

## EGRESOS POR CATEGORÍA (evalúa cada una)
%s`,
		data.Period,
		data.TotalIncome,
		data.TotalExpenses,
		data.TotalIncome-data.TotalExpenses,
		data.SavingsRate*100,
		data.IncomeStability,
		data.FinancialScore,
		formatExpensesByCategory(data.ExpensesByCategory),
	))

	if len(data.SavingsGoals) > 0 {
		sb.WriteString("\n## METAS DE AHORRO ACTIVAS\n")
		sb.WriteString(formatSavingsGoals(data.SavingsGoals))
	}

	if data.BudgetsSummary != nil && data.BudgetsSummary.TotalBudgets > 0 {
		sb.WriteString("\n## CUMPLIMIENTO DE PRESUPUESTOS\n")
		sb.WriteString(formatBudgetsSummary(data.BudgetsSummary))
	}

	sb.WriteString(`

## INSTRUCCIONES PARA EL ANÁLISIS:
1. PRIMERO clasifica los egresos por categoría en: CONSUMO vs PRODUCTIVOS (inversión/ahorro/activos/seguros/educación).
2. Calcula el ratio de CONSUMO NETO (solo gastos de consumo / ingresos). Ese es el indicador real de salud.
3. Si hay egresos productivos significativos → genera un insight destacándolos con datos exactos ($monto, % de ingresos).
4. Identifica las 2-3 categorías de CONSUMO de mayor impacto y evalúa si son optimizables.
5. Si hay metas de ahorro: evalúa progreso real, señala cuáles están en riesgo de no cumplirse.
6. Si hay presupuestos: evalúa cuáles excedidos o en alerta, y qué acción concreta tomar.
7. Usa números exactos ($montos, %) en todas las descripciones. Nunca generalices.
8. Prioriza insights por impacto real en el patrimonio del usuario.

Genera exactamente 4 a 6 insights. Responde SOLO con el array JSON:
[
  {
    "title": "Título conciso (máx 60 caracteres)",
    "description": "Análisis detallado con datos específicos y acción concreta recomendada",
    "impact": "Positivo|Negativo|Neutro",
    "score": 0-100,
    "action_type": "maintain|improve|optimize|alert|invest",
    "category": "savings|expenses|income|debt|investment|budget|goals"
  }
]`)

	return sb.String()
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
		return "  - Sin datos de egresos por categoría"
	}

	var sb strings.Builder
	for category, amount := range expenses {
		sb.WriteString(fmt.Sprintf("  - %s: $%.2f\n", category, amount))
	}
	return sb.String()
}

// formatSavingsGoals formats the savings goals list for prompt injection.
func formatSavingsGoals(goals []domain.SavingsGoalInfo) string {
	if len(goals) == 0 {
		return "  - Sin metas de ahorro activas\n"
	}

	var sb strings.Builder
	for _, g := range goals {
		progressPct := g.Progress * 100
		remaining := g.TargetAmount - g.CurrentAmount
		deadline := ""
		if !g.TargetDate.IsZero() {
			daysLeft := int(time.Until(g.TargetDate).Hours() / 24)
			if daysLeft > 0 {
				deadline = fmt.Sprintf(", vence en %d días", daysLeft)
			} else {
				deadline = " ⚠️ VENCIDA"
			}
		}
		sb.WriteString(fmt.Sprintf("  - %s: $%.0f / $%.0f (%.0f%% completado, falta $%.0f%s)\n",
			g.Name, g.CurrentAmount, g.TargetAmount, progressPct, remaining, deadline))
	}
	return sb.String()
}

// formatBudgetsSummary formats the budgets summary for prompt injection.
func formatBudgetsSummary(b *domain.BudgetsSummaryInfo) string {
	if b == nil || b.TotalBudgets == 0 {
		return "  - Sin presupuestos configurados\n"
	}
	usagePct := 0.0
	if b.TotalAllocated > 0 {
		usagePct = (b.TotalSpent / b.TotalAllocated) * 100
	}
	return fmt.Sprintf(
		"  - Total: %d presupuestos | Asignado: $%.0f | Gastado: $%.0f (%.0f%% del total)\n"+
			"  - En control: %d | En alerta: %d | Excedidos: %d\n",
		b.TotalBudgets, b.TotalAllocated, b.TotalSpent, usagePct,
		b.OnTrackCount, b.WarningCount, b.ExceededCount,
	)
}
