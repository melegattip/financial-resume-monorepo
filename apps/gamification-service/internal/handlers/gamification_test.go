package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/melegattip/financial-gamification-service/internal/core/usecases"
	"github.com/melegattip/financial-gamification-service/internal/infrastructure/http/middleware"
	"github.com/melegattip/financial-gamification-service/testutil"
)

// TestBasicHandlerSetup valida la configuración básica de handlers
func TestBasicHandlerSetup(t *testing.T) {
	// Setup mock repository
	repository := testutil.NewMockGamificationRepository()

	if repository == nil {
		t.Fatal("Mock repository no debería ser nil")
	}

	// Verificar que se puede crear un usuario
	repository.SetupUser("test_user", 100, 2)

	// Test básico para verificar que el setup funciona
	userID := "test_user"
	ctx := context.Background()

	user, err := repository.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Error al obtener usuario: %v", err)
	}

	if user.TotalXP != 100 {
		t.Errorf("XP total debería ser 100 pero es %d", user.TotalXP)
	}

	if user.CurrentLevel != 2 {
		t.Errorf("Nivel debería ser 2 pero es %d", user.CurrentLevel)
	}
}

// TestHandlerWithMultipleUsers valida el manejo de múltiples usuarios
func TestHandlerWithMultipleUsers(t *testing.T) {
	repository := testutil.NewMockGamificationRepository()
	ctx := context.Background()

	// Setup múltiples usuarios
	users := []struct {
		id    string
		xp    int
		level int
	}{
		{"user_1", 50, 1},
		{"user_2", 150, 2},
		{"user_3", 300, 3},
	}

	for _, u := range users {
		repository.SetupUser(u.id, u.xp, u.level)
	}

	// Verificar que cada usuario mantiene sus datos
	for _, u := range users {
		user, err := repository.GetByUserID(ctx, u.id)
		if err != nil {
			t.Fatalf("Error al obtener usuario %s: %v", u.id, err)
		}

		if user.TotalXP != u.xp {
			t.Errorf("Usuario %s debería tener %d XP pero tiene %d", u.id, u.xp, user.TotalXP)
		}

		if user.CurrentLevel != u.level {
			t.Errorf("Usuario %s debería estar en nivel %d pero está en %d", u.id, u.level, user.CurrentLevel)
		}
	}
}

// TestAchievementHandling valida el manejo básico de achievements
func TestAchievementHandling(t *testing.T) {
	repository := testutil.NewMockGamificationRepository()
	ctx := context.Background()
	userID := "test_achievements_user"

	// Setup usuario
	repository.SetupUser(userID, 0, 1)

	// Obtener achievements iniciales
	achievements, err := repository.GetAchievementsByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Error al obtener achievements: %v", err)
	}

	if len(achievements) == 0 {
		t.Error("Debería haber achievements inicializados")
	}

	// Verificar estructura básica de achievement
	for _, ach := range achievements {
		if ach.UserID != userID {
			t.Errorf("Achievement debería pertenecer al usuario %s pero pertenece a %s", userID, ach.UserID)
		}

		if ach.Target <= 0 {
			t.Errorf("Achievement debería tener un target mayor a 0 pero tiene %d", ach.Target)
		}

		if ach.Progress < 0 {
			t.Errorf("Achievement no debería tener progreso negativo pero tiene %d", ach.Progress)
		}

		if ach.Name == "" {
			t.Error("Achievement debería tener un nombre")
		}

		if ach.Type == "" {
			t.Error("Achievement debería tener un tipo")
		}
	}
}

// TestRealHTTPHandlers tests real HTTP handlers with actual HTTP requests
func TestRealHTTPHandlers(t *testing.T) {
	// Setup real dependencies
	mockRepo := testutil.NewMockGamificationRepository()
	service := usecases.NewGamificationUseCase(mockRepo)
	handlers := NewGamificationHandlers(service)

	// Setup router
	router := mux.NewRouter()

	// Register actual routes (using simplified paths for testing)
	router.HandleFunc("/api/v1/gamification/profile", handlers.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/v1/gamification/actions", handlers.RecordUserAction).Methods("POST")
	router.HandleFunc("/api/v1/gamification/achievements", handlers.GetUserAchievements).Methods("GET")
	router.HandleFunc("/api/v1/gamification/stats", handlers.GetUserStats).Methods("GET")

	// Test data setup
	userID := "test_user_real"
	mockRepo.SetupUser(userID, 100, 2)

	t.Run("GET_user_profile_real_http", func(t *testing.T) {
		// Create real HTTP request
		req, err := http.NewRequest("GET", "/api/v1/gamification/profile", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		// Add user ID to context (simulating JWT middleware)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		// Create real HTTP response recorder
		rr := httptest.NewRecorder()

		// Execute real HTTP handler
		router.ServeHTTP(rr, req)

		// Verify HTTP response
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
		}

		// Verify response body
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Verify actual data
		if response["user_id"] != userID {
			t.Errorf("Expected user_id %s, got %v", userID, response["user_id"])
		}

		if int(response["total_xp"].(float64)) != 100 {
			t.Errorf("Expected total_xp 100, got %v", response["total_xp"])
		}

		if int(response["current_level"].(float64)) != 2 {
			t.Errorf("Expected current_level 2, got %v", response["current_level"])
		}

		t.Logf("✅ Real HTTP GET worked: Status %d, User %s, XP %v, Level %v",
			rr.Code, response["user_id"], response["total_xp"], response["current_level"])
	})

	t.Run("POST_record_action_real_http", func(t *testing.T) {
		// Create real action request
		actionRequest := map[string]interface{}{
			"action_type": "create_expense",
			"entity_type": "expense",
			"entity_id":   "expense_123",
			"description": "Test expense creation",
		}

		requestBody, _ := json.Marshal(actionRequest)

		// Create real HTTP POST request
		req, err := http.NewRequest("POST", "/api/v1/gamification/actions", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Add user ID to context (simulating JWT middleware)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		// Create real HTTP response recorder
		rr := httptest.NewRecorder()

		// Execute real HTTP handler
		router.ServeHTTP(rr, req)

		// Verify HTTP response
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		// Verify response body
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Verify XP was earned (create_expense should give 8 XP)
		if int(response["xp_earned"].(float64)) != 8 {
			t.Errorf("Expected xp_earned 8, got %v", response["xp_earned"])
		}

		// Verify total XP updated (100 + 8 = 108)
		if int(response["total_xp"].(float64)) != 108 {
			t.Errorf("Expected total_xp 108, got %v", response["total_xp"])
		}

		t.Logf("✅ Real HTTP POST worked: Status %d, XP Earned %v, Total XP %v",
			rr.Code, response["xp_earned"], response["total_xp"])
	})

	t.Run("GET_achievements_real_http", func(t *testing.T) {
		// Create real HTTP request
		req, err := http.NewRequest("GET", "/api/v1/gamification/achievements", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		// Add user ID to context (simulating JWT middleware)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		// Create real HTTP response recorder
		rr := httptest.NewRecorder()

		// Execute real HTTP handler
		router.ServeHTTP(rr, req)

		// Verify HTTP response
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
		}

		// Verify response is valid JSON array
		var achievements []map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &achievements)
		if err != nil {
			t.Fatalf("Error parsing achievements response: %v", err)
		}

		// Verify we have achievements
		if len(achievements) == 0 {
			t.Error("Expected at least some achievements")
		}

		t.Logf("✅ Real HTTP GET achievements worked: Status %d, Count %d", rr.Code, len(achievements))
	})

	t.Run("GET_stats_real_http", func(t *testing.T) {
		// Create real HTTP request
		req, err := http.NewRequest("GET", "/api/v1/gamification/stats", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		// Add user ID to context (simulating JWT middleware)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
		req = req.WithContext(ctx)

		// Create real HTTP response recorder
		rr := httptest.NewRecorder()

		// Execute real HTTP handler
		router.ServeHTTP(rr, req)

		// Verify HTTP response
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
		}

		// Verify response body
		var stats map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &stats)
		if err != nil {
			t.Fatalf("Error parsing stats response: %v", err)
		}

		// Verify stats contain expected fields
		expectedFields := []string{"user_id", "total_xp", "current_level", "xp_to_next_level", "progress_percent"}
		for _, field := range expectedFields {
			if _, exists := stats[field]; !exists {
				t.Errorf("Expected field %s in stats response", field)
			}
		}

		t.Logf("✅ Real HTTP GET stats worked: Status %d, Level %v, XP %v",
			rr.Code, stats["current_level"], stats["total_xp"])
	})
}

// TestCompleteUserJourney tests end-to-end user flow through real HTTP handlers
func TestCompleteUserJourney(t *testing.T) {
	// Setup real dependencies
	mockRepo := testutil.NewMockGamificationRepository()
	service := usecases.NewGamificationUseCase(mockRepo)
	handlers := NewGamificationHandlers(service)

	// Setup router with real routes
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/gamification/profile", handlers.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/v1/gamification/actions", handlers.RecordUserAction).Methods("POST")
	router.HandleFunc("/api/v1/gamification/achievements", handlers.GetUserAchievements).Methods("GET")
	router.HandleFunc("/api/v1/gamification/stats", handlers.GetUserStats).Methods("GET")

	userID := "test_user_journey"

	// Step 1: Initialize new user (should create automatically)
	t.Run("step_1_initialize_user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/gamification/profile", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to initialize user: %d", rr.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		// Verify initial state
		if int(response["total_xp"].(float64)) != 0 {
			t.Errorf("Expected initial XP 0, got %v", response["total_xp"])
		}

		if int(response["current_level"].(float64)) != 1 {
			t.Errorf("Expected initial level 1, got %v", response["current_level"])
		}

		t.Logf("✅ User initialized: XP %v, Level %v", response["total_xp"], response["current_level"])
	})

	// Step 2: Perform actions to gain XP and test level progression
	actions := []struct {
		actionType  string
		entityType  string
		expectedXP  int
		description string
	}{
		{"view_dashboard", "dashboard", 2, "First dashboard view"},
		{"create_expense", "expense", 8, "First expense"},
		{"create_category", "category", 10, "First category"},
		{"view_analytics", "analytics", 3, "Analytics view"},
		{"create_expense", "expense", 8, "Second expense"},
		{"daily_login", "user", 5, "Daily login"},
		{"create_budget", "budget", 20, "First budget"},  // Total: 56 XP (Level 1)
		{"create_budget", "budget", 20, "Second budget"}, // Total: 76 XP (Level 2!)
	}

	totalExpectedXP := 0
	for i, action := range actions {
		t.Run(fmt.Sprintf("step_2_action_%d_%s", i+1, action.actionType), func(t *testing.T) {
			actionRequest := map[string]interface{}{
				"action_type": action.actionType,
				"entity_type": action.entityType,
				"entity_id":   fmt.Sprintf("entity_%d", i+1),
				"description": action.description,
			}

			requestBody, _ := json.Marshal(actionRequest)
			req, _ := http.NewRequest("POST", "/api/v1/gamification/actions", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Add user ID to context (simulating JWT middleware)
			ctx := req.Context()
			ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("Action %d failed: %d - %s", i+1, rr.Code, rr.Body.String())
			}

			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)

			// Verify XP earned
			if int(response["xp_earned"].(float64)) != action.expectedXP {
				t.Errorf("Action %d: Expected XP %d, got %v", i+1, action.expectedXP, response["xp_earned"])
			}

			totalExpectedXP += action.expectedXP

			// Verify total XP
			if int(response["total_xp"].(float64)) != totalExpectedXP {
				t.Errorf("Action %d: Expected total XP %d, got %v", i+1, totalExpectedXP, response["total_xp"])
			}

			// Check for level up at action 8 (76 XP reaches level 2)
			if i == 7 { // 8th action (0-indexed)
				if !response["level_up"].(bool) {
					t.Error("Expected level_up to be true at 76 XP")
				}
				if int(response["new_level"].(float64)) != 2 {
					t.Errorf("Expected new level 2, got %v", response["new_level"])
				}
			}

			t.Logf("✅ Action %d (%s): +%d XP = %d total XP (Level up: %v)",
				i+1, action.actionType, action.expectedXP, totalExpectedXP,
				response["level_up"])
		})
	}

	// Step 3: Verify final state
	t.Run("step_3_verify_final_state", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/gamification/stats", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to get final stats: %d", rr.Code)
		}

		var stats map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &stats)

		finalXP := int(stats["total_xp"].(float64))
		finalLevel := int(stats["current_level"].(float64))

		t.Logf("✅ Final state: %d XP, Level %d", finalXP, finalLevel)

		// With 76 XP, user should be at level 2 (threshold is 75)
		if finalXP != 76 {
			t.Errorf("Expected final XP 76, got %d", finalXP)
		}

		if finalLevel != 2 {
			t.Errorf("Expected final level 2, got %d", finalLevel)
		}

		// Verify XP to next level (Level 3 requires 200 XP, so 200-76=124)
		xpToNext := int(stats["xp_to_next_level"].(float64))
		expectedXPToNext := 200 - 76 // 124 XP to reach level 3
		if xpToNext != expectedXPToNext {
			t.Errorf("Expected XP to next level %d, got %d", expectedXPToNext, xpToNext)
		}
	})
}

// TestHandlerErrorConditions tests error handling in HTTP handlers
func TestHandlerErrorConditions(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := usecases.NewGamificationUseCase(mockRepo)
	handlers := NewGamificationHandlers(service)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/gamification/profile", handlers.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/v1/gamification/actions", handlers.RecordUserAction).Methods("POST")

	t.Run("missing_user_context", func(t *testing.T) {
		// Request without user context should return error
		req, _ := http.NewRequest("GET", "/api/v1/gamification/profile", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK {
			t.Error("Expected error for missing user context, got success")
		}

		t.Logf("✅ Error handling verified: Status %d for missing user context", rr.Code)
	})

	t.Run("invalid_json_payload", func(t *testing.T) {
		// POST with invalid JSON should return error
		invalidJSON := []byte("{invalid json}")
		req, _ := http.NewRequest("POST", "/api/v1/gamification/actions", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, "test_user")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK {
			t.Error("Expected error for invalid JSON, got success")
		}

		t.Logf("✅ Error handling verified: Status %d for invalid JSON", rr.Code)
	})
}
