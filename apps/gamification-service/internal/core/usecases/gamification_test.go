package usecases

import (
	"context"
	"testing"

	"github.com/melegattip/financial-gamification-service/internal/core/domain"
	"github.com/melegattip/financial-gamification-service/testutil"
)

// TestBasicUserSetup valida la configuración básica de usuarios para testing
func TestBasicUserSetup(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	ctx := context.Background()
	userID := "test_user_basic"

	// Setup: crear usuario inicial
	mockRepo.SetupUser(userID, 50, 1)

	// Act: obtener usuario
	user, err := mockRepo.GetByUserID(ctx, userID)

	// Assert: verificar configuración
	if err != nil {
		t.Fatalf("No debería haber error al obtener usuario: %v", err)
	}

	if user == nil {
		t.Fatal("El usuario no debería ser nil")
	}

	if user.TotalXP != 50 {
		t.Errorf("XP total debería ser 50 pero es %d", user.TotalXP)
	}

	if user.CurrentLevel != 1 {
		t.Errorf("Nivel actual debería ser 1 pero es %d", user.CurrentLevel)
	}
}

// TestAchievementCreation valida la creación de achievements
func TestAchievementCreation(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	ctx := context.Background()
	userID := "test_user_achievements"

	// Setup: crear usuario
	mockRepo.SetupUser(userID, 0, 1)

	// Verificar que se inicializaron achievements básicos
	achievements, err := mockRepo.GetAchievementsByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Error al obtener achievements: %v", err)
	}

	if len(achievements) == 0 {
		t.Error("Debería haber achievements inicializados")
	}

	// Verificar que hay al menos un achievement de tipo ai_partner
	found := false
	for _, ach := range achievements {
		if ach.Type == "ai_partner" {
			found = true
			if ach.UserID != userID {
				t.Errorf("Achievement debería pertenecer al usuario %s", userID)
			}
			if ach.Progress != 0 {
				t.Errorf("Achievement nuevo debería tener progreso 0, pero tiene %d", ach.Progress)
			}
			break
		}
	}

	if !found {
		t.Error("Debería haber un achievement de tipo ai_partner")
	}
}

// TestMultipleUsersIsolation valida que los datos de usuarios estén aislados
func TestMultipleUsersIsolation(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	ctx := context.Background()

	user1ID := "user_1"
	user2ID := "user_2"

	// Setup: crear dos usuarios con diferentes configuraciones
	mockRepo.SetupUser(user1ID, 100, 2)
	mockRepo.SetupUser(user2ID, 200, 3)

	// Verificar aislamiento
	user1, err := mockRepo.GetByUserID(ctx, user1ID)
	if err != nil {
		t.Fatalf("Error al obtener user1: %v", err)
	}

	user2, err := mockRepo.GetByUserID(ctx, user2ID)
	if err != nil {
		t.Fatalf("Error al obtener user2: %v", err)
	}

	if user1.TotalXP != 100 {
		t.Errorf("User1 debería tener 100 XP, pero tiene %d", user1.TotalXP)
	}

	if user2.TotalXP != 200 {
		t.Errorf("User2 debería tener 200 XP, pero tiene %d", user2.TotalXP)
	}

	if user1.CurrentLevel != 2 {
		t.Errorf("User1 debería estar en nivel 2, pero está en %d", user1.CurrentLevel)
	}

	if user2.CurrentLevel != 3 {
		t.Errorf("User2 debería estar en nivel 3, pero está en %d", user2.CurrentLevel)
	}

	// Verificar que los achievements son independientes
	achievements1, err := mockRepo.GetAchievementsByUserID(ctx, user1ID)
	if err != nil {
		t.Fatalf("Error al obtener achievements de user1: %v", err)
	}

	achievements2, err := mockRepo.GetAchievementsByUserID(ctx, user2ID)
	if err != nil {
		t.Fatalf("Error al obtener achievements de user2: %v", err)
	}

	// Cada usuario debería tener sus propios achievements
	if len(achievements1) == 0 {
		t.Error("User1 debería tener achievements")
	}

	if len(achievements2) == 0 {
		t.Error("User2 debería tener achievements")
	}

	// Verificar que los achievements pertenecen al usuario correcto
	for _, ach := range achievements1 {
		if ach.UserID != user1ID {
			t.Errorf("Achievement de user1 tiene UserID incorrecto: %s", ach.UserID)
		}
	}

	for _, ach := range achievements2 {
		if ach.UserID != user2ID {
			t.Errorf("Achievement de user2 tiene UserID incorrecto: %s", ach.UserID)
		}
	}
}

// TestActionLogging valida el registro de acciones
func TestActionLogging(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	ctx := context.Background()
	userID := "test_user_actions"

	// Setup: crear usuario
	mockRepo.SetupUser(userID, 0, 1)

	// Crear una acción manualmente (usamos el tipo del dominio)
	action := &domain.UserAction{
		ID:          "action_1",
		UserID:      userID,
		ActionType:  "view_insight",
		EntityType:  "insight",
		EntityID:    "insight_123",
		XPEarned:    5,
		Description: "Test action",
	}

	// Simular el registro de una acción
	err := mockRepo.CreateAction(ctx, action)
	if err != nil {
		t.Fatalf("Error al crear acción: %v", err)
	}

	// Verificar que la acción se registró
	actions, err := mockRepo.GetActionsByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Error al obtener acciones: %v", err)
	}

	if len(actions) != 1 {
		t.Errorf("Debería haber 1 acción registrada, pero hay %d", len(actions))
	}

	if len(actions) > 0 {
		if actions[0].ActionType != "view_insight" {
			t.Errorf("Tipo de acción debería ser 'view_insight', pero es '%s'", actions[0].ActionType)
		}
		if actions[0].UserID != userID {
			t.Errorf("UserID de la acción debería ser '%s', pero es '%s'", userID, actions[0].UserID)
		}
	}
}
