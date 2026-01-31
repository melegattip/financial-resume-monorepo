package ports

import (
	"context"
	"errors"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// Errores específicos de gamificación
var (
	ErrGamificationNotFound = errors.New("gamification not found")
	ErrAchievementNotFound  = errors.New("achievement not found")
	ErrActionNotFound       = errors.New("action not found")
)

// GamificationRepository define las operaciones de persistencia para gamificación
type GamificationRepository interface {
	// UserGamification operations
	Create(ctx context.Context, gamification *domain.UserGamification) error
	GetByUserID(ctx context.Context, userID string) (*domain.UserGamification, error)
	Update(ctx context.Context, gamification *domain.UserGamification) error
	Delete(ctx context.Context, userID string) error

	// Achievement operations
	CreateAchievement(ctx context.Context, achievement *domain.Achievement) error
	GetAchievementsByUserID(ctx context.Context, userID string) ([]domain.Achievement, error)
	GetAchievementByID(ctx context.Context, achievementID string) (*domain.Achievement, error)
	UpdateAchievement(ctx context.Context, achievement *domain.Achievement) error
	DeleteAchievement(ctx context.Context, achievementID string) error

	// UserAction operations
	CreateAction(ctx context.Context, action *domain.UserAction) error
	GetActionsByUserID(ctx context.Context, userID string) ([]domain.UserAction, error)
	GetActionsByUserIDAndPeriod(ctx context.Context, userID string, startDate, endDate string) ([]domain.UserAction, error)
	DeleteAction(ctx context.Context, actionID string) error
}
