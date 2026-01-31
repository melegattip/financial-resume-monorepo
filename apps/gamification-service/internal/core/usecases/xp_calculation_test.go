package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/melegattip/financial-gamification-service/testutil"
)

// TestXPCalculationForActions valida que cada acción otorgue el XP correcto según la implementación
func TestXPCalculationForActions(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()

	// Tests basados en la implementación real en calculateXPForAction
	tests := []struct {
		name        string
		actionType  string
		entityType  string
		expectedXP  int
		description string
	}{
		// 🏠 ACCIONES BÁSICAS
		{
			name:        "view_dashboard",
			actionType:  "view_dashboard",
			entityType:  "dashboard",
			expectedXP:  2,
			description: "Ver dashboard debe otorgar 2 XP",
		},
		{
			name:        "view_expenses",
			actionType:  "view_expenses",
			entityType:  "expense",
			expectedXP:  1,
			description: "Ver gastos debe otorgar 1 XP",
		},
		{
			name:        "view_analytics",
			actionType:  "view_analytics",
			entityType:  "analytics",
			expectedXP:  3,
			description: "Ver analytics debe otorgar 3 XP",
		},

		// 💰 TRANSACCIONES (Motor principal de XP)
		{
			name:        "create_expense",
			actionType:  "create_expense",
			entityType:  "expense",
			expectedXP:  8,
			description: "Crear gasto debe otorgar 8 XP",
		},
		{
			name:        "create_income",
			actionType:  "create_income",
			entityType:  "income",
			expectedXP:  8,
			description: "Crear ingreso debe otorgar 8 XP",
		},
		{
			name:        "update_expense",
			actionType:  "update_expense",
			entityType:  "expense",
			expectedXP:  5,
			description: "Actualizar gasto debe otorgar 5 XP",
		},
		{
			name:        "delete_expense",
			actionType:  "delete_expense",
			entityType:  "expense",
			expectedXP:  3,
			description: "Eliminar gasto debe otorgar 3 XP",
		},

		// 🏷️ ORGANIZACIÓN
		{
			name:        "create_category",
			actionType:  "create_category",
			entityType:  "category",
			expectedXP:  10,
			description: "Crear categoría debe otorgar 10 XP",
		},
		{
			name:        "update_category",
			actionType:  "update_category",
			entityType:  "category",
			expectedXP:  5,
			description: "Actualizar categoría debe otorgar 5 XP",
		},
		{
			name:        "assign_category",
			actionType:  "assign_category",
			entityType:  "category",
			expectedXP:  3,
			description: "Asignar categoría debe otorgar 3 XP",
		},

		// 🎯 ENGAGEMENT Y STREAKS
		{
			name:        "daily_login",
			actionType:  "daily_login",
			entityType:  "user",
			expectedXP:  5,
			description: "Login diario debe otorgar 5 XP",
		},
		{
			name:        "weekly_streak",
			actionType:  "weekly_streak",
			entityType:  "streak",
			expectedXP:  25,
			description: "Racha semanal debe otorgar 25 XP",
		},
		{
			name:        "monthly_streak",
			actionType:  "monthly_streak",
			entityType:  "streak",
			expectedXP:  100,
			description: "Racha mensual debe otorgar 100 XP",
		},

		// 🏆 CHALLENGES
		{
			name:        "daily_challenge_complete",
			actionType:  "daily_challenge_complete",
			entityType:  "challenge",
			expectedXP:  20,
			description: "Completar challenge diario debe otorgar 20 XP",
		},
		{
			name:        "weekly_challenge_complete",
			actionType:  "weekly_challenge_complete",
			entityType:  "challenge",
			expectedXP:  75,
			description: "Completar challenge semanal debe otorgar 75 XP",
		},

		// 🔓 FEATURES DESBLOQUEABLES
		{
			name:        "create_savings_goal",
			actionType:  "create_savings_goal",
			entityType:  "goal",
			expectedXP:  15,
			description: "Crear meta de ahorro debe otorgar 15 XP",
		},
		{
			name:        "create_budget",
			actionType:  "create_budget",
			entityType:  "budget",
			expectedXP:  20,
			description: "Crear presupuesto debe otorgar 20 XP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: crear usuario inicial
			userID := "test_user_" + tt.name
			mockRepo.SetupUser(userID, 0, 1)

			// Act: registrar acción
			result, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID,
				ActionType:  tt.actionType,
				EntityType:  tt.entityType,
				EntityID:    "test_entity_" + tt.name,
				Description: tt.description,
			})

			// Assert: verificar resultado
			if err != nil {
				t.Fatalf("No debería haber error al registrar acción: %v", err)
			}

			if result == nil {
				t.Fatal("El resultado no debería ser nil")
			}

			if result.XPEarned != tt.expectedXP {
				t.Errorf("%s: XP ganado debería ser %d pero fue %d",
					tt.description, tt.expectedXP, result.XPEarned)
			}

			if result.TotalXP != tt.expectedXP {
				t.Errorf("%s: XP total debería ser %d después de primera acción pero fue %d",
					tt.description, tt.expectedXP, result.TotalXP)
			}
		})
	}
}

// TestXPAccumulationAcrossMultipleActions valida la acumulación de XP por múltiples acciones
func TestXPAccumulationAcrossMultipleActions(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()
	userID := "test_user_accumulation"

	// Setup: usuario inicial
	mockRepo.SetupUser(userID, 0, 1)

	// Secuencia de acciones con XP esperado
	actionSequence := []struct {
		actionType    string
		entityType    string
		expectedXP    int
		expectedTotal int
		expectedLevel int
		description   string
	}{
		{"view_dashboard", "dashboard", 2, 2, 1, "Primera acción: ver dashboard"},
		{"create_expense", "expense", 8, 10, 1, "Segunda acción: crear gasto"},
		{"create_category", "category", 10, 20, 1, "Tercera acción: crear categoría"},
		{"view_analytics", "analytics", 3, 23, 1, "Cuarta acción: ver analytics"},
		{"update_expense", "expense", 5, 28, 1, "Quinta acción: actualizar gasto"},
		{"daily_login", "user", 5, 33, 1, "Sexta acción: login diario"},
		{"create_expense", "expense", 8, 41, 1, "Séptima acción: otro gasto"},
		{"create_expense", "expense", 8, 49, 1, "Octava acción: otro gasto"},
		{"create_category", "category", 10, 59, 1, "Novena acción: otra categoría"},
		{"create_expense", "expense", 8, 67, 1, "Décima acción: otro gasto"},
		{"create_budget", "budget", 20, 87, 2, "Onceava acción: presupuesto -> NIVEL 2"},
	}

	for i, action := range actionSequence {
		t.Run(fmt.Sprintf("action_%d_%s", i+1, action.actionType), func(t *testing.T) {
			// Act: registrar acción
			result, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID,
				ActionType:  action.actionType,
				EntityType:  action.entityType,
				EntityID:    fmt.Sprintf("entity_%d", i+1),
				Description: action.description,
			})

			// Assert: verificar resultado
			if err != nil {
				t.Fatalf("Error en acción %d: %v", i+1, err)
			}

			if result.XPEarned != action.expectedXP {
				t.Errorf("Acción %d (%s): XP ganado esperado %d, obtenido %d",
					i+1, action.description, action.expectedXP, result.XPEarned)
			}

			if result.TotalXP != action.expectedTotal {
				t.Errorf("Acción %d (%s): XP total esperado %d, obtenido %d",
					i+1, action.description, action.expectedTotal, result.TotalXP)
			}

			if result.NewLevel != action.expectedLevel {
				t.Errorf("Acción %d (%s): Nivel esperado %d, obtenido %d",
					i+1, action.description, action.expectedLevel, result.NewLevel)
			}

			// Verificar detección de level up
			if action.expectedLevel == 2 && i > 0 {
				if !result.LevelUp {
					t.Errorf("Acción %d debería detectar level up", i+1)
				}
			}

			t.Logf("Acción %d: %s -> +%d XP = %d XP total (Nivel %d)",
				i+1, action.description, result.XPEarned, result.TotalXP, result.NewLevel)
		})
	}
}

// TestUserLevelProgression valida la progresión completa de un usuario a través de niveles
func TestUserLevelProgression(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()
	userID := "test_user_progression"

	// Setup: usuario inicial
	mockRepo.SetupUser(userID, 0, 1)

	// Sesiones de usuario simulando uso real
	sessions := []struct {
		name          string
		actions       []string
		expectedMinXP int
		expectedLevel int
		description   string
	}{
		{
			name:          "session_1_onboarding",
			actions:       []string{"view_dashboard", "create_category", "create_expense", "view_analytics"},
			expectedMinXP: 23, // 2 + 10 + 8 + 3 = 23
			expectedLevel: 1,
			description:   "Sesión 1: Onboarding inicial",
		},
		{
			name:          "session_2_daily_use",
			actions:       []string{"daily_login", "create_expense", "create_expense", "update_expense"},
			expectedMinXP: 49, // 23 + 5 + 8 + 8 + 5 = 49
			expectedLevel: 1,
			description:   "Sesión 2: Uso diario",
		},
		{
			name:          "session_3_organization",
			actions:       []string{"create_category", "assign_category", "create_category", "view_analytics"},
			expectedMinXP: 71, // 49 + 10 + 3 + 10 + 3 = 75 -> NIVEL 2
			expectedLevel: 2,
			description:   "Sesión 3: Organización -> Nivel 2",
		},
		{
			name:          "session_4_advanced_features",
			actions:       []string{"create_savings_goal", "create_budget", "weekly_streak"},
			expectedMinXP: 131, // 75 + 15 + 20 + 25 = 135
			expectedLevel: 2,
			description:   "Sesión 4: Features avanzadas",
		},
		{
			name:          "session_5_consistent_use",
			actions:       []string{"daily_login", "create_expense", "create_expense", "create_budget", "create_savings_goal", "weekly_challenge_complete"},
			expectedMinXP: 251, // 135 + 5 + 8 + 8 + 20 + 15 + 75 = 266 -> NIVEL 3
			expectedLevel: 3,
			description:   "Sesión 5: Uso consistente -> Nivel 3",
		},
	}

	cumulativeXP := 0
	for sessionNum, session := range sessions {
		t.Run(session.name, func(t *testing.T) {
			sessionXP := 0

			// Ejecutar todas las acciones de la sesión
			for actionNum, actionType := range session.actions {
				// Determinar entityType basado en actionType
				entityType := getEntityTypeForAction(actionType)

				result, err := service.RecordUserAction(ctx, RecordActionParams{
					UserID:      userID,
					ActionType:  actionType,
					EntityType:  entityType,
					EntityID:    fmt.Sprintf("entity_s%d_a%d", sessionNum+1, actionNum+1),
					Description: fmt.Sprintf("%s - %s", session.description, actionType),
				})

				if err != nil {
					t.Fatalf("Error en sesión %d, acción %s: %v", sessionNum+1, actionType, err)
				}

				sessionXP += result.XPEarned
				cumulativeXP += result.XPEarned
			}

			// Verificar XP mínimo esperado (puede ser mayor debido a multiplicadores)
			if cumulativeXP < session.expectedMinXP {
				t.Errorf("%s: XP acumulado esperado mínimo %d, obtenido %d",
					session.description, session.expectedMinXP, cumulativeXP)
			}

			// Verificar nivel final
			profile, err := service.GetGamificationStats(ctx, userID)
			if err != nil {
				t.Fatalf("Error obteniendo stats: %v", err)
			}

			if profile.CurrentLevel != session.expectedLevel {
				t.Errorf("%s: Nivel esperado %d, obtenido %d",
					session.description, session.expectedLevel, profile.CurrentLevel)
			}

			t.Logf("%s: +%d XP esta sesión = %d XP total (Nivel %d)",
				session.description, sessionXP, profile.TotalXP, profile.CurrentLevel)
		})
	}
}

// Helper function para determinar entityType basado en actionType
func getEntityTypeForAction(actionType string) string {
	entityMap := map[string]string{
		"view_dashboard":            "dashboard",
		"create_expense":            "expense",
		"create_income":             "income",
		"update_expense":            "expense",
		"create_category":           "category",
		"assign_category":           "category",
		"view_analytics":            "analytics",
		"daily_login":               "user",
		"weekly_streak":             "streak",
		"monthly_streak":            "streak",
		"create_savings_goal":       "goal",
		"create_budget":             "budget",
		"daily_challenge_complete":  "challenge",
		"weekly_challenge_complete": "challenge",
	}

	if entityType, exists := entityMap[actionType]; exists {
		return entityType
	}
	return "default"
}
