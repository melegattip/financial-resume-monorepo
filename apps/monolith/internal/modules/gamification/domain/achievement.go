package domain

import (
	"time"

	"github.com/google/uuid"
)

// Achievement represents a single gamification achievement for a user.
type Achievement struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Points      int        `json:"points"`
	Progress    int        `json:"progress"`
	Target      int        `json:"target"`
	Completed   bool       `json:"completed"`
	UnlockedAt  *time.Time `json:"unlocked_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// IsCompleted returns true when the achievement's progress has reached its target.
func (a *Achievement) IsCompleted() bool {
	return a.Progress >= a.Target
}

// UpdateProgress sets a new progress value. If the achievement reaches its target
// for the first time it is marked as completed with the current timestamp.
func (a *Achievement) UpdateProgress(newProgress int) {
	a.Progress = newProgress
	if a.Progress >= a.Target && !a.Completed {
		a.Completed = true
		now := time.Now().UTC()
		a.UnlockedAt = &now
	}
	a.UpdatedAt = time.Now().UTC()
}

// DefaultAchievements returns the standard achievements every new user starts with.
func DefaultAchievements(userID string) []Achievement {
	now := time.Now().UTC()
	defs := []struct {
		achievementType string
		name            string
		description     string
		target          int
		points          int
	}{
		{"transaction_starter", "🌱 Primer Paso", "Registra tu primera transacción", 1, 25},
		{"transaction_apprentice", "📝 Aprendiz Financiero", "Registra 10 transacciones", 10, 50},
		{"transaction_master", "💎 Maestro de Transacciones", "Registra 100 transacciones", 100, 200},
		{"category_creator", "🎨 Creador de Categorías", "Crea 5 categorías personalizadas", 5, 75},
		{"organization_expert", "📊 Expert en Organización", "Asigna categorías a 50 transacciones", 50, 100},
		{"weekly_warrior", "⚡ Guerrero Semanal", "Mantén una racha de 7 días consecutivos", 7, 100},
		{"monthly_legend", "👑 Leyenda Mensual", "Mantén una racha de 30 días consecutivos", 30, 500},
		{"data_explorer", "🔍 Explorador de Datos", "Usa la app durante 25 días", 25, 75},
		{"savings_starter", "🐖 Primer Ahorro", "Realiza tu primer depósito a una meta", 1, 50},
		{"savings_champion", "🏆 Campeón del Ahorro", "Completa una meta de ahorro", 1, 300},
		{"planner_pro", "🗓️ Planificador Pro", "Configura 3 transacciones recurrentes", 3, 150},
		{"budget_beginner", "📋 Primer Presupuesto", "Crea tu primer presupuesto", 1, 50},
		{"budget_disciplined", "💪 Disciplina Presupuestaria", "Cumple el presupuesto 3 meses seguidos", 3, 200},
		{"ai_executor", "🤖 Ejecutor de IA", "Aplica 5 recomendaciones de la IA", 5, 100},
	}

	achievements := make([]Achievement, 0, len(defs))
	for _, d := range defs {
		achievements = append(achievements, Achievement{
			ID:          uuid.New().String(),
			UserID:      userID,
			Type:        d.achievementType,
			Name:        d.name,
			Description: d.description,
			Points:      d.points,
			Progress:    0,
			Target:      d.target,
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return achievements
}
