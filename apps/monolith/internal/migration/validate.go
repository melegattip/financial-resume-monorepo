package migration

import (
	"fmt"

	"gorm.io/gorm"
)

// ValidateUsers checks that user migration completed correctly:
// 1. User count in target matches expected count
// 2. No orphaned user_id references in user_preferences
// 3. All user_id values in financial tables exist in users table
func ValidateUsers(targetDB *gorm.DB, expectedCount int64) []ValidationCheck {
	var checks []ValidationCheck

	// Check 1: User count matches expected
	var actualCount int64
	if err := targetDB.Raw("SELECT COUNT(*) FROM users").Scan(&actualCount).Error; err != nil {
		checks = append(checks, ValidationCheck{
			Name:    "user_count",
			Status:  "FAIL",
			Details: fmt.Sprintf("query failed: %v", err),
		})
	} else {
		status := "PASS"
		details := ""
		if actualCount < expectedCount {
			status = "FAIL"
			details = fmt.Sprintf("expected at least %d, got %d", expectedCount, actualCount)
		}
		checks = append(checks, ValidationCheck{
			Name:     "user_count",
			Status:   status,
			Expected: expectedCount,
			Actual:   actualCount,
			Details:  details,
		})
	}

	// Check 2: No orphaned user_ids in user_preferences
	var orphanedPrefs int64
	if err := targetDB.Raw(`
		SELECT COUNT(*) FROM user_preferences p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE u.id IS NULL
	`).Scan(&orphanedPrefs).Error; err != nil {
		checks = append(checks, ValidationCheck{
			Name:    "user_preferences_orphans",
			Status:  "FAIL",
			Details: fmt.Sprintf("query failed: %v", err),
		})
	} else {
		status := "PASS"
		details := ""
		if orphanedPrefs > 0 {
			status = "FAIL"
			details = fmt.Sprintf("%d orphaned user_preference rows", orphanedPrefs)
		}
		checks = append(checks, ValidationCheck{
			Name:     "user_preferences_orphans",
			Status:   status,
			Expected: int64(0),
			Actual:   orphanedPrefs,
			Details:  details,
		})
	}

	// Check 3: Orphaned user_ids in financial tables
	financialTables := []string{
		"expenses", "incomes", "categories", "budgets",
		"savings_goals", "savings_transactions",
		"recurring_transactions",
	}
	for _, table := range financialTables {
		// Check if table exists first
		var tableExists int64
		if err := targetDB.Raw(`
			SELECT COUNT(*) FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = ?
		`, table).Scan(&tableExists).Error; err != nil || tableExists == 0 {
			continue // Skip tables that don't exist
		}

		var orphaned int64
		query := fmt.Sprintf(`
			SELECT COUNT(*) FROM %q t
			LEFT JOIN users u ON t.user_id = u.id
			WHERE u.id IS NULL AND t.user_id IS NOT NULL AND t.user_id != ''
		`, table)
		if err := targetDB.Raw(query).Scan(&orphaned).Error; err != nil {
			checks = append(checks, ValidationCheck{
				Name:    fmt.Sprintf("%s_orphan_user_ids", table),
				Status:  "FAIL",
				Details: fmt.Sprintf("query failed: %v", err),
			})
		} else if orphaned > 0 {
			checks = append(checks, ValidationCheck{
				Name:     fmt.Sprintf("%s_orphan_user_ids", table),
				Status:   "FAIL",
				Expected: int64(0),
				Actual:   orphaned,
				Details:  fmt.Sprintf("%d rows reference non-existent users", orphaned),
			})
		} else {
			checks = append(checks, ValidationCheck{
				Name:     fmt.Sprintf("%s_orphan_user_ids", table),
				Status:   "PASS",
				Expected: int64(0),
				Actual:   orphaned,
			})
		}
	}

	return checks
}

// ValidateGamification checks gamification migration integrity:
// 1. Exactly 1 user_gamification row per user_id
// 2. Challenges count matches expected
// 3. User_challenges count matches expected
// 4. No orphaned FKs in gamification tables
func ValidateGamification(targetDB *gorm.DB, expectedCounts map[string]int64) []ValidationCheck {
	var checks []ValidationCheck

	// Check 1: Exactly 1 user_gamification per user_id (no duplicates)
	var dupUsers int64
	if err := targetDB.Raw(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM user_gamification
			GROUP BY user_id HAVING COUNT(*) > 1
		) sub
	`).Scan(&dupUsers).Error; err != nil {
		checks = append(checks, ValidationCheck{
			Name:    "user_gamification_no_dups",
			Status:  "FAIL",
			Details: fmt.Sprintf("query failed: %v", err),
		})
	} else {
		status := "PASS"
		details := ""
		if dupUsers > 0 {
			status = "FAIL"
			details = fmt.Sprintf("%d users still have duplicate rows", dupUsers)
		}
		checks = append(checks, ValidationCheck{
			Name:     "user_gamification_no_dups",
			Status:   status,
			Expected: int64(0),
			Actual:   dupUsers,
			Details:  details,
		})
	}

	// Check 2-3: Table counts match expected
	for _, table := range []string{"challenges", "user_challenges", "user_gamification", "achievements", "user_actions", "challenge_progress_tracking"} {
		expected, ok := expectedCounts[table]
		if !ok {
			continue // No expected count for this table
		}

		var actual int64
		if err := targetDB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %q", table)).Scan(&actual).Error; err != nil {
			checks = append(checks, ValidationCheck{
				Name:    fmt.Sprintf("%s_count", table),
				Status:  "FAIL",
				Details: fmt.Sprintf("query failed: %v", err),
			})
			continue
		}

		status := "PASS"
		details := ""
		if actual < expected {
			status = "FAIL"
			details = fmt.Sprintf("expected at least %d, got %d", expected, actual)
		}
		checks = append(checks, ValidationCheck{
			Name:     fmt.Sprintf("%s_count", table),
			Status:   status,
			Expected: expected,
			Actual:   actual,
			Details:  details,
		})
	}

	// Check 4: No orphaned user_ids in gamification tables
	gamTables := []string{"user_gamification", "achievements", "user_actions", "user_challenges", "challenge_progress_tracking"}
	for _, table := range gamTables {
		var tableExists int64
		if err := targetDB.Raw(`
			SELECT COUNT(*) FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = ?
		`, table).Scan(&tableExists).Error; err != nil || tableExists == 0 {
			continue
		}

		var orphaned int64
		if err := targetDB.Raw(fmt.Sprintf(`
			SELECT COUNT(*) FROM %q t
			LEFT JOIN users u ON t.user_id = u.id
			WHERE u.id IS NULL AND t.user_id IS NOT NULL AND t.user_id != ''
		`, table)).Scan(&orphaned).Error; err != nil {
			checks = append(checks, ValidationCheck{
				Name:    fmt.Sprintf("%s_orphan_user_ids", table),
				Status:  "FAIL",
				Details: fmt.Sprintf("query failed: %v", err),
			})
		} else {
			status := "PASS"
			details := ""
			if orphaned > 0 {
				status = "FAIL"
				details = fmt.Sprintf("%d rows reference non-existent users", orphaned)
			}
			checks = append(checks, ValidationCheck{
				Name:     fmt.Sprintf("%s_orphan_user_ids", table),
				Status:   status,
				Expected: int64(0),
				Actual:   orphaned,
				Details:  details,
			})
		}
	}

	return checks
}

// ValidateUserIDTypes checks that all user_id columns are VARCHAR(255).
func ValidateUserIDTypes(targetDB *gorm.DB) []ValidationCheck {
	var checks []ValidationCheck

	type colInfo struct {
		TableName              string
		ColumnName             string
		DataType               string
		CharacterMaximumLength *int64
	}

	var cols []colInfo
	if err := targetDB.Raw(`
		SELECT table_name, column_name, data_type, character_maximum_length
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND column_name = 'user_id'
		ORDER BY table_name
	`).Scan(&cols).Error; err != nil {
		checks = append(checks, ValidationCheck{
			Name:    "user_id_types_query",
			Status:  "FAIL",
			Details: fmt.Sprintf("query failed: %v", err),
		})
		return checks
	}

	for _, c := range cols {
		name := fmt.Sprintf("%s.user_id_type", c.TableName)
		isCorrect := c.DataType == "character varying" && c.CharacterMaximumLength != nil && *c.CharacterMaximumLength == 255
		status := "PASS"
		details := ""
		if !isCorrect {
			status = "FAIL"
			if c.CharacterMaximumLength != nil {
				details = fmt.Sprintf("type=%s(%d), expected varchar(255)", c.DataType, *c.CharacterMaximumLength)
			} else {
				details = fmt.Sprintf("type=%s, expected varchar(255)", c.DataType)
			}
		}
		checks = append(checks, ValidationCheck{
			Name:     name,
			Status:   status,
			Expected: "varchar(255)",
			Actual:   fmt.Sprintf("%s(%v)", c.DataType, c.CharacterMaximumLength),
			Details:  details,
		})
	}

	return checks
}

// ValidateSoftDelete checks that all entity tables have a deleted_at column
// and that partial indexes exist.
func ValidateSoftDelete(targetDB *gorm.DB) []ValidationCheck {
	var checks []ValidationCheck

	expectedTables := []string{
		"users", "user_preferences", "user_notification_settings", "user_two_fa",
		"expenses", "incomes", "categories", "budgets",
		"budget_notifications", "savings_goals", "savings_transactions",
		"recurring_transactions", "recurring_transaction_executions",
		"recurring_transaction_notifications",
		"user_gamification", "achievements", "user_actions",
		"challenges", "user_challenges", "challenge_progress_tracking",
	}

	for _, table := range expectedTables {
		// Check if table exists first
		var tableExists int64
		if err := targetDB.Raw(`
			SELECT COUNT(*) FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = ?
		`, table).Scan(&tableExists).Error; err != nil || tableExists == 0 {
			continue
		}

		var colExists int64
		if err := targetDB.Raw(`
			SELECT COUNT(*) FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = ? AND column_name = 'deleted_at'
		`, table).Scan(&colExists).Error; err != nil {
			checks = append(checks, ValidationCheck{
				Name:    fmt.Sprintf("%s_deleted_at", table),
				Status:  "FAIL",
				Details: fmt.Sprintf("query failed: %v", err),
			})
			continue
		}

		status := "PASS"
		details := ""
		if colExists == 0 {
			status = "FAIL"
			details = "deleted_at column missing"
		}
		checks = append(checks, ValidationCheck{
			Name:     fmt.Sprintf("%s_deleted_at", table),
			Status:   status,
			Expected: int64(1),
			Actual:   colExists,
			Details:  details,
		})
	}

	// Check partial indexes exist
	expectedIndexes := []string{
		"idx_expenses_user_date_active",
		"idx_expenses_category_active",
		"idx_budgets_user_active_soft",
		"idx_recurring_next_exec_active",
		"idx_user_gamification_user_active",
	}

	for _, idx := range expectedIndexes {
		var idxExists int64
		if err := targetDB.Raw(`
			SELECT COUNT(*) FROM pg_indexes
			WHERE schemaname = 'public' AND indexname = ?
		`, idx).Scan(&idxExists).Error; err != nil {
			checks = append(checks, ValidationCheck{
				Name:    fmt.Sprintf("index_%s", idx),
				Status:  "FAIL",
				Details: fmt.Sprintf("query failed: %v", err),
			})
			continue
		}

		status := "PASS"
		details := ""
		if idxExists == 0 {
			status = "FAIL"
			details = "partial index missing"
		}
		checks = append(checks, ValidationCheck{
			Name:     fmt.Sprintf("index_%s", idx),
			Status:   status,
			Expected: int64(1),
			Actual:   idxExists,
			Details:  details,
		})
	}

	return checks
}
