package domain

import (
	"time"
)

// UserGamification representa el estado de gamificación de un usuario
type UserGamification struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	TotalXP           int       `json:"total_xp"`
	CurrentLevel      int       `json:"current_level"`
	InsightsViewed    int       `json:"insights_viewed"`
	ActionsCompleted  int       `json:"actions_completed"`
	AchievementsCount int       `json:"achievements_count"`
	CurrentStreak     int       `json:"current_streak"`
	LastActivity      time.Time `json:"last_activity"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Achievement representa un logro desbloqueado
type Achievement struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Type        string     `json:"type"`        // "ai_partner", "action_taker", etc.
	Name        string     `json:"name"`        // "🤖 AI Partner"
	Description string     `json:"description"` // "100 insights de IA utilizados"
	Points      int        `json:"points"`      // XP otorgados
	Progress    int        `json:"progress"`    // Progreso actual (ej: 67/100)
	Target      int        `json:"target"`      // Objetivo a alcanzar (ej: 100)
	Completed   bool       `json:"completed"`
	UnlockedAt  *time.Time `json:"unlocked_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserAction representa una acción del usuario que otorga XP
type UserAction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ActionType  string    `json:"action_type"` // "view_insight", "understand_insight", "complete_action"
	EntityType  string    `json:"entity_type"` // "insight", "suggestion", "pattern"
	EntityID    string    `json:"entity_id"`   // ID del insight/suggestion específico
	XPEarned    int       `json:"xp_earned"`
	Description string    `json:"description"` // "Viewed insight: Gastos Elevados"
	CreatedAt   time.Time `json:"created_at"`
}

// GamificationStats representa estadísticas agregadas de gamificación
type GamificationStats struct {
	UserID                string    `json:"user_id"`
	TotalXP               int       `json:"total_xp"`
	CurrentLevel          int       `json:"current_level"`
	XPToNextLevel         int       `json:"xp_to_next_level"`
	ProgressPercent       int       `json:"progress_percent"`
	TotalAchievements     int       `json:"total_achievements"`
	CompletedAchievements int       `json:"completed_achievements"`
	CurrentStreak         int       `json:"current_streak"`
	LastActivity          time.Time `json:"last_activity"`
}

// ActionType constants
const (
	ActionTypeViewInsight       = "view_insight"
	ActionTypeUnderstandInsight = "understand_insight"
	ActionTypeCompleteAction    = "complete_action"
	ActionTypeViewPattern       = "view_pattern"
	ActionTypeUseSuggestion     = "use_suggestion"
)

// AchievementType constants
const (
	AchievementTypeAIPartner     = "ai_partner"
	AchievementTypeActionTaker   = "action_taker"
	AchievementTypeDataExplorer  = "data_explorer"
	AchievementTypeQuickLearner  = "quick_learner"
	AchievementTypeInsightMaster = "insight_master"
	AchievementTypeStreakKeeper  = "streak_keeper"
)

// EntityType constants
const (
	EntityTypeInsight    = "insight"
	EntityTypeSuggestion = "suggestion"
	EntityTypePattern    = "pattern"
)

// CalculateLevel calcula el nivel basado en XP total
func (ug *UserGamification) CalculateLevel() int {
	levels := []int{0, 100, 250, 500, 1000, 2000, 4000, 8000, 16000, 32000}

	for i := len(levels) - 1; i >= 0; i-- {
		if ug.TotalXP >= levels[i] {
			return i
		}
	}
	return 0
}

// XPToNextLevel calcula XP necesario para el siguiente nivel
func (ug *UserGamification) XPToNextLevel() int {
	levels := []int{0, 100, 250, 500, 1000, 2000, 4000, 8000, 16000, 32000}
	currentLevel := ug.CalculateLevel()

	if currentLevel >= len(levels)-1 {
		return 0 // Nivel máximo alcanzado
	}

	return levels[currentLevel+1] - ug.TotalXP
}

// ProgressToNextLevel calcula el porcentaje de progreso al siguiente nivel
func (ug *UserGamification) ProgressToNextLevel() int {
	levels := []int{0, 100, 250, 500, 1000, 2000, 4000, 8000, 16000, 32000}
	currentLevel := ug.CalculateLevel()

	if currentLevel >= len(levels)-1 {
		return 100 // Nivel máximo
	}

	currentLevelXP := levels[currentLevel]
	nextLevelXP := levels[currentLevel+1]
	progressXP := ug.TotalXP - currentLevelXP

	return int((float64(progressXP) / float64(nextLevelXP-currentLevelXP)) * 100)
}

// IsCompleted verifica si un achievement está completado
func (a *Achievement) IsCompleted() bool {
	return a.Progress >= a.Target
}

// UpdateProgress actualiza el progreso de un achievement
func (a *Achievement) UpdateProgress(newProgress int) {
	a.Progress = newProgress
	if a.Progress >= a.Target && !a.Completed {
		a.Completed = true
		now := time.Now()
		a.UnlockedAt = &now
	}
	a.UpdatedAt = time.Now()
}

// GetLevelName retorna el nombre del nivel actual
func (ug *UserGamification) GetLevelName() string {
	levelNames := []string{
		"Financial Newbie",
		"Money Aware",
		"Budget Tracker",
		"Savings Starter",
		"Financial Explorer",
		"Money Manager",
		"Investment Learner",
		"Financial Guru",
		"Money Master",
		"Financial Magnate",
	}

	level := ug.CalculateLevel()
	if level >= len(levelNames) {
		return "Financial Legend"
	}

	return levelNames[level]
}
