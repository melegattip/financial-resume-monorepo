package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

// ---------------------------------------------------------------------------
// detectSophistication
// ---------------------------------------------------------------------------

func TestDetectSophistication_Nil(t *testing.T) {
	assert.Equal(t, "BÁSICO", detectSophistication(nil))
}

func TestDetectSophistication_Basic(t *testing.T) {
	b := &domain.BehaviorProfileContext{DisciplineScore: 0}
	assert.Equal(t, "BÁSICO", detectSophistication(b))
}

func TestDetectSophistication_Advanced(t *testing.T) {
	b := &domain.BehaviorProfileContext{DisciplineScore: 70}
	assert.Equal(t, "AVANZADO", detectSophistication(b))
}

func TestDetectSophistication_Executor(t *testing.T) {
	b := &domain.BehaviorProfileContext{DisciplineScore: 50, AIRecommendationsApplied: 3}
	assert.Equal(t, "EJECUTOR", detectSophistication(b))
}

// ---------------------------------------------------------------------------
// buildMonthlyCoachingPrompt
// ---------------------------------------------------------------------------

func TestBuildMonthlyCoachingPrompt_ContainsPreviousMonth(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{TotalIncome: 4000, TotalExpenses: 2500}
	prompt := svc.buildMonthlyCoachingPrompt(data, "2026-02")
	assert.Contains(t, prompt, "2026-02")
	assert.Contains(t, prompt, "4000")
}

func TestBuildMonthlyCoachingPrompt_WithBehaviorProfile(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{
		BehaviorProfile: &domain.BehaviorProfileContext{DisciplineScore: 80},
	}
	prompt := svc.buildMonthlyCoachingPrompt(data, "2026-02")
	assert.Contains(t, prompt, "PERFIL CONDUCTUAL")
}

// ---------------------------------------------------------------------------
// buildEducationCardsPrompt
// ---------------------------------------------------------------------------

func TestBuildEducationCardsPrompt_ContainsFinancialScore(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{FinancialScore: 650, SavingsRate: 0.25}
	prompt := svc.buildEducationCardsPrompt(data)
	assert.Contains(t, prompt, "650")
	assert.Contains(t, prompt, "25.0")
}

func TestBuildEducationCardsPrompt_WithBehaviorProfile(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{
		BehaviorProfile: &domain.BehaviorProfileContext{DisciplineScore: 75},
	}
	prompt := svc.buildEducationCardsPrompt(data)
	assert.Contains(t, prompt, "AVANZADO")
}

// ---------------------------------------------------------------------------
// GenerateMonthlyCoaching — HTTP mock tests
// ---------------------------------------------------------------------------

func makeOpenAIResponse(t *testing.T, content string) []byte {
	t.Helper()
	resp := map[string]interface{}{
		"choices": []map[string]interface{}{
			{"message": map[string]string{"content": content}},
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

func TestGenerateMonthlyCoaching_ValidResponse(t *testing.T) {
	mockContent := `{
      "sentiment": "positivo",
      "summary": "Buen mes",
      "wins": [{"title":"Ahorraste","description":"$500 al fondo"}],
      "improvements": [{"title":"Delivery","description":"Gasto elevado"}],
      "actions": [
        {"title":"Crear presupuesto","detail":"Para delivery","deep_link":"/budgets"},
        {"title":"Depositar","detail":"$500","deep_link":"/savings-goals"},
        {"title":"Revisar","detail":"Categorías","deep_link":"/categories"}
      ],
      "behavior_note": "Consistencia mejoró"
    }`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(makeOpenAIResponse(t, mockContent))
	}))
	defer ts.Close()

	orig := openAIURL
	openAIURL = ts.URL
	defer func() { openAIURL = orig }()

	svc := NewAnalysisService(NewOpenAIClient("test-key"))
	report, err := svc.GenerateMonthlyCoaching(context.Background(), domain.FinancialAnalysisData{}, "2026-02")
	require.NoError(t, err)
	assert.Equal(t, "positivo", report.Sentiment)
	assert.Equal(t, "2026-02", report.Month)
	assert.Len(t, report.Wins, 1)
	assert.Len(t, report.Actions, 3)
	assert.False(t, report.GeneratedAt.IsZero())
}

func TestGenerateMonthlyCoaching_ParseFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(makeOpenAIResponse(t, "not valid json at all"))
	}))
	defer ts.Close()

	orig := openAIURL
	openAIURL = ts.URL
	defer func() { openAIURL = orig }()

	svc := NewAnalysisService(NewOpenAIClient("test-key"))
	report, err := svc.GenerateMonthlyCoaching(context.Background(), domain.FinancialAnalysisData{}, "2026-01")
	require.NoError(t, err) // fallback, no error
	assert.Equal(t, "neutral", report.Sentiment)
	assert.NotEmpty(t, report.Actions) // at least one default action
}

func TestGenerateMonthlyCoaching_OpenAIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"server error"}`))
	}))
	defer ts.Close()

	orig := openAIURL
	openAIURL = ts.URL
	defer func() { openAIURL = orig }()

	svc := NewAnalysisService(NewOpenAIClient("test-key"))
	_, err := svc.GenerateMonthlyCoaching(context.Background(), domain.FinancialAnalysisData{}, "2026-01")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// GenerateEducationCards — HTTP mock tests
// ---------------------------------------------------------------------------

func TestGenerateEducationCards_ValidResponse(t *testing.T) {
	mockContent := `{"cards":[
      {"topic":"ahorro","title":"Fondo","summary":"Ahorra 3 meses","key_concept":"Fondo de emergencia","cta":"Crear meta","deep_link":"/savings-goals","difficulty":"básico"},
      {"topic":"presupuesto","title":"Presupuesto","summary":"Controla gastos","key_concept":"50/30/20","cta":"Crear presupuesto","deep_link":"/budgets","difficulty":"básico"},
      {"topic":"inversión","title":"Invertir","summary":"Haz crecer","key_concept":"Interés compuesto","cta":"Explorar","deep_link":"/insights","difficulty":"avanzado"}
    ]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(makeOpenAIResponse(t, mockContent))
	}))
	defer ts.Close()

	orig := openAIURL
	openAIURL = ts.URL
	defer func() { openAIURL = orig }()

	svc := NewAnalysisService(NewOpenAIClient("test-key"))
	cards, err := svc.GenerateEducationCards(context.Background(), domain.FinancialAnalysisData{})
	require.NoError(t, err)
	assert.Len(t, cards, 3)
	assert.Equal(t, "ahorro", cards[0].Topic)
}

func TestGenerateEducationCards_ParseFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(makeOpenAIResponse(t, "invalid json"))
	}))
	defer ts.Close()

	orig := openAIURL
	openAIURL = ts.URL
	defer func() { openAIURL = orig }()

	svc := NewAnalysisService(NewOpenAIClient("test-key"))
	cards, err := svc.GenerateEducationCards(context.Background(), domain.FinancialAnalysisData{})
	assert.NoError(t, err)
	assert.Empty(t, cards)
}

func TestGenerateEducationCards_OpenAIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	orig := openAIURL
	openAIURL = ts.URL
	defer func() { openAIURL = orig }()

	svc := NewAnalysisService(NewOpenAIClient("test-key"))
	_, err := svc.GenerateEducationCards(context.Background(), domain.FinancialAnalysisData{})
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// isProductiveCoachingCategory
// ---------------------------------------------------------------------------

func TestIsProductiveCoachingCategory_Productive(t *testing.T) {
	cases := []string{
		"Inversiones", "Ahorro emergencia", "Seguro de vida",
		"Educación", "Fondo retiro", "Pension plan",
		"Activos fijos", "Propiedad", "Inmueble",
		"Capital de trabajo", "Patrimonio neto",
		"Cripto", "Bitcoin", "ETF", "Acciones", "Bono",
		"Plazo fijo",
	}
	for _, name := range cases {
		assert.True(t, isProductiveCoachingCategory(name), "expected productive: %s", name)
	}
}

func TestIsProductiveCoachingCategory_NonProductive(t *testing.T) {
	cases := []string{
		"Comida", "Transporte", "Entretenimiento", "Ropa", "Salud",
		"Servicios", "Alquiler", "Restaurantes",
	}
	for _, name := range cases {
		assert.False(t, isProductiveCoachingCategory(name), "expected non-productive: %s", name)
	}
}

func TestBuildMonthlyCoachingPrompt_SeparatesProductiveExpenses(t *testing.T) {
	svc := &AnalysisService{}
	data := domain.FinancialAnalysisData{
		TotalIncome:   10000,
		TotalExpenses: 7000,
		SavingsRate:   0.30,
		ExpensesByCategory: map[string]float64{
			"Comida":      3000,
			"Inversiones": 4000, // productive — should NOT count as consumption
		},
	}
	prompt := svc.buildMonthlyCoachingPrompt(data, "2026-04")

	// The prompt must call out the productive/consumption split.
	assert.Contains(t, prompt, "Construcción patrimonial")
	assert.Contains(t, prompt, "Consumo")
	// True consumption is 3000, so savings rate from income(10000) - consumption(3000) = 70%.
	assert.Contains(t, prompt, "70.0%")
}
