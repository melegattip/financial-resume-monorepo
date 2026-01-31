package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// SavingsGoalRepository implements the savings goal repository interface (Repository pattern)
type SavingsGoalRepository struct {
	db *sql.DB
}

// NewSavingsGoalRepository creates a new instance of SavingsGoalRepository (Factory pattern)
func NewSavingsGoalRepository(db *sql.DB) ports.SavingsGoalRepository {
	return &SavingsGoalRepository{
		db: db,
	}
}

// Create creates a new savings goal in the database
func (r *SavingsGoalRepository) Create(ctx context.Context, goal *domain.SavingsGoal) error {
	query := `
		INSERT INTO savings_goals (
			id, user_id, name, description, target_amount, current_amount,
			category, priority, target_date, status, monthly_target,
			weekly_target, daily_target, progress, remaining_amount,
			days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			image_url, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22
		)`

	// auto_save_frequency debe ser NULL cuando is_auto_save = false para no violar el CHECK
	var autoSaveFreq sql.NullString
	if goal.IsAutoSave && goal.AutoSaveFrequency != "" {
		autoSaveFreq = sql.NullString{String: goal.AutoSaveFrequency, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		goal.ID, goal.UserID, goal.Name, goal.Description, goal.TargetAmount,
		goal.CurrentAmount, goal.Category, goal.Priority, goal.TargetDate,
		goal.Status, goal.MonthlyTarget, goal.WeeklyTarget, goal.DailyTarget,
		goal.Progress, goal.RemainingAmount, goal.DaysRemaining, goal.IsAutoSave,
		goal.AutoSaveAmount, autoSaveFreq, goal.ImageURL,
		goal.CreatedAt, goal.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create savings goal: %w", err)
	}

	return nil
}

// GetByID retrieves a savings goal by ID and user ID
func (r *SavingsGoalRepository) GetByID(ctx context.Context, userID, goalID string) (*domain.SavingsGoal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			   category, priority, target_date, status, monthly_target,
			   weekly_target, daily_target, progress, remaining_amount,
			   days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			   image_url, created_at, updated_at, achieved_at
		FROM savings_goals
		WHERE id = $1 AND user_id = $2`

	row := r.db.QueryRowContext(ctx, query, goalID, userID)

	goal := &domain.SavingsGoal{}
	var achievedAt sql.NullTime
	var description sql.NullString
	var autoSaveFrequency sql.NullString
	var imageURL sql.NullString

	err := row.Scan(
		&goal.ID, &goal.UserID, &goal.Name, &description, &goal.TargetAmount,
		&goal.CurrentAmount, &goal.Category, &goal.Priority, &goal.TargetDate,
		&goal.Status, &goal.MonthlyTarget, &goal.WeeklyTarget, &goal.DailyTarget,
		&goal.Progress, &goal.RemainingAmount, &goal.DaysRemaining, &goal.IsAutoSave,
		&goal.AutoSaveAmount, &autoSaveFrequency, &imageURL,
		&goal.CreatedAt, &goal.UpdatedAt, &achievedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get savings goal: %w", err)
	}

	if description.Valid {
		goal.Description = description.String
	}
	if autoSaveFrequency.Valid {
		goal.AutoSaveFrequency = autoSaveFrequency.String
	}
	if imageURL.Valid {
		goal.ImageURL = imageURL.String
	}
	if achievedAt.Valid {
		goal.AchievedAt = &achievedAt.Time
	}

	return goal, nil
}

// List retrieves all savings goals for a user
func (r *SavingsGoalRepository) List(ctx context.Context, userID string) ([]*domain.SavingsGoal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			   category, priority, target_date, status, monthly_target,
			   weekly_target, daily_target, progress, remaining_amount,
			   days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			   image_url, created_at, updated_at, achieved_at
		FROM savings_goals
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list savings goals: %w", err)
	}
	defer rows.Close()

	var goals []*domain.SavingsGoal
	for rows.Next() {
		goal := &domain.SavingsGoal{}
		var achievedAt sql.NullTime
		var description sql.NullString
		var autoSaveFrequency sql.NullString
		var imageURL sql.NullString

		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Name, &description, &goal.TargetAmount,
			&goal.CurrentAmount, &goal.Category, &goal.Priority, &goal.TargetDate,
			&goal.Status, &goal.MonthlyTarget, &goal.WeeklyTarget, &goal.DailyTarget,
			&goal.Progress, &goal.RemainingAmount, &goal.DaysRemaining, &goal.IsAutoSave,
			&goal.AutoSaveAmount, &autoSaveFrequency, &imageURL,
			&goal.CreatedAt, &goal.UpdatedAt, &achievedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan savings goal: %w", err)
		}

		if description.Valid {
			goal.Description = description.String
		}
		if autoSaveFrequency.Valid {
			goal.AutoSaveFrequency = autoSaveFrequency.String
		}
		if imageURL.Valid {
			goal.ImageURL = imageURL.String
		}
		if achievedAt.Valid {
			goal.AchievedAt = &achievedAt.Time
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate savings goals: %w", err)
	}

	return goals, nil
}

// ListByStatus retrieves savings goals by status for a user
func (r *SavingsGoalRepository) ListByStatus(ctx context.Context, userID string, status domain.SavingsGoalStatus) ([]*domain.SavingsGoal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			   category, priority, target_date, status, monthly_target,
			   weekly_target, daily_target, progress, remaining_amount,
			   days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			   image_url, created_at, updated_at, achieved_at
		FROM savings_goals
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list savings goals by status: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// ListByCategory retrieves savings goals by category for a user
func (r *SavingsGoalRepository) ListByCategory(ctx context.Context, userID string, category domain.SavingsGoalCategory) ([]*domain.SavingsGoal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			   category, priority, target_date, status, monthly_target,
			   weekly_target, daily_target, progress, remaining_amount,
			   days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			   image_url, created_at, updated_at, achieved_at
		FROM savings_goals
		WHERE user_id = $1 AND category = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list savings goals by category: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// Update updates an existing savings goal
func (r *SavingsGoalRepository) Update(ctx context.Context, goal *domain.SavingsGoal) error {
	query := `
		UPDATE savings_goals SET
			name = $3, description = $4, target_amount = $5, current_amount = $6,
			category = $7, priority = $8, target_date = $9, status = $10,
			monthly_target = $11, weekly_target = $12, daily_target = $13,
			progress = $14, remaining_amount = $15, days_remaining = $16,
			is_auto_save = $17, auto_save_amount = $18, auto_save_frequency = $19,
			image_url = $20, updated_at = $21, achieved_at = $22
		WHERE id = $1 AND user_id = $2`

	// auto_save_frequency NULL si no aplica
	var updAutoSaveFreq sql.NullString
	if goal.IsAutoSave && goal.AutoSaveFrequency != "" {
		updAutoSaveFreq = sql.NullString{String: goal.AutoSaveFrequency, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		goal.ID, goal.UserID, goal.Name, goal.Description, goal.TargetAmount,
		goal.CurrentAmount, goal.Category, goal.Priority, goal.TargetDate,
		goal.Status, goal.MonthlyTarget, goal.WeeklyTarget, goal.DailyTarget,
		goal.Progress, goal.RemainingAmount, goal.DaysRemaining, goal.IsAutoSave,
		goal.AutoSaveAmount, updAutoSaveFreq, goal.ImageURL,
		goal.UpdatedAt, goal.AchievedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update savings goal: %w", err)
	}

	return nil
}

// Delete removes a savings goal from the database
func (r *SavingsGoalRepository) Delete(ctx context.Context, userID, goalID string) error {
	// First, delete all related transactions
	_, err := r.db.ExecContext(ctx, "DELETE FROM savings_transactions WHERE goal_id = $1 AND user_id = $2", goalID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete savings transactions: %w", err)
	}

	// Then delete the goal
	_, err = r.db.ExecContext(ctx, "DELETE FROM savings_goals WHERE id = $1 AND user_id = $2", goalID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete savings goal: %w", err)
	}

	return nil
}

// CreateTransaction creates a new savings transaction
func (r *SavingsGoalRepository) CreateTransaction(ctx context.Context, transaction *domain.SavingsTransaction) error {
	query := `
		INSERT INTO savings_transactions (
			id, goal_id, user_id, amount, type, description, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`

	_, err := r.db.ExecContext(ctx, query,
		transaction.ID, transaction.GoalID, transaction.UserID,
		transaction.Amount, transaction.Type, transaction.Description,
		transaction.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create savings transaction: %w", err)
	}

	return nil
}

// GetTransactionsByGoal retrieves all transactions for a specific goal
func (r *SavingsGoalRepository) GetTransactionsByGoal(ctx context.Context, userID, goalID string) ([]*domain.SavingsTransaction, error) {
	query := `
		SELECT id, goal_id, user_id, amount, type, description, created_at
		FROM savings_transactions
		WHERE goal_id = $1 AND user_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, goalID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by goal: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetTransactionsByUser retrieves all transactions for a user
func (r *SavingsGoalRepository) GetTransactionsByUser(ctx context.Context, userID string) ([]*domain.SavingsTransaction, error) {
	query := `
		SELECT id, goal_id, user_id, amount, type, description, created_at
		FROM savings_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by user: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// Helper methods for scanning database results

// scanGoals scans database rows into savings goal objects
func (r *SavingsGoalRepository) scanGoals(rows *sql.Rows) ([]*domain.SavingsGoal, error) {
	var goals []*domain.SavingsGoal

	for rows.Next() {
		goal := &domain.SavingsGoal{}
		var achievedAt sql.NullTime
		var description sql.NullString
		var autoSaveFrequency sql.NullString
		var imageURL sql.NullString

		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Name, &description, &goal.TargetAmount,
			&goal.CurrentAmount, &goal.Category, &goal.Priority, &goal.TargetDate,
			&goal.Status, &goal.MonthlyTarget, &goal.WeeklyTarget, &goal.DailyTarget,
			&goal.Progress, &goal.RemainingAmount, &goal.DaysRemaining, &goal.IsAutoSave,
			&goal.AutoSaveAmount, &autoSaveFrequency, &imageURL,
			&goal.CreatedAt, &goal.UpdatedAt, &achievedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan savings goal: %w", err)
		}

		if description.Valid {
			goal.Description = description.String
		}
		if autoSaveFrequency.Valid {
			goal.AutoSaveFrequency = autoSaveFrequency.String
		}
		if imageURL.Valid {
			goal.ImageURL = imageURL.String
		}
		if achievedAt.Valid {
			goal.AchievedAt = &achievedAt.Time
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate savings goals: %w", err)
	}

	return goals, nil
}

// scanTransactions scans database rows into savings transaction objects
func (r *SavingsGoalRepository) scanTransactions(rows *sql.Rows) ([]*domain.SavingsTransaction, error) {
	var transactions []*domain.SavingsTransaction

	for rows.Next() {
		transaction := &domain.SavingsTransaction{}

		err := rows.Scan(
			&transaction.ID, &transaction.GoalID, &transaction.UserID,
			&transaction.Amount, &transaction.Type, &transaction.Description,
			&transaction.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan savings transaction: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate savings transactions: %w", err)
	}

	return transactions, nil
}

// GetActiveGoalsWithAutoSave retrieves all active goals with auto-save enabled
func (r *SavingsGoalRepository) GetActiveGoalsWithAutoSave(ctx context.Context) ([]*domain.SavingsGoal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			   category, priority, target_date, status, monthly_target,
			   weekly_target, daily_target, progress, remaining_amount,
			   days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			   image_url, created_at, updated_at, achieved_at
		FROM savings_goals
		WHERE status = 'active' AND is_auto_save = true
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active auto-save goals: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// GetOverdueGoals retrieves all overdue goals for a user
func (r *SavingsGoalRepository) GetOverdueGoals(ctx context.Context, userID string) ([]*domain.SavingsGoal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			   category, priority, target_date, status, monthly_target,
			   weekly_target, daily_target, progress, remaining_amount,
			   days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			   image_url, created_at, updated_at, achieved_at
		FROM savings_goals
		WHERE user_id = $1 
		  AND status IN ('active', 'paused')
		  AND target_date < NOW()
		  AND current_amount < target_amount
		ORDER BY target_date ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue goals: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// GetGoalsByPriority retrieves goals by priority for a user
func (r *SavingsGoalRepository) GetGoalsByPriority(ctx context.Context, userID string, priority domain.SavingsGoalPriority) ([]*domain.SavingsGoal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			   category, priority, target_date, status, monthly_target,
			   weekly_target, daily_target, progress, remaining_amount,
			   days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
			   image_url, created_at, updated_at, achieved_at
		FROM savings_goals
		WHERE user_id = $1 AND priority = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, priority)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals by priority: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// GetGoalsSummary retrieves summary statistics for a user's goals
func (r *SavingsGoalRepository) GetGoalsSummary(ctx context.Context, userID string) (*domain.SavingsGoalSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_goals,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_goals,
			COUNT(CASE WHEN status = 'achieved' THEN 1 END) as achieved_goals,
			COUNT(CASE WHEN status = 'paused' THEN 1 END) as paused_goals,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_goals,
			COALESCE(SUM(target_amount), 0) as total_target,
			COALESCE(SUM(current_amount), 0) as total_saved,
			COALESCE(SUM(remaining_amount), 0) as total_remaining,
			COALESCE(AVG(progress), 0) as average_progress,
			COUNT(CASE WHEN target_date < NOW() AND status IN ('active', 'paused') AND current_amount < target_amount THEN 1 END) as overdue_goals,
			COUNT(CASE WHEN status = 'active' AND progress >= 0.8 THEN 1 END) as on_track_goals
		FROM savings_goals
		WHERE user_id = $1`

	row := r.db.QueryRowContext(ctx, query, userID)

	summary := &domain.SavingsGoalSummary{}
	err := row.Scan(
		&summary.TotalGoals, &summary.ActiveGoals, &summary.AchievedGoals,
		&summary.PausedGoals, &summary.CancelledGoals, &summary.TotalTarget,
		&summary.TotalSaved, &summary.TotalRemaining, &summary.AverageProgress,
		&summary.OverdueGoals, &summary.OnTrackGoals,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get goals summary: %w", err)
	}

	return summary, nil
}
