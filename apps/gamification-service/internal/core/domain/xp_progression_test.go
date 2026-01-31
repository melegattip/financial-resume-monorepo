package domain

import (
	"fmt"
	"testing"
)

// TestXPAccumulation valida la acumulación progresiva de XP
func TestXPAccumulation(t *testing.T) {
	user := &UserGamification{
		UserID:  "test_user_accumulation",
		TotalXP: 0,
	}

	// Test de acumulación paso a paso
	tests := []struct {
		name          string
		addXP         int
		expectedTotal int
		expectedLevel int
		description   string
	}{
		{"initial", 0, 0, 1, "Estado inicial"},
		{"first_action", 8, 8, 1, "Primera acción (create_expense)"},
		{"second_action", 10, 18, 1, "Segunda acción (create_category)"},
		{"third_action", 2, 20, 1, "Tercera acción (view_dashboard)"},
		{"accumulating", 30, 50, 1, "Acumulando más XP"},
		{"almost_level_2", 24, 74, 1, "Casi nivel 2"},
		{"level_up_to_2", 1, 75, 2, "Subida a nivel 2"},
		{"more_xp_level_2", 100, 175, 2, "Más XP en nivel 2"},
		{"level_up_to_3", 25, 200, 3, "Subida a nivel 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Agregar XP
			user.TotalXP += tt.addXP

			// Recalcular nivel
			user.CurrentLevel = user.CalculateLevel()

			// Verificar total XP
			if user.TotalXP != tt.expectedTotal {
				t.Errorf("%s: XP total esperado %d, obtenido %d",
					tt.description, tt.expectedTotal, user.TotalXP)
			}

			// Verificar nivel
			if user.CurrentLevel != tt.expectedLevel {
				t.Errorf("%s: Nivel esperado %d, obtenido %d",
					tt.description, tt.expectedLevel, user.CurrentLevel)
			}
		})
	}
}

// TestXPProgressionThroughLevels valida la progresión completa a través de niveles
func TestXPProgressionThroughLevels(t *testing.T) {
	user := &UserGamification{
		UserID:  "test_user_progression",
		TotalXP: 0,
	}

	// Simular sesiones de usuario con diferentes cantidades de acciones
	sessions := []struct {
		name          string
		actions       []int // XP por cada acción en la sesión
		expectedLevel int
		expectedXP    int
		description   string
	}{
		{
			name:          "session_1_basic_actions",
			actions:       []int{2, 8, 10}, // dashboard + expense + category
			expectedLevel: 1,
			expectedXP:    20,
			description:   "Sesión 1: Acciones básicas",
		},
		{
			name:          "session_2_more_transactions",
			actions:       []int{8, 8, 8, 8, 8}, // 5 transacciones
			expectedLevel: 1,
			expectedXP:    60, // 20 + 40 = 60
			description:   "Sesión 2: Más transacciones",
		},
		{
			name:          "session_3_push_to_level_2",
			actions:       []int{8, 5, 2}, // expense + update + dashboard
			expectedLevel: 2,
			expectedXP:    75, // 60 + 15 = 75 (nivel 2)
			description:   "Sesión 3: Alcanzar nivel 2",
		},
		{
			name:          "session_4_level_2_actions",
			actions:       []int{10, 10, 20, 15}, // categories + budget + savings
			expectedLevel: 2,
			expectedXP:    130, // 75 + 55 = 130
			description:   "Sesión 4: Acciones en nivel 2",
		},
		{
			name:          "session_5_push_to_level_3",
			actions:       []int{20, 20, 30}, // budgets and advanced actions
			expectedLevel: 3,
			expectedXP:    200, // 130 + 70 = 200 (nivel 3)
			description:   "Sesión 5: Alcanzar nivel 3",
		},
	}

	for _, session := range sessions {
		t.Run(session.name, func(t *testing.T) {
			// Simular acciones de la sesión
			sessionXP := 0
			for _, actionXP := range session.actions {
				user.TotalXP += actionXP
				sessionXP += actionXP
			}

			// Recalcular nivel
			user.CurrentLevel = user.CalculateLevel()

			// Verificar XP total
			if user.TotalXP != session.expectedXP {
				t.Errorf("%s: XP total esperado %d, obtenido %d",
					session.description, session.expectedXP, user.TotalXP)
			}

			// Verificar nivel
			if user.CurrentLevel != session.expectedLevel {
				t.Errorf("%s: Nivel esperado %d, obtenido %d",
					session.description, session.expectedLevel, user.CurrentLevel)
			}

			t.Logf("%s: Ganó %d XP, Total: %d XP, Nivel: %d",
				session.description, sessionXP, user.TotalXP, user.CurrentLevel)
		})
	}
}

// TestLevelUpDetection valida la detección de subidas de nivel
func TestLevelUpDetection(t *testing.T) {
	tests := []struct {
		name          string
		startXP       int
		addXP         int
		expectLevelUp bool
		oldLevel      int
		newLevel      int
		description   string
	}{
		{
			name:          "no_level_up",
			startXP:       50,
			addXP:         20,
			expectLevelUp: false,
			oldLevel:      1,
			newLevel:      1,
			description:   "XP 50->70: Sin cambio de nivel",
		},
		{
			name:          "level_up_1_to_2",
			startXP:       70,
			addXP:         10,
			expectLevelUp: true,
			oldLevel:      1,
			newLevel:      2,
			description:   "XP 70->80: Subida nivel 1->2",
		},
		{
			name:          "level_up_2_to_3",
			startXP:       190,
			addXP:         20,
			expectLevelUp: true,
			oldLevel:      2,
			newLevel:      3,
			description:   "XP 190->210: Subida nivel 2->3",
		},
		{
			name:          "big_jump_multiple_levels",
			startXP:       50,
			addXP:         400,
			expectLevelUp: true,
			oldLevel:      1,
			newLevel:      4,
			description:   "XP 50->450: Salto múltiples niveles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &UserGamification{
				UserID:  "test_user_levelup",
				TotalXP: tt.startXP,
			}

			// Calcular nivel inicial
			oldLevel := user.CalculateLevel()

			// Verificar nivel inicial esperado
			if oldLevel != tt.oldLevel {
				t.Errorf("%s: Nivel inicial esperado %d, obtenido %d",
					tt.description, tt.oldLevel, oldLevel)
			}

			// Agregar XP
			user.TotalXP += tt.addXP
			newLevel := user.CalculateLevel()

			// Verificar nuevo nivel
			if newLevel != tt.newLevel {
				t.Errorf("%s: Nuevo nivel esperado %d, obtenido %d",
					tt.description, tt.newLevel, newLevel)
			}

			// Verificar detección de level up
			actualLevelUp := newLevel > oldLevel
			if actualLevelUp != tt.expectLevelUp {
				t.Errorf("%s: Level up esperado %v, obtenido %v",
					tt.description, tt.expectLevelUp, actualLevelUp)
			}

			t.Logf("%s: %d + %d XP = %d XP, Nivel %d->%d, LevelUp: %v",
				tt.description, tt.startXP, tt.addXP, user.TotalXP, oldLevel, newLevel, actualLevelUp)
		})
	}
}

// TestXPToNextLevelProgression valida el cálculo de XP restante durante progresión
func TestXPToNextLevelProgression(t *testing.T) {
	user := &UserGamification{
		UserID:  "test_user_xp_to_next",
		TotalXP: 0,
	}

	progressionSteps := []struct {
		addXP          int
		expectedXPLeft int
		expectedLevel  int
		description    string
	}{
		{0, 75, 1, "Inicial: 0 XP, faltan 75 para nivel 2"},
		{25, 50, 1, "25 XP: faltan 50 para nivel 2"},
		{25, 25, 1, "50 XP: faltan 25 para nivel 2"},
		{25, 125, 2, "75 XP: Nivel 2, faltan 125 para nivel 3"},
		{50, 75, 2, "125 XP: faltan 75 para nivel 3"},
		{75, 200, 3, "200 XP: Nivel 3, faltan 200 para nivel 4"},
		{200, 300, 4, "400 XP: Nivel 4, faltan 300 para nivel 5"},
	}

	for i, step := range progressionSteps {
		t.Run(fmt.Sprintf("step_%d", i+1), func(t *testing.T) {
			// Agregar XP
			user.TotalXP += step.addXP

			// Recalcular valores
			level := user.CalculateLevel()
			xpToNext := user.XPToNextLevel()

			// Verificar nivel
			if level != step.expectedLevel {
				t.Errorf("%s: Nivel esperado %d, obtenido %d",
					step.description, step.expectedLevel, level)
			}

			// Verificar XP restante
			if xpToNext != step.expectedXPLeft {
				t.Errorf("%s: XP restante esperado %d, obtenido %d",
					step.description, step.expectedXPLeft, xpToNext)
			}

			t.Logf("Step %d: %s -> Total XP: %d, Nivel: %d, XP restante: %d",
				i+1, step.description, user.TotalXP, level, xpToNext)
		})
	}
}

// TestProgressPercentageCalculation valida el cálculo de porcentaje de progreso
func TestProgressPercentageCalculation(t *testing.T) {
	tests := []struct {
		totalXP         int
		expectedPercent int
		expectedLevel   int
		description     string
	}{
		{0, 0, 1, "Inicio nivel 1: 0%"},
		{37, 49, 1, "Mitad nivel 1: ~49%"},
		{74, 98, 1, "Casi nivel 2: ~98%"},
		{75, 0, 2, "Inicio nivel 2: 0%"},
		{137, 49, 2, "Mitad nivel 2: ~49%"},
		{199, 99, 2, "Casi nivel 3: ~99%"},
		{200, 0, 3, "Inicio nivel 3: 0%"},
		{5500, 100, 10, "Nivel máximo: 100%"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("xp_%d", tt.totalXP), func(t *testing.T) {
			user := &UserGamification{
				UserID:  "test_user_progress",
				TotalXP: tt.totalXP,
			}

			level := user.CalculateLevel()
			progress := user.ProgressToNextLevel()

			// Verificar nivel
			if level != tt.expectedLevel {
				t.Errorf("%s: Nivel esperado %d, obtenido %d",
					tt.description, tt.expectedLevel, level)
			}

			// Verificar porcentaje (con tolerancia de ±2%)
			diff := progress - tt.expectedPercent
			if diff < -2 || diff > 2 {
				t.Errorf("%s: Progreso esperado ~%d%%, obtenido %d%%",
					tt.description, tt.expectedPercent, progress)
			}

			t.Logf("%s: %d XP = Nivel %d, Progreso %d%%",
				tt.description, tt.totalXP, level, progress)
		})
	}
}
