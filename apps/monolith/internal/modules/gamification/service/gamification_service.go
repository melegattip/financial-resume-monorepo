package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/repository"
)

// GamificationStats is the summary response returned by GetGamificationStats.
type GamificationStats struct {
	UserID                string    `json:"user_id"`
	TotalXP               int       `json:"total_xp"`
	CurrentLevel          int       `json:"current_level"`
	LevelName             string    `json:"level_name"`
	XPToNextLevel         int       `json:"xp_to_next_level"`
	ProgressPercent       int       `json:"progress_percent"`
	TotalAchievements     int       `json:"total_achievements"`
	CompletedAchievements int       `json:"completed_achievements"`
	CurrentStreak         int       `json:"current_streak"`
	LastActivity          time.Time `json:"last_activity"`
}

// RecordActionResult is returned by RecordAction so that callers can compute
// XP gained and whether the user levelled up during the action.
type RecordActionResult struct {
	XPEarned     int
	TotalXP      int
	CurrentLevel int
	LevelUp      bool
}

// GamificationService contains all gamification business logic.
type GamificationService struct {
	repo   *repository.GamificationRepo
	logger zerolog.Logger
}

// NewGamificationService creates a new GamificationService.
func NewGamificationService(repo *repository.GamificationRepo, logger zerolog.Logger) *GamificationService {
	return &GamificationService{repo: repo, logger: logger}
}

// ---------------------------------------------------------------------------
// Public methods
// ---------------------------------------------------------------------------

// InitializeUserGamification creates the gamification profile and default
// achievements for a user. It is idempotent; a second call for the same user
// is a safe no-op.
func (s *GamificationService) InitializeUserGamification(ctx context.Context, userID string) error {
	existing, err := s.repo.FindUserGamificationByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if existing != nil {
		// Already initialized.
		return nil
	}

	g := domain.NewUserGamification(userID)
	if err := s.repo.CreateUserGamification(ctx, g); err != nil {
		return err
	}

	achievements := domain.DefaultAchievements(userID)
	if err := s.repo.CreateAchievements(ctx, achievements); err != nil {
		return err
	}

	s.logger.Info().Str("user_id", userID).Msg("gamification initialized for user")
	return nil
}

// RecordAction records a user action, awards XP, updates the streak, and
// refreshes achievement progress. It returns a RecordActionResult so the
// caller can surface level-up information.
func (s *GamificationService) RecordAction(ctx context.Context, userID, actionType, entityID string) (*RecordActionResult, error) {
	// Ensure the user has a gamification profile.
	g, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Anti-farming: award view_dashboard XP only once per calendar day.
	if actionType == domain.ActionViewDashboard {
		todaysActions, err := s.repo.FindActionsByUserIDAndDay(ctx, userID, time.Now().UTC())
		if err != nil {
			return nil, err
		}
		for _, a := range todaysActions {
			if a.ActionType == domain.ActionViewDashboard {
				// Already recorded today — skip silently.
				return &RecordActionResult{
					XPEarned:     0,
					TotalXP:      g.TotalXP,
					CurrentLevel: g.CurrentLevel,
					LevelUp:      false,
				}, nil
			}
		}
	}

	xp := domain.XPForAction(actionType)
	previousLevel := g.CurrentLevel

	// Update streak for daily login actions.
	if actionType == domain.ActionDailyLogin {
		s.updateStreak(g)
	}

	// Persist the action record.
	action := &domain.UserAction{
		ID:         uuid.New().String(),
		UserID:     userID,
		ActionType: actionType,
		EntityID:   entityID,
		XPEarned:   xp,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.CreateUserAction(ctx, action); err != nil {
		return nil, err
	}

	// Update XP, level, and activity timestamp on the gamification aggregate.
	g.TotalXP += xp
	g.ActionsCompleted++
	g.CurrentLevel = g.CalculateLevel()
	g.LastActivity = time.Now().UTC()
	g.UpdatedAt = time.Now().UTC()

	// Fetch all actions to compute achievement progress.
	allActions, err := s.repo.FindActionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	newlyCompleted, err := s.updateAchievementsProgress(ctx, g, allActions)
	if err != nil {
		s.logger.Warn().Err(err).Str("user_id", userID).Msg("failed to update achievements")
	}

	// Award bonus XP for newly completed achievements.
	for _, ach := range newlyCompleted {
		g.TotalXP += ach.Points
		g.AchievementsCount++
	}

	// Recompute level in case achievement XP pushed the user up.
	g.CurrentLevel = g.CalculateLevel()

	if err := s.repo.UpdateUserGamification(ctx, g); err != nil {
		return nil, err
	}

	return &RecordActionResult{
		XPEarned:     xp,
		TotalXP:      g.TotalXP,
		CurrentLevel: g.CurrentLevel,
		LevelUp:      g.CurrentLevel > previousLevel,
	}, nil
}

// GetUserGamification returns the gamification profile for a user, initialising
// it on first access.
func (s *GamificationService) GetUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error) {
	g, err := s.repo.FindUserGamificationByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if g == nil {
		if err := s.InitializeUserGamification(ctx, userID); err != nil {
			return nil, err
		}
		g, err = s.repo.FindUserGamificationByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}

// GetGamificationStats returns a summary view of the user's gamification state.
func (s *GamificationService) GetGamificationStats(ctx context.Context, userID string) (*GamificationStats, error) {
	g, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, err
	}

	achievements, err := s.repo.FindAchievementsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	completed := 0
	for _, a := range achievements {
		if a.Completed {
			completed++
		}
	}

	return &GamificationStats{
		UserID:                userID,
		TotalXP:               g.TotalXP,
		CurrentLevel:          g.CurrentLevel,
		LevelName:             g.GetLevelName(),
		XPToNextLevel:         g.XPToNextLevel(),
		ProgressPercent:       g.ProgressToNextLevel(),
		TotalAchievements:     len(achievements),
		CompletedAchievements: completed,
		CurrentStreak:         g.CurrentStreak,
		LastActivity:          g.LastActivity,
	}, nil
}

// GetAchievements returns the list of achievements for a user.
func (s *GamificationService) GetAchievements(ctx context.Context, userID string) ([]domain.Achievement, error) {
	return s.repo.FindAchievementsByUserID(ctx, userID)
}

// ---------------------------------------------------------------------------
// Private helpers
// ---------------------------------------------------------------------------

// updateStreak adjusts the user's login streak based on the last activity date.
// A 1-day grace period is applied: missing exactly one day does not reset the
// streak, it simply does not increment it for that gap day.
func (s *GamificationService) updateStreak(g *domain.UserGamification) {
	now := time.Now().UTC()
	lastDay := time.Date(g.LastActivity.Year(), g.LastActivity.Month(), g.LastActivity.Day(), 0, 0, 0, 0, time.UTC)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	daysDiff := int(today.Sub(lastDay).Hours() / 24)

	switch {
	case daysDiff == 0:
		// Same day — do not change the streak.
	case daysDiff == 1:
		// Consecutive day — extend the streak.
		g.CurrentStreak++
	case daysDiff == 2:
		// 1-day grace period: user missed one day but streak is preserved.
		// We still increment because today counts as the next active day.
		g.CurrentStreak++
	default:
		// Missed more than one day — reset to 1.
		g.CurrentStreak = 1
	}
}

// updateAchievementsProgress recalculates progress for every achievement and
// persists any changes. It returns the achievements that were newly completed
// during this call.
func (s *GamificationService) updateAchievementsProgress(
	ctx context.Context,
	g *domain.UserGamification,
	actions []domain.UserAction,
) ([]domain.Achievement, error) {
	achievements, err := s.repo.FindAchievementsByUserID(ctx, g.UserID)
	if err != nil {
		return nil, err
	}

	// Pre-compute action counts we need for multiple achievements.
	txCount := 0        // create_expense + create_income
	categoryCount := 0  // create_category
	assignCount := 0    // assign_category

	for _, a := range actions {
		switch a.ActionType {
		case domain.ActionCreateExpense, domain.ActionCreateIncome:
			txCount++
		case domain.ActionCreateCategory:
			categoryCount++
		case domain.ActionAssignCategory:
			assignCount++
		}
	}

	daysUsed := int(time.Since(g.CreatedAt).Hours() / 24)

	var newlyCompleted []domain.Achievement

	for i := range achievements {
		ach := &achievements[i]
		if ach.Completed {
			// Already done — skip.
			continue
		}

		wasCompleted := ach.Completed
		var newProgress int

		switch ach.Type {
		case "transaction_starter", "transaction_apprentice", "transaction_master":
			newProgress = txCount
		case "category_creator":
			newProgress = categoryCount
		case "organization_expert":
			newProgress = assignCount
		case "weekly_warrior", "monthly_legend":
			newProgress = g.CurrentStreak
		case "data_explorer":
			newProgress = daysUsed
		default:
			continue
		}

		if newProgress == ach.Progress {
			// No change.
			continue
		}

		ach.UpdateProgress(newProgress)
		if err := s.repo.UpdateAchievement(ctx, ach); err != nil {
			s.logger.Warn().Err(err).Str("achievement_id", ach.ID).Msg("failed to update achievement")
			continue
		}

		if !wasCompleted && ach.Completed {
			newlyCompleted = append(newlyCompleted, *ach)
			s.logger.Info().
				Str("user_id", g.UserID).
				Str("achievement", ach.Type).
				Int("points", ach.Points).
				Msg("achievement unlocked")
		}
	}

	return newlyCompleted, nil
}
