package domain

import (
	"fmt"
	"testing"
)

// TestLevelCalculationCurrentImplementation valida que los niveles
// se calculen según la implementación actual del sistema
func TestLevelCalculationCurrentImplementation(t *testing.T) {
	// Tabla de niveles según la implementación actual
	// Basado en los thresholds reales: [0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500]
	tests := []struct {
		xp            int
		expectedLevel int
		description   string
	}{
		{xp: 0, expectedLevel: 1, description: "Nivel inicial"},
		{xp: 74, expectedLevel: 1, description: "Antes del threshold nivel 2"},
		{xp: 75, expectedLevel: 2, description: "Threshold nivel 2"},
		{xp: 199, expectedLevel: 2, description: "Antes del threshold nivel 3"},
		{xp: 200, expectedLevel: 3, description: "Threshold nivel 3"},
		{xp: 399, expectedLevel: 3, description: "Antes del threshold nivel 4"},
		{xp: 400, expectedLevel: 4, description: "Threshold nivel 4"},
		{xp: 5499, expectedLevel: 9, description: "Antes del nivel máximo"},
		{xp: 5500, expectedLevel: 10, description: "Nivel máximo"},
		{xp: 10000, expectedLevel: 10, description: "Por encima del máximo"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("xp_%d_should_be_level_%d", tt.xp, tt.expectedLevel), func(t *testing.T) {
			// Create user gamification with specific XP
			userGamification := &UserGamification{
				UserID:  "test_user",
				TotalXP: tt.xp,
			}

			// Calculate level
			actualLevel := userGamification.CalculateLevel()

			// Assert
			if actualLevel != tt.expectedLevel {
				t.Errorf("%s: Usuario con %d XP debería estar en nivel %d, pero está en nivel %d",
					tt.description, tt.xp, tt.expectedLevel, actualLevel)
			}
		})
	}
}

// TestXPToNextLevelCalculation valida el cálculo de XP necesario para próximo nivel
func TestXPToNextLevelCalculation(t *testing.T) {
	tests := []struct {
		currentXP        int
		expectedXPToNext int
		description      string
	}{
		{currentXP: 0, expectedXPToNext: 75, description: "Nivel 1 -> 2"},
		{currentXP: 50, expectedXPToNext: 25, description: "Nivel 1 -> 2 (faltan 25)"},
		{currentXP: 75, expectedXPToNext: 125, description: "Nivel 2 -> 3"},
		{currentXP: 200, expectedXPToNext: 200, description: "Nivel 3 -> 4"},
		{currentXP: 5500, expectedXPToNext: 0, description: "Nivel máximo (10)"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("xp_%d_needs_%d_to_next", tt.currentXP, tt.expectedXPToNext), func(t *testing.T) {
			userGamification := &UserGamification{
				UserID:  "test_user",
				TotalXP: tt.currentXP,
			}

			actualXPToNext := userGamification.XPToNextLevel()

			if actualXPToNext != tt.expectedXPToNext {
				t.Errorf("%s: Usuario con %d XP debería necesitar %d XP para próximo nivel, pero necesita %d",
					tt.description, tt.currentXP, tt.expectedXPToNext, actualXPToNext)
			}
		})
	}
}

// TestProgressToNextLevelCalculation valida el cálculo de porcentaje de progreso
func TestProgressToNextLevelCalculation(t *testing.T) {
	tests := []struct {
		currentXP       int
		expectedPercent int
		description     string
	}{
		{currentXP: 0, expectedPercent: 0, description: "Inicio nivel 1"},
		{currentXP: 37, expectedPercent: 49, description: "Mitad nivel 1"},
		{currentXP: 75, expectedPercent: 0, description: "Inicio nivel 2"},
		{currentXP: 137, expectedPercent: 49, description: "Mitad nivel 2"},
		{currentXP: 5500, expectedPercent: 100, description: "Nivel máximo"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("xp_%d_progress_%d", tt.currentXP, tt.expectedPercent), func(t *testing.T) {
			userGamification := &UserGamification{
				UserID:  "test_user",
				TotalXP: tt.currentXP,
			}

			actualPercent := userGamification.ProgressToNextLevel()

			if actualPercent != tt.expectedPercent {
				t.Errorf("%s: Usuario con %d XP debería tener %d%% de progreso, pero tiene %d%%",
					tt.description, tt.currentXP, tt.expectedPercent, actualPercent)
			}
		})
	}
}

// TestLevelThresholds valida que los umbrales de niveles sean los correctos del sistema actual
func TestLevelThresholds(t *testing.T) {
	tests := []struct {
		xp            int
		expectedLevel int
	}{
		{0, 1},     // 0 XP = Nivel 1
		{74, 1},    // 74 XP = Nivel 1
		{75, 2},    // 75 XP = Nivel 2
		{199, 2},   // 199 XP = Nivel 2
		{200, 3},   // 200 XP = Nivel 3
		{5500, 10}, // 5500 XP = Nivel 10 (máximo)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("xp_%d_level_%d", tt.xp, tt.expectedLevel), func(t *testing.T) {
			userGamification := &UserGamification{
				UserID:  "test_user",
				TotalXP: tt.xp,
			}
			actualLevel := userGamification.CalculateLevel()
			if actualLevel != tt.expectedLevel {
				t.Errorf("XP %d debería resultar en nivel %d pero resultó en nivel %d",
					tt.xp, tt.expectedLevel, actualLevel)
			}
		})
	}
}

// TestAchievementCompletion valida la lógica de completado de achievements
func TestAchievementCompletion(t *testing.T) {
	tests := []struct {
		name     string
		progress int
		target   int
		expected bool
	}{
		{"not_completed_partial", 5, 10, false},
		{"completed_exact", 10, 10, true},
		{"completed_exceeded", 15, 10, true},
		{"zero_progress", 0, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achievement := &Achievement{
				Progress: tt.progress,
				Target:   tt.target,
			}

			actual := achievement.IsCompleted()
			if actual != tt.expected {
				t.Errorf("Achievement con progreso %d/%d debería IsCompleted()=%v pero fue %v",
					tt.progress, tt.target, tt.expected, actual)
			}
		})
	}
}
