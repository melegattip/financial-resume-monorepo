package ports

import (
	"context"
	"errors"
	"time"

	"github.com/melegattip/financial-gamification-service/internal/core/domain"
)

// Errores personalizados
var (
	ErrGamificationNotFound = errors.New("gamification record not found")
	ErrAchievementNotFound  = errors.New("achievement not found")
	ErrActionNotFound       = errors.New("action not found")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrInvalidActionType    = errors.New("invalid action type")
)

// GamificationRepository define las operaciones de persistencia para gamificación
type GamificationRepository interface {
	// User gamification operations
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

	// Action operations
	CreateAction(ctx context.Context, action *domain.UserAction) error
	GetActionsByUserID(ctx context.Context, userID string) ([]domain.UserAction, error)
	GetActionsByUserIDAndPeriod(ctx context.Context, userID string, startDate, endDate string) ([]domain.UserAction, error)
	DeleteAction(ctx context.Context, actionID string) error

	// Challenge operations
	GetActiveChallenges(ctx context.Context, challengeType string) ([]domain.Challenge, error)
	GetChallengeByKey(ctx context.Context, challengeKey string) (*domain.Challenge, error)

	// User challenge operations
	GetUserChallengesForDate(ctx context.Context, userID string, challengeDate time.Time, challengeType string) ([]domain.UserChallenge, error)
	CreateOrUpdateUserChallenge(ctx context.Context, userChallenge *domain.UserChallenge) error
	GetUserChallengeByID(ctx context.Context, userChallengeID string) (*domain.UserChallenge, error)

	// Challenge progress tracking operations
	UpdateChallengeProgressTracking(ctx context.Context, tracking *domain.ChallengeProgressTracking) error
	GetChallengeProgressTracking(ctx context.Context, userID string, challengeDate time.Time) ([]domain.ChallengeProgressTracking, error)
}
