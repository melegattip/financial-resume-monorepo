package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/melegattip/financial-gamification-service/testutil"
)

// TestAchievementProgressionWithXP valida el progreso de achievements y su contribución de XP
func TestAchievementProgressionWithXP(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()
	userID := "test_user_achievement_progression"

	// Setup: usuario inicial
	mockRepo.SetupUser(userID, 0, 1)

	// Test: view_insight no otorga XP ni progresa achievements
	t.Run("ai_partner_achievement_no_progress_on_view", func(t *testing.T) {
		totalXPBefore := 0
		// Simular 10 view_insight, que ahora no deben otorgar XP ni registrar acción
		for i := 0; i < 10; i++ {
			result, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID,
				ActionType:  "view_insight",
				EntityType:  "insight",
				EntityID:    fmt.Sprintf("insight_%d", i+1),
				Description: fmt.Sprintf("Viewing insight %d", i+1),
			})
			if err != nil {
				t.Fatalf("Error en insight %d: %v", i+1, err)
			}
			if result.XPEarned != 0 {
				t.Errorf("view_insight debería otorgar 0 XP, obtuvo %d", result.XPEarned)
			}
			if result.TotalXP != totalXPBefore {
				t.Errorf("XP total no debe cambiar con view_insight. Esperado %d, obtenido %d", totalXPBefore, result.TotalXP)
			}
		}

		// Verificar achievements: AI Partner no progresa por views
		achievements, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("Error obteniendo achievements: %v", err)
		}
		aiAchievement := testutil.FindAchievementByType(achievements, "ai_partner")
		if aiAchievement == nil {
			t.Fatal("AI Partner achievement no encontrado")
		}
		if aiAchievement.Progress != 0 || aiAchievement.Completed {
			t.Errorf("AI Partner no debe progresar con view_insight. Progreso=%d, Completado=%v", aiAchievement.Progress, aiAchievement.Completed)
		}
	})

	// Reset usuario para siguiente test
	mockRepo.SetupUser(userID+"_2", 0, 1)

	// Test progresión de Quick Learner achievement
	t.Run("quick_learner_achievement_progression", func(t *testing.T) {
		userID2 := userID + "_2"

		// Simular 5 understand_insight para completar Quick Learner
		for i := 0; i < 5; i++ {
			result, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID2,
				ActionType:  "understand_insight",
				EntityType:  "insight",
				EntityID:    fmt.Sprintf("insight_understand_%d", i+1),
				Description: fmt.Sprintf("Understanding insight %d", i+1),
			})

			if err != nil {
				t.Fatalf("Error en understand insight %d: %v", i+1, err)
			}

			// Verificar XP acumulado (understand_insight da 15 XP según implementación)
			expectedXP := (i + 1) * 15 // 15 XP por understand_insight
			if result.TotalXP != expectedXP {
				t.Errorf("Understand %d: XP total esperado %d, obtenido %d",
					i+1, expectedXP, result.TotalXP)
			}

			t.Logf("Understand %d: +15 XP = %d XP total", i+1, result.TotalXP)
		}

		// Verificar Quick Learner achievement
		achievements, err := service.GetUserAchievements(ctx, userID2)
		if err != nil {
			t.Fatalf("Error obteniendo achievements: %v", err)
		}

		quickAchievement := testutil.FindAchievementByType(achievements, "quick_learner")
		if quickAchievement == nil {
			t.Fatal("Quick Learner achievement no encontrado")
		}

		if quickAchievement.Progress != 5 {
			t.Errorf("Quick Learner debería tener progreso 5/5, tiene %d", quickAchievement.Progress)
		}

		if !quickAchievement.Completed {
			t.Error("Quick Learner debería estar completado")
		}
	})
}

// TestMultipleAchievementsSimultaneousProgression valida progreso simultáneo de múltiples achievements
func TestMultipleAchievementsSimultaneousProgression(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()
	userID := "test_user_multiple_achievements"

	// Setup: usuario inicial
	mockRepo.SetupUser(userID, 0, 1)

	// Secuencia de acciones que afectan múltiples achievements (actualizada: view_insight = 0 XP, sin progreso AI)
	actionSequence := []struct {
		actionType     string
		entityType     string
		expectedXP     int
		aiProgress     int
		quickProgress  int
		budgetProgress int
		description    string
	}{
		{"view_insight", "insight", 0, 0, 0, 0, "Primera vista de insight"},         // view_insight = 0 XP y no progresa AI
		{"understand_insight", "insight", 15, 0, 1, 0, "Primer insight entendido"},  // understand_insight = 15 XP
		{"view_insight", "insight", 0, 0, 1, 0, "Segunda vista de insight"},         // view_insight = 0 XP
		{"create_budget", "budget", 20, 0, 1, 1, "Primer presupuesto"},              // create_budget = 20 XP (no afecta AI Partner)
		{"understand_insight", "insight", 15, 0, 2, 1, "Segundo insight entendido"}, // understand_insight = 15 XP
		{"view_insight", "insight", 0, 0, 2, 1, "Tercera vista de insight"},         // view_insight = 0 XP
		{"create_budget", "budget", 20, 0, 2, 2, "Segundo presupuesto"},             // create_budget = 20 XP (no afecta AI Partner)
	}

	cumulativeXP := 0
	for i, action := range actionSequence {
		t.Run(fmt.Sprintf("action_%d_%s", i+1, action.actionType), func(t *testing.T) {
			result, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID,
				ActionType:  action.actionType,
				EntityType:  action.entityType,
				EntityID:    fmt.Sprintf("entity_%d", i+1),
				Description: action.description,
			})

			if err != nil {
				t.Fatalf("Error en acción %d: %v", i+1, err)
			}

			// Verificar XP de la acción
			if result.XPEarned != action.expectedXP {
				t.Errorf("Acción %d: XP esperado %d, obtenido %d",
					i+1, action.expectedXP, result.XPEarned)
			}

			cumulativeXP += action.expectedXP
			if result.TotalXP != cumulativeXP {
				t.Errorf("Acción %d: XP total esperado %d, obtenido %d",
					i+1, cumulativeXP, result.TotalXP)
			}

			// Verificar progreso de achievements
			achievements, err := service.GetUserAchievements(ctx, userID)
			if err != nil {
				t.Fatalf("Error obteniendo achievements en acción %d: %v", i+1, err)
			}

			// Verificar AI Partner
			aiAchievement := testutil.FindAchievementByType(achievements, "ai_partner")
			if aiAchievement != nil && aiAchievement.Progress != action.aiProgress {
				t.Errorf("Acción %d: AI Partner progress esperado %d, obtenido %d",
					i+1, action.aiProgress, aiAchievement.Progress)
			}

			// Verificar Quick Learner
			quickAchievement := testutil.FindAchievementByType(achievements, "quick_learner")
			if quickAchievement != nil && quickAchievement.Progress != action.quickProgress {
				t.Errorf("Acción %d: Quick Learner progress esperado %d, obtenido %d",
					i+1, action.quickProgress, quickAchievement.Progress)
			}

			// Verificar Budget Master
			budgetAchievement := testutil.FindAchievementByType(achievements, "budget_master")
			if budgetAchievement != nil && budgetAchievement.Progress != action.budgetProgress {
				t.Errorf("Acción %d: Budget Master progress esperado %d, obtenido %d",
					i+1, action.budgetProgress, budgetAchievement.Progress)
			}

			t.Logf("Acción %d (%s): +%d XP = %d XP total | AI: %d, Quick: %d, Budget: %d",
				i+1, action.description, result.XPEarned, result.TotalXP,
				action.aiProgress, action.quickProgress, action.budgetProgress)
		})
	}
}

// TestAchievementCompletionBonusXP valida XP bonus al completar achievements
func TestAchievementCompletionBonusXP(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()
	userID := "test_user_completion_bonus"

	// Setup: usuario inicial
	mockRepo.SetupUser(userID, 0, 1)

	// Completar Quick Learner achievement (5 insights entendidos)
	t.Run("complete_quick_learner_for_bonus", func(t *testing.T) {
		baseXP := 0

		// Hacer 4 understand_insight (casi completar)
		for i := 0; i < 4; i++ {
			result, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID,
				ActionType:  "understand_insight",
				EntityType:  "insight",
				EntityID:    fmt.Sprintf("insight_%d", i+1),
				Description: fmt.Sprintf("Understanding insight %d", i+1),
			})

			if err != nil {
				t.Fatalf("Error en insight %d: %v", i+1, err)
			}

			baseXP += result.XPEarned
		}

		// Verificar que está casi completo
		achievements, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("Error obteniendo achievements: %v", err)
		}

		quickAchievement := testutil.FindAchievementByType(achievements, "quick_learner")
		if quickAchievement == nil {
			t.Fatal("Quick Learner achievement no encontrado")
		}

		if quickAchievement.Progress != 4 {
			t.Errorf("Quick Learner debería tener progreso 4/5, tiene %d", quickAchievement.Progress)
		}

		if quickAchievement.Completed {
			t.Error("Quick Learner NO debería estar completado aún")
		}

		// Hacer el quinto understand_insight que debería completar el achievement
		result, err := service.RecordUserAction(ctx, RecordActionParams{
			UserID:      userID,
			ActionType:  "understand_insight",
			EntityType:  "insight",
			EntityID:    "insight_final",
			Description: "Final insight to complete achievement",
		})

		if err != nil {
			t.Fatalf("Error en insight final: %v", err)
		}

		// Verificar que el achievement está completado
		achievements, err = service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("Error obteniendo achievements después de completar: %v", err)
		}

		quickAchievement = testutil.FindAchievementByType(achievements, "quick_learner")
		if quickAchievement == nil {
			t.Fatal("Quick Learner achievement no encontrado después de completar")
		}

		if quickAchievement.Progress != 5 {
			t.Errorf("Quick Learner debería tener progreso 5/5, tiene %d", quickAchievement.Progress)
		}

		if !quickAchievement.Completed {
			t.Error("Quick Learner debería estar completado")
		}

		if quickAchievement.UnlockedAt == nil {
			t.Error("Quick Learner debería tener UnlockedAt definido")
		}

		// Verificar XP total
		expectedBaseXP := 5 * 10 // 5 understand_insight * 10 XP cada uno
		if result.TotalXP < expectedBaseXP {
			t.Errorf("XP total debería ser al menos %d (base XP), obtenido %d",
				expectedBaseXP, result.TotalXP)
		}

		t.Logf("Achievement completado: Quick Learner | XP total: %d | Unlock time: %v",
			result.TotalXP, quickAchievement.UnlockedAt)
	})
}

// TestXPContributionFromDifferentSources valida XP de diferentes fuentes
func TestXPContributionFromDifferentSources(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()
	userID := "test_user_xp_sources"

	// Setup: usuario inicial
	mockRepo.SetupUser(userID, 0, 1)

	// Categorías de acciones con diferentes valores de XP
	actionCategories := []struct {
		name          string
		actions       []string
		expectedMinXP int
		description   string
	}{
		{
			name:          "basic_actions",
			actions:       []string{"view_dashboard", "view_expenses", "view_analytics"},
			expectedMinXP: 6, // 2 + 1 + 3 = 6
			description:   "Acciones básicas de navegación",
		},
		{
			name:          "transaction_actions",
			actions:       []string{"create_expense", "create_income", "update_expense"},
			expectedMinXP: 21, // 8 + 8 + 5 = 21
			description:   "Acciones de transacciones",
		},
		{
			name:          "organization_actions",
			actions:       []string{"create_category", "update_category", "assign_category"},
			expectedMinXP: 18, // 10 + 5 + 3 = 18
			description:   "Acciones de organización",
		},
		{
			name:          "engagement_actions",
			actions:       []string{"daily_login", "weekly_streak"},
			expectedMinXP: 30, // 5 + 25 = 30
			description:   "Acciones de engagement",
		},
		{
			name:          "advanced_features",
			actions:       []string{"create_budget", "create_savings_goal"},
			expectedMinXP: 35, // 20 + 15 = 35
			description:   "Features avanzadas",
		},
	}

	cumulativeXP := 0
	for categoryNum, category := range actionCategories {
		t.Run(category.name, func(t *testing.T) {
			categoryXP := 0

			for actionNum, actionType := range category.actions {
				entityType := getEntityTypeForAction(actionType)

				result, err := service.RecordUserAction(ctx, RecordActionParams{
					UserID:      userID,
					ActionType:  actionType,
					EntityType:  entityType,
					EntityID:    fmt.Sprintf("entity_c%d_a%d", categoryNum+1, actionNum+1),
					Description: fmt.Sprintf("%s - %s", category.description, actionType),
				})

				if err != nil {
					t.Fatalf("Error en %s - %s: %v", category.name, actionType, err)
				}

				categoryXP += result.XPEarned
			}

			cumulativeXP += categoryXP

			// Verificar XP mínimo de la categoría
			if categoryXP < category.expectedMinXP {
				t.Errorf("%s: XP esperado mínimo %d, obtenido %d",
					category.description, category.expectedMinXP, categoryXP)
			}

			// Verificar XP acumulativo
			profile, err := service.GetGamificationStats(ctx, userID)
			if err != nil {
				t.Fatalf("Error obteniendo stats: %v", err)
			}

			if profile.TotalXP < cumulativeXP {
				t.Errorf("XP acumulado esperado mínimo %d, obtenido %d",
					cumulativeXP, profile.TotalXP)
			}

			t.Logf("%s: +%d XP esta categoría = %d XP total (Nivel %d)",
				category.description, categoryXP, profile.TotalXP, profile.CurrentLevel)
		})
	}
}
