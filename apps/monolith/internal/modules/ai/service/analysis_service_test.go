package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/domain"
)

// ---------------------------------------------------------------------------
// formatBehaviorProfile
// ---------------------------------------------------------------------------

func TestFormatBehaviorProfile_NilReturnsEmpty(t *testing.T) {
	result := formatBehaviorProfile(nil)
	assert.Empty(t, result)
}

func TestFormatBehaviorProfile_ContainsSection(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		CurrentLevel: 3,
		LevelName:    "Intermedio",
		CurrentStreak: 10,
		DaysActive:   30,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "PERFIL CONDUCTUAL")
	assert.Contains(t, result, "INSTRUCCIÓN DE PERSONALIZACIÓN")
}

func TestFormatBehaviorProfile_ContainsLevelAndStreak(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		CurrentLevel:  5,
		LevelName:     "Avanzado",
		CurrentStreak: 20,
		DaysActive:    60,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "5")
	assert.Contains(t, result, "Avanzado")
	assert.Contains(t, result, "20")
	assert.Contains(t, result, "60")
}

func TestFormatBehaviorProfile_ContainsBudgetsAndSavings(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		BudgetsCreated:         3,
		BudgetComplianceEvents: 2,
		SavingsGoalsCreated:    4,
		SavingsDeposits:        7,
		SavingsGoalsAchieved:   1,
		RecurringSetups:        2,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "3")
	assert.Contains(t, result, "2")
	assert.Contains(t, result, "4")
	assert.Contains(t, result, "7")
}

func TestFormatBehaviorProfile_ContainsDimensionScores(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		ConsistencyScore: 75,
		DisciplineScore:  60,
		EngagementScore:  40,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "75")
	assert.Contains(t, result, "60")
	assert.Contains(t, result, "40")
}

// ---------------------------------------------------------------------------
// formatBehaviorProfile — personalization instructions
// ---------------------------------------------------------------------------

func TestFormatBehaviorProfile_BasicUserInstruction(t *testing.T) {
	// Low discipline + no budgets → BÁSICO instruction.
	b := &domain.BehaviorProfileContext{
		DisciplineScore: 10,
		BudgetsCreated:  0,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "BÁSICO")
	assert.Contains(t, result, "presupuesto")
}

func TestFormatBehaviorProfile_AdvancedUserInstruction(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		DisciplineScore: 80,
		BudgetsCreated:  5,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "AVANZADO")
}

func TestFormatBehaviorProfile_ExecutorUserInstruction(t *testing.T) {
	// DisciplineScore < 70, BudgetsCreated > 0, AIRecommendationsApplied >= 3.
	b := &domain.BehaviorProfileContext{
		DisciplineScore:          50,
		BudgetsCreated:           1,
		AIRecommendationsApplied: 3,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "EJECUTOR")
}

func TestFormatBehaviorProfile_SavingsAchieverInstruction(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		DisciplineScore:      50,
		BudgetsCreated:       1,
		SavingsGoalsAchieved: 1,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "ejecución probada")
}

func TestFormatBehaviorProfile_IntermediateDefault(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		DisciplineScore:          40,
		BudgetsCreated:           1,
		AIRecommendationsApplied: 1,
		SavingsGoalsAchieved:     0,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "INTERMEDIO")
}

func TestFormatBehaviorProfile_ReEngagementWarningWhenStreakZero(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		CurrentStreak: 0,
	}
	result := formatBehaviorProfile(b)
	assert.Contains(t, result, "re-engagement")
}

func TestFormatBehaviorProfile_NoReEngagementWarningWhenStreakPositive(t *testing.T) {
	b := &domain.BehaviorProfileContext{
		CurrentStreak: 5,
	}
	result := formatBehaviorProfile(b)
	assert.NotContains(t, result, "re-engagement")
}

// ---------------------------------------------------------------------------
// buildInsightsPrompt — BehaviorProfile injected when present
// ---------------------------------------------------------------------------

func TestBuildInsightsPrompt_WithBehaviorProfile(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{
		Period:      "this_month",
		TotalIncome: 5000,
		BehaviorProfile: &domain.BehaviorProfileContext{
			CurrentLevel:  4,
			LevelName:     "Avanzado",
			DisciplineScore: 80,
		},
	}
	prompt := svc.buildInsightsPrompt(data)
	assert.Contains(t, prompt, "PERFIL CONDUCTUAL")
	assert.Contains(t, prompt, "Avanzado")
}

func TestBuildInsightsPrompt_WithoutBehaviorProfile(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{
		Period:          "this_month",
		TotalIncome:     5000,
		BehaviorProfile: nil,
	}
	prompt := svc.buildInsightsPrompt(data)
	assert.NotContains(t, prompt, "PERFIL CONDUCTUAL")
}

func TestBuildInsightsPrompt_ContainsPeriodAndIncomes(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{
		Period:          "last_month",
		TotalIncome:     3000,
		TotalExpenses:   2000,
		ExpensesByCategory: map[string]float64{"Comida": 800},
	}
	prompt := svc.buildInsightsPrompt(data)
	assert.Contains(t, prompt, "last_month")
	assert.Contains(t, prompt, "3000")
}

// ---------------------------------------------------------------------------
// cleanJSONResponse
// ---------------------------------------------------------------------------

func TestCleanJSONResponse_StripsMarkdownFence(t *testing.T) {
	raw := "```json\n{\"key\":\"val\"}\n```"
	cleaned := cleanJSONResponse(raw)
	assert.Equal(t, `{"key":"val"}`, cleaned)
}

func TestCleanJSONResponse_StripsPlainFence(t *testing.T) {
	raw := "```\n[]\n```"
	cleaned := cleanJSONResponse(raw)
	assert.Equal(t, "[]", cleaned)
}

func TestCleanJSONResponse_NoFenceUnchanged(t *testing.T) {
	raw := `{"score":500}`
	cleaned := cleanJSONResponse(raw)
	assert.Equal(t, raw, cleaned)
}

func TestCleanJSONResponse_TrimsWhitespace(t *testing.T) {
	raw := "  \n{}\n  "
	cleaned := cleanJSONResponse(raw)
	assert.Equal(t, "{}", cleaned)
}

// ---------------------------------------------------------------------------
// formatExpensesByCategory
// ---------------------------------------------------------------------------

func TestFormatExpensesByCategory_Empty(t *testing.T) {
	result := formatExpensesByCategory(nil)
	assert.Contains(t, result, "Sin datos")
}

func TestFormatExpensesByCategory_WithData(t *testing.T) {
	expenses := map[string]float64{"Comida": 500.0}
	result := formatExpensesByCategory(expenses)
	assert.Contains(t, result, "Comida")
	assert.Contains(t, result, "500.00")
}

// ---------------------------------------------------------------------------
// formatExpensesByCategoryWithPct
// ---------------------------------------------------------------------------

func TestFormatExpensesByCategoryWithPct_NoIncome(t *testing.T) {
	expenses := map[string]float64{"Transporte": 200.0}
	result := formatExpensesByCategoryWithPct(expenses, 0)
	assert.Contains(t, result, "Transporte")
	assert.True(t, strings.Contains(result, "200") && !strings.Contains(result, "%"))
}

func TestFormatExpensesByCategoryWithPct_WithIncome(t *testing.T) {
	expenses := map[string]float64{"Comida": 250.0}
	result := formatExpensesByCategoryWithPct(expenses, 1000)
	assert.Contains(t, result, "25.0%")
}
