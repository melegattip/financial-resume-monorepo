package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics/domain"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ---------------------------------------------------------------------------
// Mock repository
// ---------------------------------------------------------------------------

type mockAnalyticsRepo struct {
	dashboard  *domain.DashboardSummary
	categories []domain.CategorySummary
	dashErr    error
}

func (m *mockAnalyticsRepo) GetDashboardSummary(_ context.Context, _ string) (*domain.DashboardSummary, error) {
	return m.dashboard, m.dashErr
}

func (m *mockAnalyticsRepo) GetExpensesByCategory(_ context.Context, _ string, _, _ time.Time) ([]domain.CategorySummary, error) {
	return m.categories, nil
}

func (m *mockAnalyticsRepo) GetExpenseSummary(_ context.Context, _ string, _, _ time.Time, _ string) (*domain.ExpenseSummary, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetIncomeSummary(_ context.Context, _ string, _, _ time.Time, _ string) (*domain.IncomeSummary, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetMonthlyExpenses(_ context.Context, _ string, _ int) ([]domain.MonthlySummary, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetMonthlyIncomes(_ context.Context, _ string, _ int) ([]domain.MonthlySummary, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetTransactionsForReport(_ context.Context, _ string, _, _ time.Time) ([]domain.ReportTransaction, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newHealthHandler(dashboard *domain.DashboardSummary, categories []domain.CategorySummary) *AnalyticsHandler {
	repo := &mockAnalyticsRepo{dashboard: dashboard, categories: categories}
	return NewAnalyticsHandler(repo, zerolog.Nop())
}

func callGetFinancialHealth(h *AnalyticsHandler, query string) *httptest.ResponseRecorder {
	r := gin.New()
	r.GET("/financial-health", h.GetFinancialHealth)

	req := httptest.NewRequest(http.MethodGet, "/financial-health"+query, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseHealthResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	return body
}

// ---------------------------------------------------------------------------
// clamp100
// ---------------------------------------------------------------------------

func TestClamp100_Below(t *testing.T) {
	assert.Equal(t, 0.0, clamp100(-50))
}

func TestClamp100_Above(t *testing.T) {
	assert.Equal(t, 100.0, clamp100(150))
}

func TestClamp100_Within(t *testing.T) {
	assert.Equal(t, 72.5, clamp100(72.5))
}

func TestClamp100_Boundaries(t *testing.T) {
	assert.Equal(t, 0.0, clamp100(0))
	assert.Equal(t, 100.0, clamp100(100))
}

// ---------------------------------------------------------------------------
// parseIntQuery
// ---------------------------------------------------------------------------

func TestParseIntQuery_MissingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	assert.Equal(t, 0, parseIntQuery(c, "missing"))
}

func TestParseIntQuery_ValidValue(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?streak=15", nil)
	assert.Equal(t, 15, parseIntQuery(c, "streak"))
}

func TestParseIntQuery_InvalidValue(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?streak=abc", nil)
	assert.Equal(t, 0, parseIntQuery(c, "streak"))
}

func TestParseIntQuery_NegativeValue(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?streak=-5", nil)
	assert.Equal(t, 0, parseIntQuery(c, "streak"))
}

// ---------------------------------------------------------------------------
// GetFinancialHealth — cash flow dimension
// ---------------------------------------------------------------------------

func TestGetFinancialHealth_NoIncomeNoExpenses(t *testing.T) {
	// cash_flow=100 * 0.4 + planning=0 + consistency=0 + engagement=0 = 40 → score 400
	h := newHealthHandler(&domain.DashboardSummary{}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	score := body["score"].(float64)
	assert.Equal(t, float64(400), score)
	assert.Equal(t, "poor", body["status"]) // finalScore=40 → "poor" (< 50)
}

func TestGetFinancialHealth_NoIncomeWithExpenses(t *testing.T) {
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthExpenses: 1000,
		CurrentMonthIncomes:  0,
	}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	// Score = cash flow(5)*0.4 + rest 0 = 2.0 → display 20
	score := body["score"].(float64)
	assert.Equal(t, float64(20), score)
	assert.Equal(t, "critical", body["status"])
}

func TestGetFinancialHealth_PerfectCashFlow(t *testing.T) {
	// Spend nothing, earn 5000 — pure 100 cash-flow score.
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthExpenses: 0,
		CurrentMonthIncomes:  5000,
		CurrentMonthBalance:  5000,
	}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	cashFlow := body["cash_flow_score"].(float64)
	assert.Equal(t, float64(100), cashFlow)
}

func TestGetFinancialHealth_HighConsumptionRatio(t *testing.T) {
	// Spend 90% of income.
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthExpenses: 900,
		CurrentMonthIncomes:  1000,
		CurrentMonthBalance:  100,
	}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	cashFlow := body["cash_flow_score"].(float64)
	assert.Equal(t, float64(10), cashFlow, "cash flow = 100 - 90 = 10")
}

func TestGetFinancialHealth_SavingsBonusApplied(t *testing.T) {
	// Spend 60% → cash flow 40 normally; balance 40% → savings rate > 20% → +10 bonus.
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthExpenses: 600,
		CurrentMonthIncomes:  1000,
		CurrentMonthBalance:  400, // 40% savings rate
	}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	cashFlow := body["cash_flow_score"].(float64)
	assert.Equal(t, float64(50), cashFlow, "40 + 10 savings bonus = 50")
}

// ---------------------------------------------------------------------------
// GetFinancialHealth — planning dimension
// ---------------------------------------------------------------------------

func TestGetFinancialHealth_PlanningWithBudget(t *testing.T) {
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthIncomes: 1000,
	}, nil)
	w := callGetFinancialHealth(h, "?budgets_created=1")
	body := parseHealthResponse(t, w)

	planning := body["planning_score"].(float64)
	assert.GreaterOrEqual(t, planning, float64(25))
}

func TestGetFinancialHealth_PlanningFullCombo(t *testing.T) {
	// budgets_created=1(+25), savings_goals=1(+20), savings_deposits=2(+10), budget_compliance=1(+10) = 65
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthIncomes: 1000,
	}, nil)
	w := callGetFinancialHealth(h, "?budgets_created=1&savings_goals=1&savings_deposits=2&budget_compliance=1")
	body := parseHealthResponse(t, w)

	planning := body["planning_score"].(float64)
	assert.Equal(t, float64(65), planning)
}

func TestGetFinancialHealth_PlanningCappedAt100(t *testing.T) {
	// Very high params should cap at 100.
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthIncomes: 1000,
	}, nil)
	w := callGetFinancialHealth(h, "?budgets_created=10&savings_goals=10&savings_deposits=20&budget_compliance=20")
	body := parseHealthResponse(t, w)

	planning := body["planning_score"].(float64)
	assert.Equal(t, float64(100), planning)
}

// ---------------------------------------------------------------------------
// GetFinancialHealth — consistency dimension
// ---------------------------------------------------------------------------

func TestGetFinancialHealth_ConsistencyZero(t *testing.T) {
	h := newHealthHandler(&domain.DashboardSummary{}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	consistency := body["consistency_score"].(float64)
	assert.Equal(t, float64(0), consistency)
}

func TestGetFinancialHealth_ConsistencyFull(t *testing.T) {
	// 30-day streak (60) + 90-day tenure (40) = 100.
	h := newHealthHandler(&domain.DashboardSummary{}, nil)
	w := callGetFinancialHealth(h, "?streak=30&days_active=90")
	body := parseHealthResponse(t, w)

	consistency := body["consistency_score"].(float64)
	assert.Equal(t, float64(100), consistency)
}

// ---------------------------------------------------------------------------
// GetFinancialHealth — engagement dimension
// ---------------------------------------------------------------------------

func TestGetFinancialHealth_EngagementZero(t *testing.T) {
	h := newHealthHandler(&domain.DashboardSummary{}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	engagement := body["engagement_score"].(float64)
	assert.Equal(t, float64(0), engagement)
}

func TestGetFinancialHealth_EngagementWithAIAndRecurring(t *testing.T) {
	// 2 AI (40) + 2 recurring (30) = 70.
	h := newHealthHandler(&domain.DashboardSummary{}, nil)
	w := callGetFinancialHealth(h, "?ai_applied=2&recurring_setups=2")
	body := parseHealthResponse(t, w)

	engagement := body["engagement_score"].(float64)
	assert.Equal(t, float64(70), engagement)
}

// ---------------------------------------------------------------------------
// GetFinancialHealth — composite score and status labels
// ---------------------------------------------------------------------------

func TestGetFinancialHealth_StatusExcellent(t *testing.T) {
	// Perfect cash flow (100) + full consistency (100) = 85*0.4 + 0 + 100*0.2 = ...
	// Easier: 0 income+expenses → cash 100, planning 0, consistency 0, engagement 0 → final 40 = "fair"
	// Use full behavioral params for excellent.
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthIncomes:  1000,
		CurrentMonthExpenses: 100,
		CurrentMonthBalance:  900,
	}, nil)
	w := callGetFinancialHealth(h, "?streak=30&days_active=90&budgets_created=1&savings_goals=1&savings_deposits=4&budget_compliance=3&ai_applied=3&recurring_setups=2")
	body := parseHealthResponse(t, w)

	assert.Equal(t, "excellent", body["status"])
}

func TestGetFinancialHealth_StatusCritical(t *testing.T) {
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthIncomes:  1000,
		CurrentMonthExpenses: 980,
		CurrentMonthBalance:  20,
	}, nil)
	w := callGetFinancialHealth(h, "")
	body := parseHealthResponse(t, w)

	assert.Equal(t, "critical", body["status"])
}

func TestGetFinancialHealth_ResponseContainsAllFields(t *testing.T) {
	h := newHealthHandler(&domain.DashboardSummary{
		CurrentMonthIncomes:  2000,
		CurrentMonthExpenses: 1000,
		CurrentMonthBalance:  1000,
	}, nil)
	w := callGetFinancialHealth(h, "?streak=10&days_active=30")
	body := parseHealthResponse(t, w)

	requiredFields := []string{
		"score", "status", "savings_rate",
		"current_month_expenses", "current_month_incomes", "current_month_balance",
		"cash_flow_score", "planning_score", "consistency_score", "engagement_score",
	}
	for _, f := range requiredFields {
		assert.Contains(t, body, f, "response missing field %q", f)
	}
}

func TestGetFinancialHealth_ScaleIs0to1000(t *testing.T) {
	// Run several scenarios and confirm score stays in [0, 1000].
	scenarios := []struct {
		dashboard *domain.DashboardSummary
		query     string
	}{
		{&domain.DashboardSummary{}, ""},
		{&domain.DashboardSummary{CurrentMonthIncomes: 5000, CurrentMonthExpenses: 4999}, ""},
		{&domain.DashboardSummary{CurrentMonthIncomes: 1000, CurrentMonthExpenses: 0, CurrentMonthBalance: 1000}, "?streak=30&days_active=90&budgets_created=5&ai_applied=5"},
	}
	for _, sc := range scenarios {
		h := newHealthHandler(sc.dashboard, nil)
		w := callGetFinancialHealth(h, sc.query)
		body := parseHealthResponse(t, w)
		score := body["score"].(float64)
		assert.GreaterOrEqual(t, score, float64(0))
		assert.LessOrEqual(t, score, float64(1000))
	}
}

// ---------------------------------------------------------------------------
// isProductiveCategory
// ---------------------------------------------------------------------------

func TestIsProductiveCategory_KnownKeywords(t *testing.T) {
	assert.True(t, isProductiveCategory("Inversiones"))
	assert.True(t, isProductiveCategory("Ahorro emergencia"))
	assert.True(t, isProductiveCategory("Seguro de vida"))
	assert.True(t, isProductiveCategory("Educación"))
	assert.True(t, isProductiveCategory("Fondo de retiro"))
	assert.True(t, isProductiveCategory("Bitcoin"))
	assert.True(t, isProductiveCategory("ETF CEDEAR"))
	assert.True(t, isProductiveCategory("Plazo Fijo"))
}

func TestIsProductiveCategory_ConsumptionCategories(t *testing.T) {
	assert.False(t, isProductiveCategory("Supermercado"))
	assert.False(t, isProductiveCategory("Restaurantes"))
	assert.False(t, isProductiveCategory("Transporte"))
	assert.False(t, isProductiveCategory("Entretenimiento"))
}

func TestIsProductiveCategory_CaseInsensitive(t *testing.T) {
	assert.True(t, isProductiveCategory("INVERSIÓN"))
	assert.True(t, isProductiveCategory("ahorro"))
}
