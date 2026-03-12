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
	systemPrompt := `Eres un asesor financiero personal de alto valor, especializado en el mercado latinoamericano.

Tu misión: generar exactamente 3 recomendaciones ACCIONABLES que el usuario pueda ejecutar HOY O ESTA SEMANA.

REGLAS ESTRICTAS:
1. Genera EXACTAMENTE 3 insights. Ni más, ni menos.
2. Cada insight debe tener un "next_action": una instrucción específica, concreta e inmediata. Ejemplo: "Transferir $X al fondo Y antes del viernes", NO "considera ahorrar más".
3. DISTINGUE siempre egresos de CONSUMO vs PRODUCTIVOS (inversión, ahorro, activos, seguros, educación, propiedades, cripto, plazo fijo). Los productivos son logros a celebrar.
4. Usa montos EXACTOS ($X) y porcentajes (X%) basados en los datos reales. Nunca generalices.
5. Prioriza por impacto real en el patrimonio: el insight más importante va primero.
6. Contexto latinoamericano: compra de propiedades, dolarización, plazo fijo, fondos comunes, son estrategias positivas.
7. Si la situación financiera es buena, di exactamente qué optimizar o potenciar — no inventes problemas.

Responde ÚNICAMENTE con un JSON array de exactamente 3 objetos. Sin texto adicional.`

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

	surplus := data.TotalIncome - data.TotalExpenses
	savingsPct := data.SavingsRate * 100
	stabilityLabel := "baja (ingresos irregulares o de una sola fuente)"
	if data.IncomeStability >= 0.8 {
		stabilityLabel = "alta (ingresos regulares y diversificados)"
	} else if data.IncomeStability >= 0.5 {
		stabilityLabel = "media (ingresos con algo de variabilidad)"
	}

	surplusPct := 0.0
	if data.TotalIncome > 0 {
		surplusPct = surplus / data.TotalIncome * 100
	}

	sb.WriteString(fmt.Sprintf(`Genera 3 recomendaciones accionables para este usuario.

## FLUJO DE EFECTIVO — PERÍODO: %s
- Ingresos totales: $%.0f
- Egresos totales: $%.0f
- Superávit/Déficit: $%.0f (%.1f%% de los ingresos)
- Tasa de ahorro neta: %.1f%%
- Estabilidad de ingresos: %s

## EGRESOS POR CATEGORÍA (monto y %% sobre ingresos)
%s`,
		data.Period,
		data.TotalIncome,
		data.TotalExpenses,
		surplus,
		surplusPct,
		savingsPct,
		stabilityLabel,
		formatExpensesByCategoryWithPct(data.ExpensesByCategory, data.TotalIncome),
	))

	if len(data.SavingsGoals) > 0 {
		sb.WriteString("\n## METAS DE AHORRO ACTIVAS\n")
		sb.WriteString(formatSavingsGoals(data.SavingsGoals))
	}

	if data.BudgetsSummary != nil && data.BudgetsSummary.TotalBudgets > 0 {
		sb.WriteString("\n## CUMPLIMIENTO DE PRESUPUESTOS\n")
		sb.WriteString(formatBudgetsSummary(data.BudgetsSummary))
	}

	if data.BehaviorProfile != nil {
		sb.WriteString(formatBehaviorProfile(data.BehaviorProfile))
	}

	sb.WriteString(`

## PROCESO:
1. Clasifica los egresos: CONSUMO vs PRODUCTIVOS (inversión/ahorro/activos/seguros/educación/propiedades).
2. Calcula ratio CONSUMO NETO / ingresos.
3. Genera exactamente 3 insights priorizados por impacto en el patrimonio.
4. El campo "next_action" debe ser una instrucción CONCRETA que el usuario puede ejecutar esta semana. Incluye montos o pasos específicos. Máx 120 caracteres.

Responde SOLO con el array JSON de exactamente 3 objetos:
[
  {
    "title": "Título conciso (máx 60 caracteres)",
    "description": "Análisis con datos exactos ($montos, %) y contexto claro. Máx 200 caracteres.",
    "impact": "Positivo|Negativo|Neutro",
    "score": 0-100,
    "action_type": "maintain|improve|optimize|alert|invest",
    "category": "savings|expenses|income|debt|investment|budget|goals",
    "next_action": "Paso concreto e inmediato a tomar esta semana. Incluir $monto si aplica. Máx 120 caracteres."
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

// formatExpensesByCategoryWithPct formats expenses with both absolute amount and % of total income.
func formatExpensesByCategoryWithPct(expenses map[string]float64, totalIncome float64) string {
	if len(expenses) == 0 {
		return "  - Sin datos de egresos por categoría"
	}

	var sb strings.Builder
	for category, amount := range expenses {
		if totalIncome > 0 {
			pct := amount / totalIncome * 100
			sb.WriteString(fmt.Sprintf("  - %s: $%.0f (%.1f%% de ingresos)\n", category, amount, pct))
		} else {
			sb.WriteString(fmt.Sprintf("  - %s: $%.0f\n", category, amount))
		}
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

// formatBehaviorProfile formats the behavioral profile for prompt injection.
func formatBehaviorProfile(b *domain.BehaviorProfileContext) string {
	if b == nil {
		return ""
	}

	reEngagement := ""
	if b.CurrentStreak == 0 {
		reEngagement = "\n  ⚠️ El usuario tiene racha en 0 — incluir una recomendación de re-engagement."
	}

	instruction := "  - Usa el nivel de sofisticación adecuado:"
	switch {
	case b.DisciplineScore < 30 && b.BudgetsCreated == 0:
		instruction += " usuario BÁSICO. Prioriza crear primer presupuesto como next_action."
	case b.DisciplineScore >= 70:
		instruction += " usuario AVANZADO. No enseñar conceptos básicos; ir directo a optimización."
	case b.AIRecommendationsApplied >= 3:
		instruction += " usuario EJECUTOR. Reconocer que actúa sobre recomendaciones y reforzar ese hábito."
	case b.SavingsGoalsAchieved > 0:
		instruction += " usuario con capacidad de ejecución probada. Proponer siguiente meta concreta."
	default:
		instruction += " usuario INTERMEDIO. Equilibrar educación con acción concreta."
	}

	return fmt.Sprintf(`

## PERFIL CONDUCTUAL DEL USUARIO
- Nivel de gamificación: %d (%s)
- Racha activa: %d días consecutivos
- Tiempo en la plataforma: %d días
- Presupuestos configurados: %d (meses respetados: %d)
- Metas de ahorro creadas: %d | Depósitos realizados: %d | Metas completadas: %d
- Transacciones recurrentes configuradas: %d
- Recomendaciones de IA aplicadas anteriormente: %d
- Score de consistencia: %d/100
- Score de disciplina financiera: %d/100
- Score de engagement: %d/100

## INSTRUCCIÓN DE PERSONALIZACIÓN
%s%s`,
		b.CurrentLevel, b.LevelName,
		b.CurrentStreak,
		b.DaysActive,
		b.BudgetsCreated, b.BudgetComplianceEvents,
		b.SavingsGoalsCreated, b.SavingsDeposits, b.SavingsGoalsAchieved,
		b.RecurringSetups,
		b.AIRecommendationsApplied,
		b.ConsistencyScore,
		b.DisciplineScore,
		b.EngagementScore,
		instruction,
		reEngagement,
	)
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
