package migration

import (
	"fmt"

	"gorm.io/gorm"
)

// createUsersAndPreferencesTables creates the users and user_preferences tables
// using direct SQL to ensure correct table names. This is needed before copying data.
func createUsersAndPreferencesTables(db *gorm.DB, dryRun bool) error {
	if dryRun {
		fmt.Println("[DRY RUN] Would create users and user_preferences tables")
		return nil
	}

	// Create users table directly with SQL
	createUsersSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			phone VARCHAR(50),
			avatar VARCHAR(500),
			is_active BOOLEAN DEFAULT true,
			is_verified BOOLEAN DEFAULT false,
			last_login TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		)
	`
	if err := db.Exec(createUsersSQL).Error; err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	// Create user_preferences table directly with SQL
	createPrefsSQL := `
		CREATE TABLE IF NOT EXISTS user_preferences (
			user_id VARCHAR(255) PRIMARY KEY,
			currency VARCHAR(10) DEFAULT 'USD',
			language VARCHAR(10) DEFAULT 'en',
			timezone VARCHAR(50) DEFAULT 'UTC',
			gamification_enabled BOOLEAN DEFAULT true,
			theme VARCHAR(20) DEFAULT 'light',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			CONSTRAINT fk_user_preferences_user FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`
	if err := db.Exec(createPrefsSQL).Error; err != nil {
		return fmt.Errorf("create user_preferences table: %w", err)
	}

	return nil
}

// ensureUsersTableVarchar verifies the users table PK is varchar(255).
// If it was previously a uuid type (from GORM AutoMigrate with uuid.UUID),
// it alters the column to varchar(255). If the table doesn't exist yet,
// GORM AutoMigrate will create it correctly from the updated User model.
func ensureUsersTableVarchar(db *gorm.DB, dryRun bool) error {
	// Check the current type of users.id
	type colType struct {
		DataType               string
		CharacterMaximumLength *int64
		UdtName                string
	}
	var ct colType
	err := db.Raw(`
		SELECT data_type, character_maximum_length, udt_name
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'users'
		  AND column_name = 'id'
	`).Scan(&ct).Error
	if err != nil {
		return fmt.Errorf("query users.id type: %w", err)
	}

	// If no result, table doesn't exist yet — AutoMigrate will handle it.
	if ct.DataType == "" {
		return nil
	}

	// If already varchar(255), nothing to do.
	if ct.DataType == "character varying" && ct.CharacterMaximumLength != nil && *ct.CharacterMaximumLength >= 255 {
		return nil
	}

	// Need to alter the column type.
	sql := `ALTER TABLE users ALTER COLUMN id TYPE VARCHAR(255)`
	if dryRun {
		fmt.Printf("[DRY RUN] Would execute: %s\n", sql)
		return nil
	}
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("alter users.id to varchar(255): %w", err)
	}
	return nil
}

// ensureUserPreferencesVarchar verifies user_preferences.user_id is varchar(255).
func ensureUserPreferencesVarchar(db *gorm.DB, dryRun bool) error {
	type colType struct {
		DataType               string
		CharacterMaximumLength *int64
	}
	var ct colType
	err := db.Raw(`
		SELECT data_type, character_maximum_length
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'user_preferences'
		  AND column_name = 'user_id'
	`).Scan(&ct).Error
	if err != nil {
		return fmt.Errorf("query user_preferences.user_id type: %w", err)
	}

	if ct.DataType == "" {
		return nil // Table doesn't exist yet
	}

	if ct.DataType == "character varying" && ct.CharacterMaximumLength != nil && *ct.CharacterMaximumLength >= 255 {
		return nil
	}

	sql := `ALTER TABLE user_preferences ALTER COLUMN user_id TYPE VARCHAR(255)`
	if dryRun {
		fmt.Printf("[DRY RUN] Would execute: %s\n", sql)
		return nil
	}
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("alter user_preferences.user_id to varchar(255): %w", err)
	}
	return nil
}

// createGamificationSchema drops old gamification tables from main-db (if they
// exist with the old XP-based schema) and creates the new gamification-db-style
// tables via GORM AutoMigrate.
func createGamificationSchema(db *gorm.DB, dryRun bool) error {
	if dryRun {
		fmt.Println("[DRY RUN] Would AutoMigrate gamification target tables")
		return nil
	}

	// AutoMigrate creates or updates tables to match the target models.
	// For existing tables it adds missing columns without dropping anything.
	return db.AutoMigrate(
		&TgtUserGamification{},
		&TgtAchievement{},
		&TgtUserAction{},
		&TgtChallenge{},
		&TgtUserChallenge{},
		&TgtChallengeProgressTracking{},
	)
}

// --- US3: Standardize user_id columns to VARCHAR(255) ---

// standardizeUserIDColumns widens all user_id columns across financial and
// gamification tables to VARCHAR(255).
func standardizeUserIDColumns(db *gorm.DB, dryRun bool) error {
	// Financial tables that need user_id widened
	financialTables := []string{
		"expenses", "incomes", "categories", "budgets",
		"budget_notifications", "savings_goals", "savings_transactions",
		"recurring_transactions", "recurring_transaction_executions",
		"recurring_transaction_notifications",
	}

	// Gamification tables (should already be varchar(255) from AutoMigrate, but verify)
	gamificationTables := []string{
		"user_gamification", "achievements", "user_actions",
		"user_challenges", "challenge_progress_tracking",
	}

	allTables := append(financialTables, gamificationTables...)

	for _, table := range allTables {
		// Check if table exists
		var tableExists int64
		if err := db.Raw(`
			SELECT COUNT(*) FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = ?
		`, table).Scan(&tableExists).Error; err != nil || tableExists == 0 {
			continue
		}

		// Check if user_id column exists and its current type
		type colType struct {
			DataType               string
			CharacterMaximumLength *int64
		}
		var ct colType
		if err := db.Raw(`
			SELECT data_type, character_maximum_length
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = ? AND column_name = 'user_id'
		`, table).Scan(&ct).Error; err != nil || ct.DataType == "" {
			continue // Column doesn't exist
		}

		// Already wide enough
		if ct.DataType == "character varying" && ct.CharacterMaximumLength != nil && *ct.CharacterMaximumLength >= 255 {
			continue
		}

		sql := fmt.Sprintf("ALTER TABLE %q ALTER COLUMN user_id TYPE VARCHAR(255)", table)
		if dryRun {
			fmt.Printf("[DRY RUN] Would execute: %s\n", sql)
			continue
		}
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("alter %s.user_id: %w", table, err)
		}
	}
	return nil
}

// --- US4: Add soft delete columns and partial indexes ---

// addSoftDeleteColumns adds "deleted_at TIMESTAMP NULL" to all entity tables.
func addSoftDeleteColumns(db *gorm.DB, dryRun bool) error {
	tables := []string{
		"users", "user_preferences", "user_notification_settings", "user_two_fa",
		"expenses", "incomes", "categories", "budgets",
		"budget_notifications", "savings_goals", "savings_transactions",
		"recurring_transactions", "recurring_transaction_executions",
		"recurring_transaction_notifications",
		"user_gamification", "achievements", "user_actions",
		"challenges", "user_challenges", "challenge_progress_tracking",
	}

	for _, table := range tables {
		// Check if table exists
		var tableExists int64
		if err := db.Raw(`
			SELECT COUNT(*) FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = ?
		`, table).Scan(&tableExists).Error; err != nil || tableExists == 0 {
			continue
		}

		sql := fmt.Sprintf("ALTER TABLE %q ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL", table)
		if dryRun {
			fmt.Printf("[DRY RUN] Would execute: %s\n", sql)
			continue
		}
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("add deleted_at to %s: %w", table, err)
		}
	}
	return nil
}

// createPartialIndexes creates performance-optimized partial indexes with
// WHERE deleted_at IS NULL for key query paths.
func createPartialIndexes(db *gorm.DB, dryRun bool) error {
	indexes := []struct {
		name    string
		table   string
		columns string
	}{
		{"idx_expenses_user_date_active", "expenses", "user_id, transaction_date"},
		{"idx_expenses_category_active", "expenses", "category_id"},
		{"idx_budgets_user_active_soft", "budgets", "user_id, is_active"},
		{"idx_recurring_next_exec_active", "recurring_transactions", "next_date, is_active"},
		{"idx_user_gamification_user_active", "user_gamification", "user_id"},
	}

	for _, idx := range indexes {
		// Check if table exists
		var tableExists int64
		if err := db.Raw(`
			SELECT COUNT(*) FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = ?
		`, idx.table).Scan(&tableExists).Error; err != nil || tableExists == 0 {
			continue
		}

		sql := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %q (%s) WHERE deleted_at IS NULL",
			idx.name, idx.table, idx.columns)
		if dryRun {
			fmt.Printf("[DRY RUN] Would execute: %s\n", sql)
			continue
		}
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("create index %s: %w", idx.name, err)
		}
	}
	return nil
}
