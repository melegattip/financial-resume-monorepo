package domain

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

// Contador global para garantizar unicidad en generateUUID
var uuidCounter int64

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

// NewID genera un nuevo ID único
func NewID() string {
	// En un microservicio real, usaríamos UUID
	return generateUUID()
}

// CalculateLevel calcula el nivel basado en XP total
func (ug *UserGamification) CalculateLevel() int {
	levels := []int{0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500}

	for i := len(levels) - 1; i >= 0; i-- {
		if ug.TotalXP >= levels[i] {
			return i + 1 // Devolver nivel 1-10 en lugar de 0-9
		}
	}
	return 1 // Nivel mínimo es 1
}

// XPToNextLevel calcula XP necesario para el siguiente nivel
func (ug *UserGamification) XPToNextLevel() int {
	levels := []int{0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500}
	currentLevel := ug.CalculateLevel()

	if currentLevel >= 10 { // Nivel máximo es 10
		return 0 // Nivel máximo alcanzado
	}

	return levels[currentLevel] - ug.TotalXP
}

// ProgressToNextLevel calcula el porcentaje de progreso al siguiente nivel
func (ug *UserGamification) ProgressToNextLevel() int {
	levels := []int{0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500}
	currentLevel := ug.CalculateLevel()

	if currentLevel >= 10 { // Nivel máximo es 10
		return 100 // Nivel máximo
	}

	currentLevelXP := levels[currentLevel-1] // Ajuste por el offset de nivel
	nextLevelXP := levels[currentLevel]
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
		"", // Índice 0 vacío
		"Financial Newbie",
		"Money Tracker",
		"Smart Saver", // 🔓 METAS DE AHORRO
		"Budget Master",
		"Financial Planner", // 🔓 PRESUPUESTOS
		"Investment Seeker",
		"Wealth Builder", // 🔓 IA FINANCIERA
		"Financial Strategist",
		"Money Mentor",
		"Financial Magnate",
	}

	level := ug.CalculateLevel()
	if level >= len(levelNames) {
		return "Financial Legend"
	}

	return levelNames[level]
}

// Challenge representa un challenge disponible
type Challenge struct {
	ID                string                 `json:"id"`
	ChallengeKey      string                 `json:"challenge_key"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	ChallengeType     string                 `json:"challenge_type"` // "daily", "weekly", "monthly"
	Icon              string                 `json:"icon"`
	XPReward          int                    `json:"xp_reward"`
	RequirementType   string                 `json:"requirement_type"` // "transaction_count", "category_variety", etc.
	RequirementTarget int                    `json:"requirement_target"`
	RequirementData   map[string]interface{} `json:"requirement_data"` // Requisitos complejos
	Active            bool                   `json:"active"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// UserChallenge representa el progreso de un usuario en un challenge
type UserChallenge struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	ChallengeID   string     `json:"challenge_id"`
	ChallengeDate time.Time  `json:"challenge_date"` // Para challenges diarios/semanales
	Progress      int        `json:"progress"`
	Target        int        `json:"target"`
	Completed     bool       `json:"completed"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Datos del challenge (join)
	Challenge *Challenge `json:"challenge,omitempty"`
}

// ChallengeProgressTracking representa el tracking detallado de acciones para challenges
type ChallengeProgressTracking struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"user_id"`
	ChallengeDate  time.Time              `json:"challenge_date"`
	ActionType     string                 `json:"action_type"`
	EntityType     string                 `json:"entity_type"`
	Count          int                    `json:"count"`
	UniqueEntities map[string]interface{} `json:"unique_entities"` // Para tracking de entities únicas
	CreatedAt      time.Time              `json:"created_at"`
}

// ChallengeResult representa el resultado de completar un challenge
type ChallengeResult struct {
	UserChallenge  *UserChallenge  `json:"user_challenge"`
	XPEarned       int             `json:"xp_earned"`
	NewlyCompleted bool            `json:"newly_completed"`
	AllChallenges  []UserChallenge `json:"all_challenges"`
}

// Challenge Types constants
const (
	ChallengeTypeDaily   = "daily"
	ChallengeTypeWeekly  = "weekly"
	ChallengeTypeMonthly = "monthly"
)

// Challenge Requirement Types constants
const (
	RequirementTypeTransactionCount = "transaction_count"
	RequirementTypeCategoryVariety  = "category_variety"
	RequirementTypeViewCombo        = "view_combo"
	RequirementTypeDailyLogin       = "daily_login"
	RequirementTypeDailyLoginCount  = "daily_login_count"
)

// Métodos helper para UserChallenge

// IsCompleted verifica si un challenge está completado
func (uc *UserChallenge) IsCompleted() bool {
	return uc.Progress >= uc.Target
}

// UpdateProgress actualiza el progreso de un challenge
func (uc *UserChallenge) UpdateProgress(newProgress int) {
	uc.Progress = newProgress
	if uc.Progress >= uc.Target && !uc.Completed {
		uc.Completed = true
		now := time.Now()
		uc.CompletedAt = &now
	}
	uc.UpdatedAt = time.Now()
}

// GetProgressPercentage retorna el porcentaje de progreso
func (uc *UserChallenge) GetProgressPercentage() int {
	if uc.Target == 0 {
		return 0
	}
	return min(100, (uc.Progress*100)/uc.Target)
}

// GetChallengeDate retorna la fecha del challenge (solo fecha, sin hora)
func GetChallengeDate(t time.Time, challengeType string) time.Time {
	switch challengeType {
	case ChallengeTypeDaily:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case ChallengeTypeWeekly:
		// Obtener el lunes de la semana
		weekday := int(t.Weekday())
		if weekday == 0 { // Domingo = 0, queremos que sea 7
			weekday = 7
		}
		daysToMonday := weekday - 1
		monday := t.AddDate(0, 0, -daysToMonday)
		return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
	case ChallengeTypeMonthly:
		// Primer día del mes
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	default:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}
}

// generateUUID genera un UUID único usando nanosegundos y un contador
func generateUUID() string {
	ns := time.Now().UnixNano()
	count := atomic.AddInt64(&uuidCounter, 1)
	random := rand.Intn(10000)

	return fmt.Sprintf("gam-%d-%d-%d", ns, count, random)
}

// min función helper
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
