package domain

import (
	"fmt"
	"testing"
)

// TestAchievementProgressCalculation valida el cálculo de progreso de achievements
func TestAchievementProgressCalculation(t *testing.T) {
	tests := []struct {
		name             string
		progress         int
		target           int
		expectedComplete bool
	}{
		{"no_progress", 0, 10, false},
		{"partial_progress", 3, 10, false},
		{"halfway", 5, 10, false},
		{"almost_complete", 9, 10, false},
		{"just_completed", 10, 10, true},
		{"exceeded", 12, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achievement := &Achievement{
				Progress: tt.progress,
				Target:   tt.target,
			}

			// Test completion status
			actualComplete := achievement.IsCompleted()
			if actualComplete != tt.expectedComplete {
				t.Errorf("Achievement con progreso %d/%d debería estar completed=%v pero está %v",
					tt.progress, tt.target, tt.expectedComplete, actualComplete)
			}
		})
	}
}

// TestAchievementBasicProperties valida las propiedades básicas de achievements
func TestAchievementBasicProperties(t *testing.T) {
	achievement := &Achievement{
		ID:          "test_1",
		UserID:      "user_1",
		Type:        "ai_partner",
		Name:        "🤖 AI Explorer",
		Description: "Utiliza 10 insights de IA",
		Points:      100,
		Progress:    0,
		Target:      10,
		Completed:   false,
	}

	// Test basic properties
	if achievement.UserID != "user_1" {
		t.Errorf("UserID debería ser 'user_1' pero es '%s'", achievement.UserID)
	}

	if achievement.Type != "ai_partner" {
		t.Errorf("Type debería ser 'ai_partner' pero es '%s'", achievement.Type)
	}

	if achievement.Target != 10 {
		t.Errorf("Target debería ser 10 pero es %d", achievement.Target)
	}

	if achievement.Completed {
		t.Error("Achievement nuevo debería empezar no completado")
	}
}

// TestAchievementProgressUpdate valida la actualización de progreso
func TestAchievementProgressUpdate(t *testing.T) {
	achievement := &Achievement{
		ID:        "test_ach_1",
		UserID:    "test_user",
		Type:      "ai_partner",
		Name:      "🤖 AI Explorer",
		Progress:  5,
		Target:    10,
		Completed: false,
	}

	// Before update
	if achievement.Completed {
		t.Error("Achievement debería empezar no completado")
	}

	// Update progress using the existing method
	achievement.UpdateProgress(10)

	// After update
	if achievement.Progress != 10 {
		t.Errorf("Progress debería ser 10 pero es %d", achievement.Progress)
	}

	if !achievement.Completed {
		t.Error("Achievement debería estar completado después de llegar al target")
	}
}

// TestMultipleAchievementProgress valida actualizaciones incrementales de progreso
func TestMultipleAchievementProgress(t *testing.T) {
	achievement := &Achievement{
		ID:       "test_ach_3",
		UserID:   "test_user",
		Type:     "action_taker",
		Progress: 0,
		Target:   25,
	}

	// Test incremental progress
	progressSteps := []int{5, 15, 20, 25, 30}

	for i, newProgress := range progressSteps {
		t.Run(fmt.Sprintf("step_%d_progress_%d", i+1, newProgress), func(t *testing.T) {
			achievement.UpdateProgress(newProgress)

			if achievement.Progress != newProgress {
				t.Errorf("Progress debería ser %d pero es %d", newProgress, achievement.Progress)
			}

			expectedComplete := newProgress >= achievement.Target
			if achievement.Completed != expectedComplete {
				t.Errorf("Completed debería ser %v para progreso %d/%d",
					expectedComplete, newProgress, achievement.Target)
			}
		})
	}
}
