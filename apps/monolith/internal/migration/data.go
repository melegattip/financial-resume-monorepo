package migration

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// sourceUser represents a user row from the old users_db.
// The old schema uses SERIAL (integer) IDs and a single 'name' field.
type sourceUser struct {
	ID           int64  `gorm:"column:id"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	Name         string `gorm:"column:name"`
	Phone        string `gorm:"column:phone"`
	Avatar       string `gorm:"column:avatar"`
	IsActive     bool   `gorm:"column:is_active"`
	IsVerified   bool   `gorm:"column:is_verified"`
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// targetUser represents a user row in the target financial_resume.users table.
type targetUser struct {
	ID           string `gorm:"column:id;type:varchar(255);primaryKey"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	FirstName    string `gorm:"column:first_name"`
	LastName     string `gorm:"column:last_name"`
	Phone        string `gorm:"column:phone"`
	Avatar       string `gorm:"column:avatar"`
	IsActive     bool   `gorm:"column:is_active"`
	IsVerified   bool   `gorm:"column:is_verified"`
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (targetUser) TableName() string { return "users" }

// sourcePreference represents a user_preferences row from users_db.
type sourcePreference struct {
	ID                   int64  `gorm:"column:id"`
	UserID               int64  `gorm:"column:user_id"`
	Currency             string `gorm:"column:currency"`
	Language             string `gorm:"column:language"`
	Timezone             string `gorm:"column:timezone"`
	GamificationEnabled  bool   `gorm:"column:gamification_enabled"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// targetPreference represents a user_preferences row in the target database.
type targetPreference struct {
	UserID              string `gorm:"column:user_id;type:varchar(255);primaryKey"`
	Currency            string `gorm:"column:currency"`
	Language            string `gorm:"column:language"`
	Timezone            string `gorm:"column:timezone"`
	GamificationEnabled bool   `gorm:"column:gamification_enabled"`
	Theme               string `gorm:"column:theme"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (targetPreference) TableName() string { return "user_preferences" }

// CopyUsers reads all users from the source users_db, converts integer IDs to
// strings, splits the 'name' field into first_name/last_name, and inserts into
// the target database with ON CONFLICT (id) DO NOTHING.
func CopyUsers(sourceDB, targetDB *gorm.DB, log zerolog.Logger, dryRun bool) (copied, skipped int64, err error) {
	var users []sourceUser
	if err := sourceDB.Table("users").Find(&users).Error; err != nil {
		return 0, 0, fmt.Errorf("read users from source: %w", err)
	}

	log.Info().Int("source_count", len(users)).Msg("read users from users_db")

	if dryRun {
		log.Info().Int("count", len(users)).Msg("[DRY RUN] would copy users")
		return int64(len(users)), 0, nil
	}

	for _, u := range users {
		firstName, lastName := splitName(u.Name)
		target := targetUser{
			ID:           fmt.Sprintf("%d", u.ID),
			Email:        u.Email,
			PasswordHash: u.PasswordHash,
			FirstName:    firstName,
			LastName:     lastName,
			Phone:        u.Phone,
			Avatar:       u.Avatar,
			IsActive:     u.IsActive,
			IsVerified:   u.IsVerified,
			LastLogin:    u.LastLogin,
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
		}

		// ON CONFLICT (id) DO NOTHING — skip if user already exists.
		result := targetDB.Exec(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, phone, avatar, is_active, is_verified, last_login, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING
		`, target.ID, target.Email, target.PasswordHash, target.FirstName, target.LastName,
			target.Phone, target.Avatar, target.IsActive, target.IsVerified,
			target.LastLogin, target.CreatedAt, target.UpdatedAt)

		if result.Error != nil {
			log.Error().Err(result.Error).Str("user_id", target.ID).Msg("failed to insert user")
			return copied, skipped, fmt.Errorf("insert user %s: %w", target.ID, result.Error)
		}

		if result.RowsAffected == 0 {
			skipped++
			log.Debug().Str("user_id", target.ID).Msg("user skipped (already exists)")
		} else {
			copied++
		}
	}

	log.Info().Int64("copied", copied).Int64("skipped", skipped).Msg("user migration complete")
	return copied, skipped, nil
}

// CopyUserPreferences reads user_preferences from users_db, converts user_id
// to string, and inserts into target with ON CONFLICT (user_id) DO NOTHING.
func CopyUserPreferences(sourceDB, targetDB *gorm.DB, log zerolog.Logger, dryRun bool) (copied, skipped int64, err error) {
	var prefs []sourcePreference
	if err := sourceDB.Table("user_preferences").Find(&prefs).Error; err != nil {
		return 0, 0, fmt.Errorf("read preferences from source: %w", err)
	}

	log.Info().Int("source_count", len(prefs)).Msg("read user_preferences from users_db")

	if dryRun {
		log.Info().Int("count", len(prefs)).Msg("[DRY RUN] would copy user_preferences")
		return int64(len(prefs)), 0, nil
	}

	for _, p := range prefs {
		userID := fmt.Sprintf("%d", p.UserID)

		result := targetDB.Exec(`
			INSERT INTO user_preferences (user_id, currency, language, timezone, gamification_enabled, theme, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 'light', ?, ?)
			ON CONFLICT (user_id) DO NOTHING
		`, userID, p.Currency, p.Language, p.Timezone, p.GamificationEnabled,
			p.CreatedAt, p.UpdatedAt)

		if result.Error != nil {
			log.Error().Err(result.Error).Str("user_id", userID).Msg("failed to insert preference")
			return copied, skipped, fmt.Errorf("insert preference user_id=%s: %w", userID, result.Error)
		}

		if result.RowsAffected == 0 {
			skipped++
		} else {
			copied++
		}
	}

	log.Info().Int64("copied", copied).Int64("skipped", skipped).Msg("user_preferences migration complete")
	return copied, skipped, nil
}

// splitName splits a full name string into (firstName, lastName).
// If there's no space, firstName gets the full name and lastName is empty.
func splitName(name string) (string, string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ""
	}
	parts := strings.SplitN(name, " ", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.TrimSpace(parts[1])
}

// --- US2: Gamification data copy + deduplication ---

// CopyGamificationCore copies user_gamification, achievements, and user_actions
// from gamification-db to the target database with ON CONFLICT DO NOTHING.
func CopyGamificationCore(gamDB, targetDB *gorm.DB, log zerolog.Logger, dryRun bool) (map[string]int64, map[string]int64, error) {
	copied := make(map[string]int64)
	skipped := make(map[string]int64)

	// user_gamification
	c, s, err := copyTableGeneric[SrcUserGamification](gamDB, targetDB, log, dryRun, "user_gamification",
		`INSERT INTO user_gamification (id, user_id, financial_health_score, engagement_component, health_component,
		 current_level, insights_viewed, actions_completed, achievements_count, current_streak,
		 last_activity, last_score_calculation, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcUserGamification) []interface{} {
			return []interface{}{r.ID, r.UserID, r.FinancialHealthScore, r.EngagementComponent, r.HealthComponent,
				r.CurrentLevel, r.InsightsViewed, r.ActionsCompleted, r.AchievementsCount, r.CurrentStreak,
				r.LastActivity, r.LastScoreCalculation, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy user_gamification: %w", err)
	}
	copied["user_gamification"] = c
	skipped["user_gamification"] = s

	// achievements
	c, s, err = copyTableGeneric[SrcAchievement](gamDB, targetDB, log, dryRun, "achievements",
		`INSERT INTO achievements (id, user_id, type, name, description, points, progress, target, completed, unlocked_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcAchievement) []interface{} {
			return []interface{}{r.ID, r.UserID, r.Type, r.Name, r.Description, r.Points, r.Progress, r.Target, r.Completed, r.UnlockedAt, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy achievements: %w", err)
	}
	copied["achievements"] = c
	skipped["achievements"] = s

	// user_actions: skipped — records already exist in target DB
	log.Info().Msg("skipping user_actions (already migrated)")

	return copied, skipped, nil
}

// CopyGamificationTables copies challenges, user_challenges, and
// challenge_progress_tracking from gamification-db to the target database.
func CopyGamificationTables(gamDB, targetDB *gorm.DB, log zerolog.Logger, dryRun bool) (map[string]int64, map[string]int64, error) {
	copied := make(map[string]int64)
	skipped := make(map[string]int64)

	// challenges
	c, s, err := copyTableGeneric[SrcChallenge](gamDB, targetDB, log, dryRun, "challenges",
		`INSERT INTO challenges (id, challenge_key, name, description, challenge_type, icon, xp_reward,
		 requirement_type, requirement_target, requirement_data, active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcChallenge) []interface{} {
			return []interface{}{r.ID, r.ChallengeKey, r.Name, r.Description, r.ChallengeType, r.Icon, r.XPReward,
				r.RequirementType, r.RequirementTarget, r.RequirementData, r.Active, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy challenges: %w", err)
	}
	copied["challenges"] = c
	skipped["challenges"] = s

	// user_challenges
	c, s, err = copyTableGeneric[SrcUserChallenge](gamDB, targetDB, log, dryRun, "user_challenges",
		`INSERT INTO user_challenges (id, user_id, challenge_id, challenge_date, progress, target, completed, completed_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcUserChallenge) []interface{} {
			return []interface{}{r.ID, r.UserID, r.ChallengeID, r.ChallengeDate, r.Progress, r.Target, r.Completed, r.CompletedAt, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy user_challenges: %w", err)
	}
	copied["user_challenges"] = c
	skipped["user_challenges"] = s

	// challenge_progress_tracking
	c, s, err = copyTableGeneric[SrcChallengeProgressTracking](gamDB, targetDB, log, dryRun, "challenge_progress_tracking",
		`INSERT INTO challenge_progress_tracking (id, user_id, challenge_date, action_type, entity_type, count, unique_entities, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?::jsonb, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcChallengeProgressTracking) []interface{} {
			// Ensure unique_entities is valid JSON (empty array if empty/null)
			uniqueEntities := r.UniqueEntities
			if uniqueEntities == "" || uniqueEntities == "null" {
				uniqueEntities = "[]"
			}
			return []interface{}{r.ID, r.UserID, r.ChallengeDate, r.ActionType, r.EntityType, r.Count, uniqueEntities, r.CreatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy challenge_progress_tracking: %w", err)
	}
	copied["challenge_progress_tracking"] = c
	skipped["challenge_progress_tracking"] = s

	return copied, skipped, nil
}

// copyTableGeneric reads all rows of type T from sourceDB and inserts into
// targetDB using the given INSERT ... ON CONFLICT query.
func copyTableGeneric[T any](sourceDB, targetDB *gorm.DB, log zerolog.Logger, dryRun bool, tableName string, insertSQL string, argsFn func(T) []interface{}) (int64, int64, error) {
	var rows []T
	if err := sourceDB.Find(&rows).Error; err != nil {
		return 0, 0, fmt.Errorf("read %s from source: %w", tableName, err)
	}

	log.Info().Str("table", tableName).Int("source_count", len(rows)).Msg("read from gamification-db")

	if dryRun {
		log.Info().Str("table", tableName).Int("count", len(rows)).Msg("[DRY RUN] would copy")
		return int64(len(rows)), 0, nil
	}

	var copied, skipped int64
	for _, row := range rows {
		args := argsFn(row)
		result := targetDB.Exec(insertSQL, args...)
		if result.Error != nil {
			return copied, skipped, fmt.Errorf("insert into %s: %w", tableName, result.Error)
		}
		if result.RowsAffected == 0 {
			skipped++
		} else {
			copied++
		}
	}

	log.Info().Str("table", tableName).Int64("copied", copied).Int64("skipped", skipped).Msg("copy complete")
	return copied, skipped, nil
}

// DeduplicateUserGamification removes duplicate user_gamification rows,
// keeping the one with the highest financial_health_score per user_id.
func DeduplicateUserGamification(db *gorm.DB, log zerolog.Logger, dryRun bool) (DedupResult, error) {
	result := DedupResult{Table: "user_gamification"}

	// Find users with duplicate rows.
	type dupUser struct {
		UserID string
		Cnt    int64
	}
	var dups []dupUser
	if err := db.Raw(`
		SELECT user_id, COUNT(*) as cnt FROM user_gamification
		GROUP BY user_id HAVING COUNT(*) > 1
	`).Scan(&dups).Error; err != nil {
		return result, fmt.Errorf("find duplicates: %w", err)
	}

	for _, d := range dups {
		result.DuplicatesSeen += d.Cnt - 1 // extra rows beyond the keeper
	}

	if result.DuplicatesSeen == 0 {
		log.Info().Msg("no user_gamification duplicates found")
		return result, nil
	}

	log.Info().Int64("duplicates", result.DuplicatesSeen).Int("users_affected", len(dups)).Msg("user_gamification duplicates found")

	if dryRun {
		log.Info().Int64("would_remove", result.DuplicatesSeen).Msg("[DRY RUN] user_gamification dedup")
		return result, nil
	}

	// Delete all but the row with the highest financial_health_score per user_id.
	res := db.Exec(`
		DELETE FROM user_gamification
		WHERE id NOT IN (
			SELECT DISTINCT ON (user_id) id
			FROM user_gamification
			ORDER BY user_id, financial_health_score DESC, updated_at DESC
		)
	`)
	if res.Error != nil {
		return result, fmt.Errorf("dedup delete: %w", res.Error)
	}
	result.Removed = res.RowsAffected
	result.Kept = result.DuplicatesSeen + int64(len(dups)) - result.Removed

	log.Info().Int64("removed", result.Removed).Msg("user_gamification dedup complete")
	return result, nil
}

// DeduplicateAchievements removes duplicate achievements per (user_id, type),
// keeping the one with the most recent updated_at.
func DeduplicateAchievements(db *gorm.DB, log zerolog.Logger, dryRun bool) (DedupResult, error) {
	result := DedupResult{Table: "achievements"}

	type dupPair struct {
		UserID string
		Type   string
		Cnt    int64
	}
	var dups []dupPair
	if err := db.Raw(`
		SELECT user_id, type, COUNT(*) as cnt FROM achievements
		GROUP BY user_id, type HAVING COUNT(*) > 1
	`).Scan(&dups).Error; err != nil {
		return result, fmt.Errorf("find duplicates: %w", err)
	}

	for _, d := range dups {
		result.DuplicatesSeen += d.Cnt - 1
	}

	if result.DuplicatesSeen == 0 {
		log.Info().Msg("no achievements duplicates found")
		return result, nil
	}

	log.Info().Int64("duplicates", result.DuplicatesSeen).Msg("achievements duplicates found")

	if dryRun {
		log.Info().Int64("would_remove", result.DuplicatesSeen).Msg("[DRY RUN] achievements dedup")
		return result, nil
	}

	res := db.Exec(`
		DELETE FROM achievements
		WHERE id NOT IN (
			SELECT DISTINCT ON (user_id, type) id
			FROM achievements
			ORDER BY user_id, type, updated_at DESC
		)
	`)
	if res.Error != nil {
		return result, fmt.Errorf("dedup delete: %w", res.Error)
	}
	result.Removed = res.RowsAffected
	result.Kept = result.DuplicatesSeen + int64(len(dups)) - result.Removed

	log.Info().Int64("removed", result.Removed).Msg("achievements dedup complete")
	return result, nil
}

// DeduplicateUserActions removes duplicate user_actions by id (PK conflicts).
// Since gamification-db data was inserted first with ON CONFLICT DO NOTHING,
// this should typically be a no-op.
func DeduplicateUserActions(db *gorm.DB, log zerolog.Logger, dryRun bool) (DedupResult, error) {
	result := DedupResult{Table: "user_actions"}

	// user_actions has a PK on id, so duplicates by id shouldn't exist after
	// ON CONFLICT DO NOTHING inserts. This is a safety check.
	var dupCount int64
	if err := db.Raw(`
		SELECT COUNT(*) - COUNT(DISTINCT id) FROM user_actions
	`).Scan(&dupCount).Error; err != nil {
		return result, fmt.Errorf("check user_actions dups: %w", err)
	}

	result.DuplicatesSeen = dupCount
	if dupCount == 0 {
		log.Info().Msg("no user_actions duplicates (expected)")
		return result, nil
	}

	log.Warn().Int64("duplicates", dupCount).Msg("unexpected user_actions duplicates found")
	return result, nil
}
