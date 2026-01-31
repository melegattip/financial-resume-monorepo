package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// GamificationUseCase define los casos de uso de gamificación
type GamificationUseCase interface {
	// User gamification management
	GetUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error)
	InitializeUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error)

	// Action tracking
	RecordUserAction(ctx context.Context, params RecordActionParams) (*ActionResult, error)

	// Achievement management
	GetUserAchievements(ctx context.Context, userID string) ([]domain.Achievement, error)
	CheckAndUpdateAchievements(ctx context.Context, userID string) ([]domain.Achievement, []domain.Achievement, error)

	// Statistics
	GetGamificationStats(ctx context.Context, userID string) (*domain.GamificationStats, error)
}

// RecordActionParams parámetros para registrar una acción del usuario
type RecordActionParams struct {
	UserID      string `json:"user_id"`
	ActionType  string `json:"action_type"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Description string `json:"description"`
}

// ActionResult resultado de registrar una acción
type ActionResult struct {
	XPEarned            int                  `json:"xp_earned"`
	NewLevel            int                  `json:"new_level"`
	CurrentLevel        int                  `json:"current_level"` // Compatibilidad con frontend
	LevelUp             bool                 `json:"level_up"`
	NewAchievements     []domain.Achievement `json:"new_achievements"`
	UpdatedAchievements []domain.Achievement `json:"updated_achievements"`
	TotalXP             int                  `json:"total_xp"`
}

// gamificationService implementa GamificationUseCase
type gamificationService struct {
	gamificationRepo ports.GamificationRepository
}

// NewGamificationUseCase crea una nueva instancia del servicio de gamificación
func NewGamificationUseCase(gamificationRepo ports.GamificationRepository) GamificationUseCase {
	return &gamificationService{
		gamificationRepo: gamificationRepo,
	}
}

// GetUserGamification obtiene el estado de gamificación del usuario
func (s *gamificationService) GetUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error) {
	gamification, err := s.gamificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Si no existe, inicializar automáticamente
		if err == ports.ErrGamificationNotFound {
			return s.InitializeUserGamification(ctx, userID)
		}
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Actualizar nivel calculado
	gamification.CurrentLevel = gamification.CalculateLevel()

	return gamification, nil
}

// InitializeUserGamification inicializa la gamificación para un nuevo usuario
func (s *gamificationService) InitializeUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error) {
	gamification := &domain.UserGamification{
		ID:                domain.NewID(),
		UserID:            userID,
		TotalXP:           0,
		CurrentLevel:      0,
		InsightsViewed:    0,
		ActionsCompleted:  0,
		AchievementsCount: 0,
		CurrentStreak:     0,
		LastActivity:      time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err := s.gamificationRepo.Create(ctx, gamification)
	if err != nil {
		return nil, fmt.Errorf("error creating user gamification: %w", err)
	}

	// Inicializar achievements básicos
	err = s.initializeBasicAchievements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error initializing achievements: %w", err)
	}

	return gamification, nil
}

// RecordUserAction registra una acción del usuario y otorga XP
func (s *gamificationService) RecordUserAction(ctx context.Context, params RecordActionParams) (*ActionResult, error) {
	// Obtener gamificación actual del usuario
	gamification, err := s.GetUserGamification(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Calcular XP para la acción
	xpEarned := s.calculateXPForAction(params.ActionType, params.EntityType)

	// Registrar la acción
	action := &domain.UserAction{
		ID:          domain.NewID(),
		UserID:      params.UserID,
		ActionType:  params.ActionType,
		EntityType:  params.EntityType,
		EntityID:    params.EntityID,
		XPEarned:    xpEarned,
		Description: params.Description,
		CreatedAt:   time.Now(),
	}

	err = s.gamificationRepo.CreateAction(ctx, action)
	if err != nil {
		return nil, fmt.Errorf("error creating user action: %w", err)
	}

	// Actualizar estadísticas del usuario
	oldLevel := gamification.CurrentLevel
	gamification.TotalXP += xpEarned
	gamification.CurrentLevel = gamification.CalculateLevel()
	gamification.LastActivity = time.Now()
	gamification.UpdatedAt = time.Now()

	// Actualizar contadores específicos
	switch params.ActionType {
	case domain.ActionTypeViewInsight:
		gamification.InsightsViewed++
	case domain.ActionTypeCompleteAction:
		gamification.ActionsCompleted++
	}

	err = s.gamificationRepo.Update(ctx, gamification)
	if err != nil {
		return nil, fmt.Errorf("error updating user gamification: %w", err)
	}

	// Verificar nuevos achievements
	newAchievements, updatedAchievements, err := s.CheckAndUpdateAchievements(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("error checking achievements: %w", err)
	}

	result := &ActionResult{
		XPEarned:            xpEarned,
		NewLevel:            gamification.CurrentLevel,
		LevelUp:             gamification.CurrentLevel > oldLevel,
		NewAchievements:     newAchievements,
		UpdatedAchievements: updatedAchievements,
		TotalXP:             gamification.TotalXP,
	}

	return result, nil
}

// GetUserAchievements obtiene todos los achievements del usuario
func (s *gamificationService) GetUserAchievements(ctx context.Context, userID string) ([]domain.Achievement, error) {
	achievements, err := s.gamificationRepo.GetAchievementsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user achievements: %w", err)
	}

	return achievements, nil
}

// CheckAndUpdateAchievements verifica y actualiza achievements del usuario
func (s *gamificationService) CheckAndUpdateAchievements(ctx context.Context, userID string) ([]domain.Achievement, []domain.Achievement, error) {
	// Obtener gamificación actual
	gamification, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Obtener achievements actuales
	achievements, err := s.GetUserAchievements(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting achievements: %w", err)
	}

	var newAchievements []domain.Achievement
	var updatedAchievements []domain.Achievement

	// Verificar cada tipo de achievement
	for _, achievement := range achievements {
		oldProgress := achievement.Progress
		newProgress := s.calculateAchievementProgress(achievement.Type, gamification)

		if newProgress > oldProgress {
			achievement.UpdateProgress(newProgress)

			// Si se completó por primera vez, es un nuevo achievement
			if achievement.Completed && oldProgress < achievement.Target {
				newAchievements = append(newAchievements, achievement)
				gamification.AchievementsCount++
			} else {
				updatedAchievements = append(updatedAchievements, achievement)
			}

			// Actualizar en base de datos
			err = s.gamificationRepo.UpdateAchievement(ctx, &achievement)
			if err != nil {
				return nil, nil, fmt.Errorf("error updating achievement: %w", err)
			}
		}
	}

	// Actualizar contador de achievements si hay nuevos
	if len(newAchievements) > 0 {
		err = s.gamificationRepo.Update(ctx, gamification)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating gamification: %w", err)
		}
	}

	return newAchievements, updatedAchievements, nil
}

// GetGamificationStats obtiene estadísticas de gamificación del usuario
func (s *gamificationService) GetGamificationStats(ctx context.Context, userID string) (*domain.GamificationStats, error) {
	gamification, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	achievements, err := s.GetUserAchievements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting achievements: %w", err)
	}

	completedAchievements := 0
	for _, achievement := range achievements {
		if achievement.Completed {
			completedAchievements++
		}
	}

	stats := &domain.GamificationStats{
		UserID:                userID,
		TotalXP:               gamification.TotalXP,
		CurrentLevel:          gamification.CurrentLevel,
		XPToNextLevel:         gamification.XPToNextLevel(),
		ProgressPercent:       gamification.ProgressToNextLevel(),
		TotalAchievements:     len(achievements),
		CompletedAchievements: completedAchievements,
		CurrentStreak:         gamification.CurrentStreak,
		LastActivity:          gamification.LastActivity,
	}

	return stats, nil
}

// calculateXPForAction calcula XP basado en el tipo de acción
func (s *gamificationService) calculateXPForAction(actionType, entityType string) int {
	basePoints := map[string]int{
		domain.ActionTypeViewInsight:       1,
		domain.ActionTypeUnderstandInsight: 3,
		domain.ActionTypeCompleteAction:    10,
		domain.ActionTypeViewPattern:       2,
		domain.ActionTypeUseSuggestion:     5,
	}

	// Multiplicadores por tipo de entidad
	multipliers := map[string]float64{
		domain.EntityTypeInsight:    1.0,
		domain.EntityTypeSuggestion: 1.2,
		domain.EntityTypePattern:    1.1,
	}

	baseXP := basePoints[actionType]
	multiplier := multipliers[entityType]
	if multiplier == 0 {
		multiplier = 1.0
	}

	return int(float64(baseXP) * multiplier)
}

// calculateAchievementProgress calcula el progreso de un achievement
func (s *gamificationService) calculateAchievementProgress(achievementType string, gamification *domain.UserGamification) int {
	switch achievementType {
	case domain.AchievementTypeAIPartner:
		return gamification.InsightsViewed
	case domain.AchievementTypeActionTaker:
		return gamification.ActionsCompleted
	case domain.AchievementTypeDataExplorer:
		return gamification.InsightsViewed / 5 // 5 insights = 1 progreso
	case domain.AchievementTypeQuickLearner:
		return gamification.InsightsViewed / 2 // 2 insights = 1 progreso
	case domain.AchievementTypeStreakKeeper:
		return gamification.CurrentStreak
	default:
		return 0
	}
}

// initializeBasicAchievements crea los achievements básicos para un usuario nuevo
func (s *gamificationService) initializeBasicAchievements(ctx context.Context, userID string) error {
	basicAchievements := []domain.Achievement{
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        domain.AchievementTypeAIPartner,
			Name:        "🤖 AI Partner",
			Description: "Utiliza 100 insights de IA",
			Points:      500,
			Progress:    0,
			Target:      100,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        domain.AchievementTypeActionTaker,
			Name:        "🎯 Action Taker",
			Description: "Completa 50 acciones sugeridas",
			Points:      300,
			Progress:    0,
			Target:      50,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        domain.AchievementTypeDataExplorer,
			Name:        "📊 Data Explorer",
			Description: "Explora insights 5 días consecutivos",
			Points:      200,
			Progress:    0,
			Target:      5,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        domain.AchievementTypeQuickLearner,
			Name:        "⚡ Quick Learner",
			Description: "Marca 10 insights como entendidos",
			Points:      100,
			Progress:    0,
			Target:      10,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, achievement := range basicAchievements {
		err := s.gamificationRepo.CreateAchievement(ctx, &achievement)
		if err != nil {
			return fmt.Errorf("error creating achievement %s: %w", achievement.Type, err)
		}
	}

	return nil
}
