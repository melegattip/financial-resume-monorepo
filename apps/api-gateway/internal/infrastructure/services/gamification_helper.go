package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// GamificationHelper ayuda a registrar acciones de gamificación
type GamificationHelper struct {
	serviceURL string
	httpClient *http.Client
}

// NewGamificationHelper crea una nueva instancia del helper
func NewGamificationHelper(serviceURL string) *GamificationHelper {
	return &GamificationHelper{
		serviceURL: serviceURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // Timeout corto para no bloquear
		},
	}
}

// ActionRequest estructura para enviar acciones al microservicio
type ActionRequest struct {
	ActionType  string `json:"action_type"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Description string `json:"description"`
}

// RecordActionAsync registra una acción de forma asíncrona (no bloquea)
func (g *GamificationHelper) RecordActionAsync(userID, actionType, entityType, entityID, description string) {
	go func() {
		err := g.RecordAction(userID, actionType, entityType, entityID, description)
		if err != nil {
			log.Printf("⚠️ Error recording gamification action: %v", err)
		}
	}()
}

// RecordAction registra una acción de gamificación
func (g *GamificationHelper) RecordAction(userID, actionType, entityType, entityID, description string) error {
	action := ActionRequest{
		ActionType:  actionType,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
	}

	jsonData, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("error marshaling action: %w", err)
	}

	url := g.serviceURL + "/gamification/actions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Headers importantes
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID) // Usar el header personalizado

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gamification service returned status %d", resp.StatusCode)
	}

	return nil
}

// IsHealthy verifica si el servicio de gamificación está disponible
func (g *GamificationHelper) IsHealthy() bool {
	// Extract base URL without /api/v1 for health check
	baseURL := strings.TrimSuffix(g.serviceURL, "/api/v1")
	resp, err := g.httpClient.Get(baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Constantes para tipos de acciones
const (
	// IA y insights
	ActionViewInsight       = "view_insight"
	ActionUnderstandInsight = "understand_insight"
	ActionCompleteAction    = "complete_action"
	ActionViewPattern       = "view_pattern"
	ActionUseSuggestion     = "use_suggestion"

	// Navegación y dashboard
	ActionViewDashboard = "view_dashboard"
	ActionViewAnalytics = "view_analytics"

	// Transacciones - Gastos
	ActionViewExpenses  = "view_expenses"
	ActionCreateExpense = "create_expense"
	ActionUpdateExpense = "update_expense"
	ActionDeleteExpense = "delete_expense"

	// Transacciones - Ingresos
	ActionViewIncomes  = "view_incomes"
	ActionCreateIncome = "create_income"
	ActionUpdateIncome = "update_income"
	ActionDeleteIncome = "delete_income"

	// Categorías
	ActionViewCategories = "view_categories"
	ActionCreateCategory = "create_category"
	ActionUpdateCategory = "update_category"
	ActionDeleteCategory = "delete_category"
	ActionAssignCategory = "assign_category"

	// Presupuestos
	ActionViewBudgets  = "view_budgets"
	ActionCreateBudget = "create_budget"
	ActionUpdateBudget = "update_budget"
	ActionDeleteBudget = "delete_budget"
	ActionExceedBudget = "exceed_budget"
	ActionStayOnBudget = "stay_on_budget"

	// Metas de ahorro
	ActionViewSavingsGoals   = "view_savings_goals"
	ActionCreateSavingsGoal  = "create_savings_goal"
	ActionUpdateSavingsGoal  = "update_savings_goal"
	ActionDeleteSavingsGoal  = "delete_savings_goal"
	ActionDepositSavings     = "deposit_savings"
	ActionWithdrawSavings    = "withdraw_savings"
	ActionAchieveSavingsGoal = "achieve_savings_goal"

	// Transacciones recurrentes
	ActionViewRecurring    = "view_recurring"
	ActionCreateRecurring  = "create_recurring"
	ActionUpdateRecurring  = "update_recurring"
	ActionDeleteRecurring  = "delete_recurring"
	ActionExecuteRecurring = "execute_recurring"

	// Engagement y persistencia
	ActionDailyLogin    = "daily_login"
	ActionWeeklyStreak  = "weekly_streak"
	ActionMonthlyStreak = "monthly_streak"
)

// Constantes para tipos de entidades
const (
	// IA y insights
	EntityInsight    = "insight"
	EntitySuggestion = "suggestion"
	EntityPattern    = "pattern"

	// Transacciones financieras
	EntityExpense  = "expense"
	EntityIncome   = "income"
	EntityCategory = "category"

	// Planificación financiera
	EntityBudget      = "budget"
	EntitySavingsGoal = "savings_goal"
	EntityRecurring   = "recurring_transaction"

	// Navegación y analytics
	EntityDashboard = "dashboard"
	EntityAnalytics = "analytics"
	EntityReport    = "report"

	// Engagement
	EntityUser    = "user"
	EntitySession = "session"
	EntityStreak  = "streak"
)

// 🎯 MÉTODOS DE CONVENIENCIA
// Estos métodos facilitan el registro de acciones comunes sin tener que recordar los valores exactos

// RecordExpenseAction registra acciones relacionadas con gastos
func (g *GamificationHelper) RecordExpenseAction(userID, action, expenseID, description string) {
	if g == nil {
		return
	}
	g.RecordActionAsync(userID, action, EntityExpense, expenseID, description)
}

// RecordIncomeAction registra acciones relacionadas con ingresos
func (g *GamificationHelper) RecordIncomeAction(userID, action, incomeID, description string) {
	if g == nil {
		return
	}
	g.RecordActionAsync(userID, action, EntityIncome, incomeID, description)
}

// RecordDashboardView registra cuando el usuario ve el dashboard
func (g *GamificationHelper) RecordDashboardView(userID string) {
	if g == nil {
		return
	}
	g.RecordActionAsync(userID, ActionViewDashboard, EntityDashboard, "main", "Usuario visitó el dashboard principal")
}

// RecordInsightInteraction registra interacciones con insights de IA
func (g *GamificationHelper) RecordInsightInteraction(userID, action, insightID, description string) {
	if g == nil {
		return
	}
	g.RecordActionAsync(userID, action, EntityInsight, insightID, description)
}

// RecordBudgetAction registra acciones relacionadas con presupuestos
func (g *GamificationHelper) RecordBudgetAction(userID, action, budgetID, description string) {
	if g == nil {
		return
	}
	g.RecordActionAsync(userID, action, EntityBudget, budgetID, description)
}

// RecordSavingsGoalAction registra acciones relacionadas con metas de ahorro
func (g *GamificationHelper) RecordSavingsGoalAction(userID, action, goalID, description string) {
	if g == nil {
		return
	}
	g.RecordActionAsync(userID, action, EntitySavingsGoal, goalID, description)
}

// RecordNavigationAction registra acciones de navegación en la app
func (g *GamificationHelper) RecordNavigationAction(userID, section, description string) {
	if g == nil {
		return
	}
	var action string
	var entity string

	switch section {
	case "analytics":
		action = ActionViewAnalytics
		entity = EntityAnalytics
	case "expenses":
		action = ActionViewExpenses
		entity = EntityExpense
	case "incomes":
		action = ActionViewIncomes
		entity = EntityIncome
	case "budgets":
		action = ActionViewBudgets
		entity = EntityBudget
	case "savings":
		action = ActionViewSavingsGoals
		entity = EntitySavingsGoal
	default:
		action = ActionViewDashboard
		entity = EntityDashboard
	}

	g.RecordActionAsync(userID, action, entity, section, description)
}
