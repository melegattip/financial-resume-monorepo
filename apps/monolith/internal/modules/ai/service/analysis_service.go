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
	systemPrompt := `Eres un asesor financiero experto especializado en el mercado latinoamericano.
Tu trabajo es evaluar la situación financiera del usuario y proporcionar un análisis detallado y accionable.

REGLA FUNDAMENTAL: Los gastos en inversión, ahorro, seguros y educación son MOVIMIENTOS POSITIVOS
que demuestran disciplina financiera y construcción de patrimonio. NUNCA los evalúes negativamente.

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
	systemPrompt := `Eres un asesor financiero experto especializado en el mercado latinoamericano.
Generas insights financieros personalizados, claros y accionables basados en datos reales del usuario.

REGLA FUNDAMENTAL: Los gastos en categorías de inversión, ahorro, fondos, seguros, educación y activos
son MOVIMIENTOS POSITIVOS que indican disciplina financiera y construcción de patrimonio.
Si el usuario gasta en estas categorías, RECONÓCELO como un logro, no como un problema.

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
1. Identifica las 2-3 categorías de mayor gasto y evalúa si son necesarias o reducibles
2. Si hay egresos en Inversión/Ahorro/Fondos/Seguros/Educación/Activos → genera un insight POSITIVO reconociendo ese comportamiento
3. Si hay metas de ahorro: evalúa su progreso, señala cuáles están en riesgo de no cumplirse
4. Si hay presupuestos: evalúa cuáles están excedidos o en alerta, y qué hacer
5. Usa números específicos (montos exactos, porcentajes) en las descripciones
6. Sé directo y concreto con las recomendaciones
7. Prioriza los insights por impacto (los más importantes primero)

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
