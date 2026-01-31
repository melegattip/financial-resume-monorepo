package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// gamificationRepository implementa ports.GamificationRepository
type gamificationRepository struct {
	db *sql.DB
}

// NewGamificationRepository crea una nueva instancia del repository
func NewGamificationRepository(db *sql.DB) ports.GamificationRepository {
	return &gamificationRepository{
		db: db,
	}
}

// Create crea un nuevo registro de gamificación para un usuario
func (r *gamificationRepository) Create(ctx context.Context, gamification *domain.UserGamification) error {
	query := `
		INSERT INTO user_gamification (
			id, user_id, total_xp, current_level, insights_viewed, 
			actions_completed, achievements_count, current_streak, 
			last_activity, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		gamification.ID,
		gamification.UserID,
		gamification.TotalXP,
		gamification.CurrentLevel,
		gamification.InsightsViewed,
		gamification.ActionsCompleted,
		gamification.AchievementsCount,
		gamification.CurrentStreak,
		gamification.LastActivity,
		gamification.CreatedAt,
		gamification.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating user gamification: %w", err)
	}

	return nil
}

// GetByUserID obtiene la gamificación de un usuario por su ID
func (r *gamificationRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserGamification, error) {
	query := `
		SELECT id, user_id, total_xp, current_level, insights_viewed,
			   actions_completed, achievements_count, current_streak,
			   last_activity, created_at, updated_at
		FROM user_gamification 
		WHERE user_id = $1
	`

	var gamification domain.UserGamification
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&gamification.ID,
		&gamification.UserID,
		&gamification.TotalXP,
		&gamification.CurrentLevel,
		&gamification.InsightsViewed,
		&gamification.ActionsCompleted,
		&gamification.AchievementsCount,
		&gamification.CurrentStreak,
		&gamification.LastActivity,
		&gamification.CreatedAt,
		&gamification.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrGamificationNotFound
		}
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	return &gamification, nil
}

// Update actualiza un registro de gamificación
func (r *gamificationRepository) Update(ctx context.Context, gamification *domain.UserGamification) error {
	query := `
		UPDATE user_gamification 
		SET total_xp = $2, current_level = $3, insights_viewed = $4,
			actions_completed = $5, achievements_count = $6, current_streak = $7,
			last_activity = $8, updated_at = $9
		WHERE user_id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		gamification.UserID,
		gamification.TotalXP,
		gamification.CurrentLevel,
		gamification.InsightsViewed,
		gamification.ActionsCompleted,
		gamification.AchievementsCount,
		gamification.CurrentStreak,
		gamification.LastActivity,
		gamification.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating user gamification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrGamificationNotFound
	}

	return nil
}

// Delete elimina un registro de gamificación
func (r *gamificationRepository) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM user_gamification WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user gamification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrGamificationNotFound
	}

	return nil
}

// CreateAchievement crea un nuevo achievement
func (r *gamificationRepository) CreateAchievement(ctx context.Context, achievement *domain.Achievement) error {
	query := `
		INSERT INTO achievements (
			id, user_id, type, name, description, points, progress, 
			target, completed, unlocked_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		achievement.ID,
		achievement.UserID,
		achievement.Type,
		achievement.Name,
		achievement.Description,
		achievement.Points,
		achievement.Progress,
		achievement.Target,
		achievement.Completed,
		achievement.UnlockedAt,
		achievement.CreatedAt,
		achievement.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating achievement: %w", err)
	}

	return nil
}

// GetAchievementsByUserID obtiene todos los achievements de un usuario
func (r *gamificationRepository) GetAchievementsByUserID(ctx context.Context, userID string) ([]domain.Achievement, error) {
	query := `
		SELECT id, user_id, type, name, description, points, progress,
			   target, completed, unlocked_at, created_at, updated_at
		FROM achievements 
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying achievements: %w", err)
	}
	defer rows.Close()

	var achievements []domain.Achievement
	for rows.Next() {
		var achievement domain.Achievement
		err := rows.Scan(
			&achievement.ID,
			&achievement.UserID,
			&achievement.Type,
			&achievement.Name,
			&achievement.Description,
			&achievement.Points,
			&achievement.Progress,
			&achievement.Target,
			&achievement.Completed,
			&achievement.UnlockedAt,
			&achievement.CreatedAt,
			&achievement.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning achievement: %w", err)
		}
		achievements = append(achievements, achievement)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating achievements: %w", err)
	}

	return achievements, nil
}

// GetAchievementByID obtiene un achievement por su ID
func (r *gamificationRepository) GetAchievementByID(ctx context.Context, achievementID string) (*domain.Achievement, error) {
	query := `
		SELECT id, user_id, type, name, description, points, progress,
			   target, completed, unlocked_at, created_at, updated_at
		FROM achievements 
		WHERE id = $1
	`

	var achievement domain.Achievement
	err := r.db.QueryRowContext(ctx, query, achievementID).Scan(
		&achievement.ID,
		&achievement.UserID,
		&achievement.Type,
		&achievement.Name,
		&achievement.Description,
		&achievement.Points,
		&achievement.Progress,
		&achievement.Target,
		&achievement.Completed,
		&achievement.UnlockedAt,
		&achievement.CreatedAt,
		&achievement.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrAchievementNotFound
		}
		return nil, fmt.Errorf("error getting achievement: %w", err)
	}

	return &achievement, nil
}

// UpdateAchievement actualiza un achievement
func (r *gamificationRepository) UpdateAchievement(ctx context.Context, achievement *domain.Achievement) error {
	query := `
		UPDATE achievements 
		SET progress = $2, completed = $3, unlocked_at = $4, updated_at = $5
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		achievement.ID,
		achievement.Progress,
		achievement.Completed,
		achievement.UnlockedAt,
		achievement.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating achievement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrAchievementNotFound
	}

	return nil
}

// DeleteAchievement elimina un achievement
func (r *gamificationRepository) DeleteAchievement(ctx context.Context, achievementID string) error {
	query := `DELETE FROM achievements WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, achievementID)
	if err != nil {
		return fmt.Errorf("error deleting achievement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrAchievementNotFound
	}

	return nil
}

// CreateAction crea una nueva acción de usuario
func (r *gamificationRepository) CreateAction(ctx context.Context, action *domain.UserAction) error {
	query := `
		INSERT INTO user_actions (
			id, user_id, action_type, entity_type, entity_id, 
			xp_earned, description, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		action.ID,
		action.UserID,
		action.ActionType,
		action.EntityType,
		action.EntityID,
		action.XPEarned,
		action.Description,
		action.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating user action: %w", err)
	}

	return nil
}

// GetActionsByUserID obtiene todas las acciones de un usuario
func (r *gamificationRepository) GetActionsByUserID(ctx context.Context, userID string) ([]domain.UserAction, error) {
	query := `
		SELECT id, user_id, action_type, entity_type, entity_id,
			   xp_earned, description, created_at
		FROM user_actions 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user actions: %w", err)
	}
	defer rows.Close()

	var actions []domain.UserAction
	for rows.Next() {
		var action domain.UserAction
		err := rows.Scan(
			&action.ID,
			&action.UserID,
			&action.ActionType,
			&action.EntityType,
			&action.EntityID,
			&action.XPEarned,
			&action.Description,
			&action.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user action: %w", err)
		}
		actions = append(actions, action)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user actions: %w", err)
	}

	return actions, nil
}

// GetActionsByUserIDAndPeriod obtiene acciones de un usuario en un período
func (r *gamificationRepository) GetActionsByUserIDAndPeriod(ctx context.Context, userID string, startDate, endDate string) ([]domain.UserAction, error) {
	query := `
		SELECT id, user_id, action_type, entity_type, entity_id,
			   xp_earned, description, created_at
		FROM user_actions 
		WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying user actions by period: %w", err)
	}
	defer rows.Close()

	var actions []domain.UserAction
	for rows.Next() {
		var action domain.UserAction
		err := rows.Scan(
			&action.ID,
			&action.UserID,
			&action.ActionType,
			&action.EntityType,
			&action.EntityID,
			&action.XPEarned,
			&action.Description,
			&action.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user action: %w", err)
		}
		actions = append(actions, action)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user actions: %w", err)
	}

	return actions, nil
}

// DeleteAction elimina una acción
func (r *gamificationRepository) DeleteAction(ctx context.Context, actionID string) error {
	query := `DELETE FROM user_actions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, actionID)
	if err != nil {
		return fmt.Errorf("error deleting user action: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrActionNotFound
	}

	return nil
}
