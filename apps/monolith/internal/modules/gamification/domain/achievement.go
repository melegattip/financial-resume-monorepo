package domain

import (
	"time"

	"github.com/google/uuid"
)

// Achievement represents a single gamification achievement for a user.
type Achievement struct {
	ID          string
	UserID      string
	Type        string
	Name        string
	Description string
	Points      int
	Progress    int
	Target      int
	Completed   bool
	UnlockedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

// DefaultAchievements returns the 8 standard achievements every new user starts with.
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
