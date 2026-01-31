package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/melegattip/financial-gamification-service/internal/core/domain"
	"github.com/melegattip/financial-gamification-service/internal/core/ports"
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
	log.Printf("🔍 [CreateGamification] Creating user gamification for user: %s", gamification.UserID)

	// Query simplificado sin updated_at para evitar problemas de columnas
	query := `
		INSERT INTO user_gamification (
			id, user_id, total_xp, current_level, insights_viewed, 
			actions_completed, achievements_count, current_streak, 
			last_activity, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
	)

	if err != nil {
		log.Printf("❌ [CreateGamification] Insert failed for user %s: %v", gamification.UserID, err)
		return fmt.Errorf("error creating user gamification: %w", err)
	}

	log.Printf("✅ [CreateGamification] Successfully created user gamification for user: %s", gamification.UserID)
	return nil
}

// GetByUserID obtiene la gamificación de un usuario por su ID
func (r *gamificationRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserGamification, error) {
	log.Printf("🔍 [GetUserGamification] Getting gamification for user: %s", userID)

	// Usar un enfoque más robusto: verificar qué columnas existen
	var gamification domain.UserGamification

	// Query que maneja cualquier estructura de tabla
	query := `
		SELECT id, user_id, total_xp, current_level, insights_viewed,
			   actions_completed, achievements_count, current_streak,
			   last_activity, created_at
		FROM user_gamification 
		WHERE user_id = $1
	`

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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("🔍 [GetUserGamification] No gamification found for user: %s", userID)
			return nil, ports.ErrGamificationNotFound
		}
		log.Printf("❌ [GetUserGamification] Query failed for user %s: %v", userID, err)
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Set updated_at to created_at for consistency
	gamification.UpdatedAt = gamification.CreatedAt
	log.Printf("✅ [GetUserGamification] Successfully retrieved gamification for user: %s", userID)
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

// GetAchievementsByUserID obtiene achievements de un usuario por su ID
func (r *gamificationRepository) GetAchievementsByUserID(ctx context.Context, userID string) ([]domain.Achievement, error) {
	log.Printf("🔍 [GetAchievements] Executing query for user: %s", userID)

	// Query con todas las columnas incluyendo updated_at
	query := `
		SELECT id, user_id, type, name, description, points, progress,
			   target, completed, unlocked_at, created_at, updated_at
		FROM achievements 
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("❌ [GetAchievements] Query failed for user %s: %v", userID, err)
		return nil, fmt.Errorf("error querying achievements: %w", err)
	}
	defer rows.Close()

	log.Printf("✅ [GetAchievements] Query successful, processing rows...")

	var achievements []domain.Achievement
	rowCount := 0
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
			log.Printf("❌ [GetAchievements] Scan failed at row %d: %v", rowCount, err)
			return nil, fmt.Errorf("error scanning achievement: %w", err)
		}
		achievements = append(achievements, achievement)
		rowCount++
	}

	if err = rows.Err(); err != nil {
		log.Printf("❌ [GetAchievements] Rows iteration error: %v", err)
		return nil, fmt.Errorf("error iterating achievements: %w", err)
	}

	log.Printf("✅ [GetAchievements] Successfully processed %d achievements for user %s", rowCount, userID)

	if achievements == nil {
		achievements = []domain.Achievement{}
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
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
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

// DeleteAction elimina una acción de usuario
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

// ====================================
// CHALLENGE OPERATIONS
// ====================================

// GetActiveChallenges obtiene todos los challenges activos de un tipo específico
func (r *gamificationRepository) GetActiveChallenges(ctx context.Context, challengeType string) ([]domain.Challenge, error) {
	query := `
		SELECT id, challenge_key, name, description, challenge_type, icon,
			   xp_reward, requirement_type, requirement_target, 
			   requirement_data, active, created_at, updated_at
		FROM challenges 
		WHERE challenge_type = $1 AND active = true
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, challengeType)
	if err != nil {
		return nil, fmt.Errorf("error querying active challenges: %w", err)
	}
	defer rows.Close()

	var challenges []domain.Challenge
	for rows.Next() {
		var challenge domain.Challenge
		var requirementDataJSON []byte
		err := rows.Scan(
			&challenge.ID,
			&challenge.ChallengeKey,
			&challenge.Name,
			&challenge.Description,
			&challenge.ChallengeType,
			&challenge.Icon,
			&challenge.XPReward,
			&challenge.RequirementType,
			&challenge.RequirementTarget,
			&requirementDataJSON,
			&challenge.Active,
			&challenge.CreatedAt,
			&challenge.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning challenge: %w", err)
		}
		// TODO: Parse requirementDataJSON to challenge.RequirementData if needed
		challenges = append(challenges, challenge)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating challenges: %w", err)
	}

	return challenges, nil
}

// GetChallengeByKey obtiene un challenge específico por su key
func (r *gamificationRepository) GetChallengeByKey(ctx context.Context, challengeKey string) (*domain.Challenge, error) {
	query := `
		SELECT id, challenge_key, name, description, challenge_type, icon,
			   xp_reward, requirement_type, requirement_target, 
			   requirement_data, active, created_at, updated_at
		FROM challenges 
		WHERE challenge_key = $1
	`

	var challenge domain.Challenge
	var requirementDataJSON []byte
	err := r.db.QueryRowContext(ctx, query, challengeKey).Scan(
		&challenge.ID,
		&challenge.ChallengeKey,
		&challenge.Name,
		&challenge.Description,
		&challenge.ChallengeType,
		&challenge.Icon,
		&challenge.XPReward,
		&challenge.RequirementType,
		&challenge.RequirementTarget,
		&requirementDataJSON,
		&challenge.Active,
		&challenge.CreatedAt,
		&challenge.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("challenge not found: %s", challengeKey)
		}
		return nil, fmt.Errorf("error getting challenge: %w", err)
	}

	// TODO: Parse requirementDataJSON to challenge.RequirementData if needed
	return &challenge, nil
}

// ====================================
// USER CHALLENGE OPERATIONS
// ====================================

// GetUserChallengesForDate obtiene los user challenges de un usuario para una fecha específica
func (r *gamificationRepository) GetUserChallengesForDate(ctx context.Context, userID string, challengeDate time.Time, challengeType string) ([]domain.UserChallenge, error) {
	// Join con challenges para filtrar por tipo
	query := `
		SELECT uc.id, uc.user_id, uc.challenge_id, uc.challenge_date,
			   uc.progress, uc.target, uc.completed, uc.completed_at,
			   uc.created_at, uc.updated_at
		FROM user_challenges uc
		JOIN challenges c ON uc.challenge_id = c.id
		WHERE uc.user_id = $1 AND uc.challenge_date = $2 AND c.challenge_type = $3
		ORDER BY uc.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, challengeDate.Format("2006-01-02"), challengeType)
	if err != nil {
		return nil, fmt.Errorf("error querying user challenges for date: %w", err)
	}
	defer rows.Close()

	var userChallenges []domain.UserChallenge
	for rows.Next() {
		var userChallenge domain.UserChallenge
		err := rows.Scan(
			&userChallenge.ID,
			&userChallenge.UserID,
			&userChallenge.ChallengeID,
			&userChallenge.ChallengeDate,
			&userChallenge.Progress,
			&userChallenge.Target,
			&userChallenge.Completed,
			&userChallenge.CompletedAt,
			&userChallenge.CreatedAt,
			&userChallenge.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user challenge: %w", err)
		}
		userChallenges = append(userChallenges, userChallenge)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user challenges: %w", err)
	}

	return userChallenges, nil
}

// CreateOrUpdateUserChallenge crea o actualiza un user challenge
func (r *gamificationRepository) CreateOrUpdateUserChallenge(ctx context.Context, userChallenge *domain.UserChallenge) error {
	// Primero intentamos hacer update
	updateQuery := `
		UPDATE user_challenges 
		SET progress = $4, target = $5, completed = $6, completed_at = $7, updated_at = $8
		WHERE user_id = $1 AND challenge_id = $2 AND challenge_date = $3
	`

	result, err := r.db.ExecContext(ctx, updateQuery,
		userChallenge.UserID,
		userChallenge.ChallengeID,
		userChallenge.ChallengeDate,
		userChallenge.Progress,
		userChallenge.Target,
		userChallenge.Completed,
		userChallenge.CompletedAt,
		userChallenge.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating user challenge: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	// Si no se actualizó ninguna fila, hacemos insert
	if rowsAffected == 0 {
		insertQuery := `
			INSERT INTO user_challenges (
				id, user_id, challenge_id, challenge_date,
				progress, target, completed, completed_at,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, err = r.db.ExecContext(ctx, insertQuery,
			userChallenge.ID,
			userChallenge.UserID,
			userChallenge.ChallengeID,
			userChallenge.ChallengeDate,
			userChallenge.Progress,
			userChallenge.Target,
			userChallenge.Completed,
			userChallenge.CompletedAt,
			userChallenge.CreatedAt,
			userChallenge.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("error creating user challenge: %w", err)
		}
	}

	return nil
}

// GetUserChallengeByID obtiene un user challenge por su ID
func (r *gamificationRepository) GetUserChallengeByID(ctx context.Context, userChallengeID string) (*domain.UserChallenge, error) {
	query := `
		SELECT id, user_id, challenge_id, challenge_date,
			   progress, target, completed, completed_at,
			   created_at, updated_at
		FROM user_challenges 
		WHERE id = $1
	`

	var userChallenge domain.UserChallenge
	var challengeDateStr string
	err := r.db.QueryRowContext(ctx, query, userChallengeID).Scan(
		&userChallenge.ID,
		&userChallenge.UserID,
		&userChallenge.ChallengeID,
		&challengeDateStr,
		&userChallenge.Progress,
		&userChallenge.Target,
		&userChallenge.Completed,
		&userChallenge.CompletedAt,
		&userChallenge.CreatedAt,
		&userChallenge.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user challenge not found: %s", userChallengeID)
		}
		return nil, fmt.Errorf("error getting user challenge: %w", err)
	}

	// Parse challenge date
	if userChallenge.ChallengeDate, err = time.Parse("2006-01-02", challengeDateStr); err != nil {
		return nil, fmt.Errorf("error parsing challenge date: %w", err)
	}

	return &userChallenge, nil
}

// ====================================
// CHALLENGE PROGRESS TRACKING OPERATIONS
// ====================================

// UpdateChallengeProgressTracking actualiza o crea el progreso de tracking de challenges
func (r *gamificationRepository) UpdateChallengeProgressTracking(ctx context.Context, tracking *domain.ChallengeProgressTracking) error {
	// Primero intentamos hacer update
	updateQuery := `
		UPDATE challenge_progress_tracking 
		SET count = count + $4
		WHERE user_id = $1 AND challenge_date = $2 AND action_type = $3
	`

	result, err := r.db.ExecContext(ctx, updateQuery,
		tracking.UserID,
		tracking.ChallengeDate.Format("2006-01-02"),
		tracking.ActionType,
		tracking.Count,
	)

	if err != nil {
		return fmt.Errorf("error updating challenge progress tracking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	// Si no se actualizó ninguna fila, hacemos insert
	if rowsAffected == 0 {
		insertQuery := `
			INSERT INTO challenge_progress_tracking (
				id, user_id, challenge_date, action_type, count,
				created_at
			) VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err = r.db.ExecContext(ctx, insertQuery,
			tracking.ID,
			tracking.UserID,
			tracking.ChallengeDate.Format("2006-01-02"),
			tracking.ActionType,
			tracking.Count,
			tracking.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf("error creating challenge progress tracking: %w", err)
		}
	}

	return nil
}

// GetChallengeProgressTracking obtiene el progreso de tracking para un usuario en una fecha
func (r *gamificationRepository) GetChallengeProgressTracking(ctx context.Context, userID string, challengeDate time.Time) ([]domain.ChallengeProgressTracking, error) {
	query := `
		SELECT id, user_id, challenge_date, action_type, count, created_at
		FROM challenge_progress_tracking 
		WHERE user_id = $1 AND challenge_date = $2
		ORDER BY action_type ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, challengeDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("error querying challenge progress tracking: %w", err)
	}
	defer rows.Close()

	var trackings []domain.ChallengeProgressTracking
	for rows.Next() {
		var tracking domain.ChallengeProgressTracking
		var challengeDateStr string
		err := rows.Scan(
			&tracking.ID,
			&tracking.UserID,
			&challengeDateStr,
			&tracking.ActionType,
			&tracking.Count,
			&tracking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning challenge progress tracking: %w", err)
		}

		// Parse challenge date
		if tracking.ChallengeDate, err = time.Parse("2006-01-02", challengeDateStr); err != nil {
			return nil, fmt.Errorf("error parsing challenge date: %w", err)
		}

		trackings = append(trackings, tracking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating challenge progress tracking: %w", err)
	}

	return trackings, nil
}
