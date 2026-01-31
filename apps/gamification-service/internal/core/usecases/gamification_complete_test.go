package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/melegattip/financial-gamification-service/testutil"
)

// TestCompleteGamificationWorkflow tests the complete workflow with all real functions
func TestCompleteGamificationWorkflow(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()

	userID := "complete_workflow_user"

	// Test 1: Initialize user gamification (exercises InitializeUserGamification)
	t.Run("initialize_user_gamification", func(t *testing.T) {
		gamification, err := service.InitializeUserGamification(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to initialize user: %v", err)
		}

		if gamification.UserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, gamification.UserID)
		}

		if gamification.TotalXP != 0 {
			t.Errorf("Expected initial XP 0, got %d", gamification.TotalXP)
		}

		if gamification.CurrentLevel != 1 {
			t.Errorf("Expected initial level 1, got %d", gamification.CurrentLevel)
		}

		t.Logf("✅ User initialized: %s with %d XP, Level %d",
			gamification.UserID, gamification.TotalXP, gamification.CurrentLevel)
	})

	// Test 2: Get user gamification (exercises GetUserGamification)
	t.Run("get_user_gamification", func(t *testing.T) {
		gamification, err := service.GetUserGamification(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get user gamification: %v", err)
		}

		if gamification.UserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, gamification.UserID)
		}

		// Verify level calculation is working
		calculatedLevel := gamification.CalculateLevel()
		if gamification.CurrentLevel != calculatedLevel {
			t.Errorf("Level mismatch: stored %d vs calculated %d",
				gamification.CurrentLevel, calculatedLevel)
		}

		t.Logf("✅ User retrieved: %s with %d XP, Level %d",
			gamification.UserID, gamification.TotalXP, gamification.CurrentLevel)
	})

	// Test 3: Record multiple actions (exercises RecordUserAction extensively)
	actionsToTest := []struct {
		actionType  string
		entityType  string
		expectedXP  int
		description string
	}{
		{"view_dashboard", "dashboard", 2, "Dashboard view"},
		{"create_expense", "expense", 8, "First expense"},
		{"create_category", "category", 10, "First category"},
		{"view_analytics", "analytics", 3, "Analytics view"},
		{"daily_login", "user", 5, "Daily login"},
		{"create_income", "income", 8, "First income"},
		{"update_expense", "expense", 5, "Update expense"},
		{"assign_category", "category", 3, "Assign category"},
		{"create_savings_goal", "goal", 15, "Savings goal"},
		{"create_budget", "budget", 20, "Budget creation"},
	}

	totalXP := 0
	for i, action := range actionsToTest {
		t.Run(fmt.Sprintf("record_action_%d_%s", i+1, action.actionType), func(t *testing.T) {
			result, err := service.RecordUserAction(ctx, RecordActionParams{
				UserID:      userID,
				ActionType:  action.actionType,
				EntityType:  action.entityType,
				EntityID:    fmt.Sprintf("entity_%d", i+1),
				Description: action.description,
			})

			if err != nil {
				t.Fatalf("Failed to record action %s: %v", action.actionType, err)
			}

			// Verify XP earned is at least the base XP (can be higher due to achievement bonuses)
			if result.XPEarned < action.expectedXP {
				t.Errorf("Action %s: expected at least %d XP, got %d",
					action.actionType, action.expectedXP, result.XPEarned)
			}

			// Use authoritative total from result (it may include achievement bonus XP)
			if result.TotalXP < totalXP+action.expectedXP {
				t.Errorf("Action %s: expected total at least %d XP, got %d",
					action.actionType, totalXP+action.expectedXP, result.TotalXP)
			}

			totalXP = result.TotalXP

			// Verify level calculation using actual total XP
			expectedLevel := calculateLevelForXP(totalXP)
			if result.NewLevel != expectedLevel {
				t.Errorf("Action %s: expected level %d, got %d",
					action.actionType, expectedLevel, result.NewLevel)
			}

			// Check for level up
			if i > 0 {
				// Recompute previous level based on the total XP before this action
				previousLevel := calculateLevelForXP(totalXP - result.XPEarned)
				expectedLevelUp := result.NewLevel > previousLevel
				if result.LevelUp != expectedLevelUp {
					t.Errorf("Action %s: level up detection mismatch. Expected %v, got %v",
						action.actionType, expectedLevelUp, result.LevelUp)
				}
			}

			t.Logf("✅ Action %d (%s): +%d XP = %d total (Level %d, LevelUp: %v)",
				i+1, action.actionType, result.XPEarned, result.TotalXP, result.NewLevel, result.LevelUp)
		})
	}

	// Test 4: Get user achievements (exercises GetUserAchievements)
	t.Run("get_user_achievements", func(t *testing.T) {
		achievements, err := service.GetUserAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get achievements: %v", err)
		}

		if len(achievements) == 0 {
			t.Error("Expected some achievements to be initialized")
		}

		// Log achievement states
		for _, achievement := range achievements {
			t.Logf("Achievement: %s (%s) - Progress: %d/%d, Completed: %v",
				achievement.Name, achievement.Type, achievement.Progress, achievement.Target, achievement.Completed)
		}

		t.Logf("✅ Retrieved %d achievements", len(achievements))
	})

	// Test 5: Check and update achievements (exercises CheckAndUpdateAchievements)
	t.Run("check_and_update_achievements", func(t *testing.T) {
		newAchievements, updatedAchievements, err := service.CheckAndUpdateAchievements(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to check achievements: %v", err)
		}

		t.Logf("✅ Achievement check: %d new, %d updated",
			len(newAchievements), len(updatedAchievements))

		// Log any new achievements
		for _, achievement := range newAchievements {
			t.Logf("New achievement: %s (%s)", achievement.Name, achievement.Type)
		}

		// Log any updated achievements
		for _, achievement := range updatedAchievements {
			t.Logf("Updated achievement: %s - Progress: %d/%d",
				achievement.Name, achievement.Progress, achievement.Target)
		}
	})

	// Test 6: Get gamification stats (exercises GetGamificationStats)
	t.Run("get_gamification_stats", func(t *testing.T) {
		stats, err := service.GetGamificationStats(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}

		if stats.UserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, stats.UserID)
		}

		if stats.TotalXP != totalXP {
			t.Errorf("Expected total XP %d, got %d", totalXP, stats.TotalXP)
		}

		// Verify calculated fields
		expectedLevel := calculateLevelForXP(totalXP)
		if stats.CurrentLevel != expectedLevel {
			t.Errorf("Expected level %d, got %d", expectedLevel, stats.CurrentLevel)
		}

		t.Logf("✅ Stats retrieved: %d XP, Level %d, %d%% progress, %d/%d achievements",
			stats.TotalXP, stats.CurrentLevel, stats.ProgressPercent,
			stats.CompletedAchievements, stats.TotalAchievements)
	})

	// Test 7: Get action types (exercises GetActionTypes)
	t.Run("get_action_types", func(t *testing.T) {
		actionTypes, err := service.GetActionTypes(ctx)
		if err != nil {
			t.Fatalf("Failed to get action types: %v", err)
		}

		if len(actionTypes) == 0 {
			t.Error("Expected some action types")
		}

		// Verify some expected action types
		expectedActions := map[string]int{
			"view_dashboard":      2,
			"create_expense":      8,
			"create_category":     10,
			"daily_login":         5,
			"create_budget":       20,
			"create_savings_goal": 15,
		}

		foundActions := make(map[string]int)
		for _, actionType := range actionTypes {
			foundActions[actionType.Type] = actionType.BaseXP
		}

		for expectedAction, expectedXP := range expectedActions {
			if foundXP, exists := foundActions[expectedAction]; !exists {
				t.Errorf("Expected action type %s not found", expectedAction)
			} else if foundXP != expectedXP {
				t.Errorf("Action %s: expected %d XP, got %d", expectedAction, expectedXP, foundXP)
			}
		}

		t.Logf("✅ Retrieved %d action types", len(actionTypes))
	})

	// Test 8: Get levels (exercises GetLevels)
	t.Run("get_levels", func(t *testing.T) {
		levels, err := service.GetLevels(ctx)
		if err != nil {
			t.Fatalf("Failed to get levels: %v", err)
		}

		if len(levels) == 0 {
			t.Error("Expected some level information")
		}

		// Verify level progression makes sense
		for i, level := range levels {
			if level.Level != i+1 {
				t.Errorf("Level %d: expected level number %d, got %d", i, i+1, level.Level)
			}

			if level.XPRequired < 0 {
				t.Errorf("Level %d: negative required XP %d", level.Level, level.XPRequired)
			}

			if i > 0 && level.XPRequired <= levels[i-1].XPRequired {
				t.Errorf("Level %d: XP requirement %d should be higher than previous level %d",
					level.Level, level.XPRequired, levels[i-1].XPRequired)
			}

			t.Logf("Level %d: %s - Requires %d XP",
				level.Level, level.Name, level.XPRequired)
		}

		t.Logf("✅ Retrieved %d levels", len(levels))
	})
}

// TestFeatureAccessAndGates exercises feature gate functionality
func TestFeatureAccessAndGates(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()

	userID := "feature_test_user"

	// Setup user with specific level
	mockRepo.SetupUser(userID, 50, 1) // Level 1 user

	// Test feature access for different levels (usando features reales)
	featureTests := []struct {
		featureKey     string
		userLevel      int
		userXP         int
		expectedAccess bool
		description    string
	}{
		{"SAVINGS_GOALS", 1, 50, false, "Savings goals require level 3+ (user is level 1)"},
		{"BUDGETS", 1, 50, false, "Budgets require level 5+ (user is level 1)"},
		{"AI_INSIGHTS", 1, 50, false, "AI insights require level 7+ (user is level 1)"},
	}

	for _, test := range featureTests {
		t.Run(fmt.Sprintf("feature_%s_level_%d", test.featureKey, test.userLevel), func(t *testing.T) {
			result, err := service.CheckFeatureAccess(ctx, userID, test.featureKey)
			if err != nil {
				t.Fatalf("Failed to check feature access for %s: %v", test.featureKey, err)
			}

			if result.HasAccess != test.expectedAccess {
				t.Errorf("%s: expected access %v, got %v",
					test.description, test.expectedAccess, result.HasAccess)
			}

			t.Logf("✅ Feature %s: Access=%v, Description=%s",
				test.featureKey, result.HasAccess, result.Description)
		})
	}

	// Test getting all user features
	t.Run("get_all_user_features", func(t *testing.T) {
		features, err := service.GetUserFeatures(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get user features: %v", err)
		}

		if len(features.UnlockedFeatures) == 0 && len(features.LockedFeatures) == 0 {
			t.Error("Expected some features to be returned")
		}

		t.Logf("✅ User features: %d unlocked, %d locked",
			len(features.UnlockedFeatures), len(features.LockedFeatures))

		for _, featureKey := range features.UnlockedFeatures {
			t.Logf("Unlocked: %s", featureKey)
		}

		for _, feature := range features.LockedFeatures {
			t.Logf("Locked: %s - %s (Requires Level %d)",
				feature.FeatureKey, feature.FeatureName, feature.RequiredLevel)
		}
	})
}

// TestChallengeSystem exercises challenge-related functionality
func TestChallengeSystem(t *testing.T) {
	mockRepo := testutil.NewMockGamificationRepository()
	service := NewGamificationUseCase(mockRepo)
	ctx := context.Background()

	userID := "challenge_test_user"
	mockRepo.SetupUser(userID, 100, 2)

	// Test getting daily challenges
	t.Run("get_daily_challenges", func(t *testing.T) {
		challenges, err := service.GetDailyChallenges(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get daily challenges: %v", err)
		}

		t.Logf("✅ Retrieved %d daily challenges", len(challenges))

		for _, challenge := range challenges {
			t.Logf("Daily Challenge: %s - Progress %d/%d (%d%% complete)",
				challenge.Name, challenge.Progress, challenge.Target, challenge.ProgressPercent)
		}
	})

	// Test getting weekly challenges
	t.Run("get_weekly_challenges", func(t *testing.T) {
		challenges, err := service.GetWeeklyChallenges(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get weekly challenges: %v", err)
		}

		t.Logf("✅ Retrieved %d weekly challenges", len(challenges))

		for _, challenge := range challenges {
			t.Logf("Weekly Challenge: %s - Progress %d/%d",
				challenge.Name, challenge.Progress, challenge.Target)
		}
	})

	// Test processing challenge progress
	t.Run("process_challenge_progress", func(t *testing.T) {
		result, err := service.ProcessChallengeProgress(ctx, userID, "create_expense", "expense", "test_expense")
		if err != nil {
			t.Fatalf("Failed to process challenge progress: %v", err)
		}

		t.Logf("✅ Challenge progress: %d updated, %d completed, %d XP earned",
			len(result.UpdatedChallenges), len(result.CompletedChallenges), result.TotalXPEarned)

		for _, challenge := range result.UpdatedChallenges {
			t.Logf("Updated: %s - Progress %d/%d",
				challenge.Name, challenge.Progress, challenge.Target)
		}

		for _, challenge := range result.CompletedChallenges {
			t.Logf("Completed: %s - Earned %d XP", challenge.Name, challenge.XPReward)
		}
	})
}

// Helper function to calculate level for given XP (matches domain logic)
func calculateLevelForXP(xp int) int {
	levels := []int{0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500}

	for i := len(levels) - 1; i >= 0; i-- {
		if xp >= levels[i] {
			return i + 1 // Return level 1-10 instead of 0-9
		}
	}
	return 1 // Minimum level is 1
}
