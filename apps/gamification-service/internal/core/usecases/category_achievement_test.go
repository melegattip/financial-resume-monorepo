package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/melegattip/financial-gamification-service/internal/core/domain"
	"github.com/melegattip/financial-gamification-service/testutil"
)

// TestCategoryAchievementBugFix verifica que el bug del achievement de categorías esté corregido
func TestCategoryAchievementBugFix(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()

	userID := "category_achievement_test_user"
	mockRepo.SetupUser(userID, 50, 2)

	t.Run("CategoryCreatorAchievement_Should_Progress_With_Real_Actions", func(t *testing.T) {
		// ✅ ANTES DEL FIX: El achievement usaba ActionsCompleted / 20 (estimación errónea)
		// ✅ DESPUÉS DEL FIX: El achievement cuenta las acciones create_category reales

		// Inicializar el sistema de gamificación del usuario (crea achievements básicos)
		_, err := service.InitializeUserGamification(ctx, userID)
		if err != nil {
			t.Fatalf("InitializeUserGamification failed: %v", err)
		}

		// Obtener achievements después de la inicialización
		achievements, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("GetUserAchievements failed: %v", err)
		}

		// Buscar el achievement category_creator
		var categoryCreatorAchievement *domain.Achievement
		for i := range achievements {
			if achievements[i].Type == "category_creator" {
				categoryCreatorAchievement = &achievements[i]
				break
			}
		}

		if categoryCreatorAchievement == nil {
			t.Fatal("Achievement category_creator not found")
		}

		// Verificar progreso inicial (debería ser 0 porque no ha creado categorías)
		initialProgress := categoryCreatorAchievement.Progress
		t.Logf("🔍 Progreso inicial del achievement category_creator: %d", initialProgress)

		// 🎯 SIMULAR CREACIÓN DE PRIMERA CATEGORÍA
		result1, err := service.RecordUserAction(ctx, RecordActionParams{
			UserID:      userID,
			ActionType:  "create_category",
			EntityType:  "category",
			EntityID:    "cat_food",
			Description: "Created category: Food",
		})
		if err != nil {
			t.Fatalf("RecordUserAction 1 failed: %v", err)
		}

		// Verificar que se otorgaron 10 XP por crear categoría
		if result1.XPEarned != 10 {
			t.Errorf("Expected 10 XP for create_category, got %d", result1.XPEarned)
		}
		t.Logf("✅ Primera categoría creada - XP ganado: %d", result1.XPEarned)

		// 🎯 SIMULAR CREACIÓN DE SEGUNDA CATEGORÍA
		result2, err := service.RecordUserAction(ctx, RecordActionParams{
			UserID:      userID,
			ActionType:  "create_category",
			EntityType:  "category",
			EntityID:    "cat_transport",
			Description: "Created category: Transport",
		})
		if err != nil {
			t.Fatalf("RecordUserAction 2 failed: %v", err)
		}

		t.Logf("✅ Segunda categoría creada - XP ganado: %d", result2.XPEarned)

		// 🎯 VERIFICAR PROGRESO DEL ACHIEVEMENT DESPUÉS DE 2 CATEGORÍAS
		achievementsAfter, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("GetUserAchievements after actions failed: %v", err)
		}

		var updatedAchievement *domain.Achievement
		for i := range achievementsAfter {
			if achievementsAfter[i].Type == "category_creator" {
				updatedAchievement = &achievementsAfter[i]
				break
			}
		}

		if updatedAchievement == nil {
			t.Fatal("Updated achievement category_creator not found")
		}

		// ✅ VERIFICACIÓN PRINCIPAL: El progreso debe ser 2 (2 categorías creadas)
		expectedProgress := 2
		if updatedAchievement.Progress != expectedProgress {
			t.Errorf("❌ Achievement progress INCORRECT: expected %d categories, got %d",
				expectedProgress, updatedAchievement.Progress)
			t.Errorf("   This indicates the bug is NOT fixed!")
		} else {
			t.Logf("✅ Achievement progress CORRECT: %d/%d categories (Bug FIXED!)",
				updatedAchievement.Progress, updatedAchievement.Target)
		}

		// Verificar que el achievement no está completado todavía (requiere 5)
		if updatedAchievement.Completed {
			t.Error("Achievement should not be completed yet (needs 5 categories)")
		}

		// 🎯 CREAR 3 CATEGORÍAS MÁS PARA COMPLETAR EL ACHIEVEMENT
		categoryNames := []string{"cat_entertainment", "cat_health", "cat_shopping"}
		for i, categoryName := range categoryNames {
			_, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID,
				ActionType:  "create_category",
				EntityType:  "category",
				EntityID:    categoryName,
				Description: fmt.Sprintf("Created category: %s", categoryName),
			})
			if err != nil {
				t.Fatalf("RecordUserAction %d failed: %v", i+3, err)
			}
		}

		// 🎯 VERIFICAR QUE EL ACHIEVEMENT SE COMPLETÓ
		finalAchievements, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("GetUserAchievements final failed: %v", err)
		}

		var finalAchievement *domain.Achievement
		for i := range finalAchievements {
			if finalAchievements[i].Type == "category_creator" {
				finalAchievement = &finalAchievements[i]
				break
			}
		}

		if finalAchievement == nil {
			t.Fatal("Final achievement category_creator not found")
		}

		// Verificar progreso final (5/5 categorías)
		if finalAchievement.Progress != 5 {
			t.Errorf("Final progress incorrect: expected 5, got %d", finalAchievement.Progress)
		}

		// Verificar que el achievement está completado
		if !finalAchievement.Completed {
			t.Error("Achievement should be completed after creating 5 categories")
		}

		t.Logf("🎉 ¡ACHIEVEMENT COMPLETADO! 5/5 categorías creadas")
		t.Logf("📊 Total XP del usuario: %d", result1.TotalXP+result2.XPEarned+(3*10))
	})
}

// TestCategoryAchievementCountsRealActions verifica que solo las acciones create_category se cuenten
func TestCategoryAchievementCountsRealActions(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()

	userID := "real_actions_test_user"
	mockRepo.SetupUser(userID, 100, 3)

	t.Run("Only_Create_Category_Actions_Should_Count", func(t *testing.T) {
		// Inicializar el sistema de gamificación del usuario
		_, err := service.InitializeUserGamification(ctx, userID)
		if err != nil {
			t.Fatalf("InitializeUserGamification failed: %v", err)
		}

		// Realizar muchas acciones que NO son create_category
		nonCategoryActions := []RecordActionParams{
			{UserID: userID, ActionType: "create_expense", EntityType: "expense", EntityID: "exp1", Description: "Expense 1"},
			{UserID: userID, ActionType: "create_income", EntityType: "income", EntityID: "inc1", Description: "Income 1"},
			{UserID: userID, ActionType: "view_dashboard", EntityType: "dashboard", EntityID: "main", Description: "View dashboard"},
			{UserID: userID, ActionType: "view_analytics", EntityType: "analytics", EntityID: "monthly", Description: "View analytics"},
			{UserID: userID, ActionType: "update_category", EntityType: "category", EntityID: "cat1", Description: "Update category"},
			{UserID: userID, ActionType: "assign_category", EntityType: "expense", EntityID: "exp2", Description: "Assign category"},
		}

		// Ejecutar todas las acciones no-categoría
		for _, action := range nonCategoryActions {
			_, err := service.RecordUserAction(ctx, action)
			if err != nil {
				t.Fatalf("Non-category action failed: %v", err)
			}
		}

		t.Logf("✅ Ejecutadas %d acciones que NO son create_category", len(nonCategoryActions))

		// Verificar que el achievement category_creator sigue en 0
		achievements, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("GetUserAchievements failed: %v", err)
		}

		var categoryAchievement *domain.Achievement
		for i := range achievements {
			if achievements[i].Type == "category_creator" {
				categoryAchievement = &achievements[i]
				break
			}
		}

		if categoryAchievement == nil {
			t.Fatal("Achievement category_creator not found")
		}

		// ✅ VERIFICACIÓN: Debe ser 0 porque NO se crearon categorías
		if categoryAchievement.Progress != 0 {
			t.Errorf("❌ Progress should be 0 after non-category actions, got %d", categoryAchievement.Progress)
			t.Error("   This indicates the achievement is counting wrong actions!")
		} else {
			t.Logf("✅ Correct: Progress is 0 after %d non-category actions", len(nonCategoryActions))
		}

		// Ahora crear UNA categoría real
		_, err = service.RecordUserAction(ctx, RecordActionParams{
			UserID:      userID,
			ActionType:  "create_category",
			EntityType:  "category",
			EntityID:    "cat_real",
			Description: "Real category creation",
		})
		if err != nil {
			t.Fatalf("Real category creation failed: %v", err)
		}

		// Verificar que ahora el progreso es exactamente 1
		updatedAchievements, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("GetUserAchievements after real category failed: %v", err)
		}

		var updatedCategoryAchievement *domain.Achievement
		for i := range updatedAchievements {
			if updatedAchievements[i].Type == "category_creator" {
				updatedCategoryAchievement = &updatedAchievements[i]
				break
			}
		}

		if updatedCategoryAchievement.Progress != 1 {
			t.Errorf("❌ Progress should be 1 after creating 1 category, got %d", updatedCategoryAchievement.Progress)
		} else {
			t.Logf("✅ Perfect: Progress is 1 after creating exactly 1 category")
		}
	})
}
