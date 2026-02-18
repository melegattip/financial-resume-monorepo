package migration

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Runner orchestrates the database consolidation migration.
type Runner struct {
	TargetDB       *gorm.DB
	UsersDB        *gorm.DB
	GamificationDB *gorm.DB
	Log            zerolog.Logger
	DryRun         bool
	Phase          int // 0 = all phases, 1-4 = specific phase
}

// Run executes the full migration: audit → schema → data → dedup → validate → report.
func (r *Runner) Run() (exitCode int) {
	report := NewReport(r.DryRun)
	r.Log.Info().Bool("dry_run", r.DryRun).Int("phase", r.Phase).Msg("migration started")

	// Phase 1: Audit
	if r.Phase == 0 || r.Phase == 1 {
		r.Log.Info().Msg("phase 1: running pre-migration audit")
		auditResults, err := r.runAudit()
		if err != nil {
			r.Log.Error().Err(err).Msg("audit failed")
			report.Overall = "FAIL"
			report.Finish()
			r.printReport(report)
			return 2
		}
		report.Audit = auditResults
	}

	// Phase 2: Schema changes
	if r.Phase == 0 || r.Phase == 2 {
		r.Log.Info().Msg("phase 2: applying schema changes")
		if err := r.runSchemaChanges(); err != nil {
			r.Log.Error().Err(err).Msg("schema migration failed")
			report.Overall = "FAIL"
			report.Finish()
			r.printReport(report)
			return 2
		}
	}

	// Phase 3: Data copy + dedup
	if r.Phase == 0 || r.Phase == 3 {
		r.Log.Info().Msg("phase 3: copying and deduplicating data")
		if err := r.runDataMigration(report); err != nil {
			r.Log.Error().Err(err).Msg("data migration failed")
			report.Overall = "FAIL"
			report.Finish()
			r.printReport(report)
			return 3
		}
	}

	// Phase 4: Validation
	if r.Phase == 0 || r.Phase == 4 {
		r.Log.Info().Msg("phase 4: running post-migration validation")
		checks, err := r.runValidation()
		if err != nil {
			r.Log.Error().Err(err).Msg("validation execution failed")
			report.Overall = "FAIL"
			report.Finish()
			r.printReport(report)
			return 4
		}
		report.Validation = checks
	}

	report.Finish()
	r.printReport(report)

	if report.Overall == "FAIL" {
		return 4
	}
	return 0
}

// Audit runs only the pre-migration audit and prints results.
func (r *Runner) Audit() (exitCode int) {
	r.Log.Info().Msg("running pre-migration audit")
	results, err := r.runAudit()
	if err != nil {
		r.Log.Error().Err(err).Msg("audit failed")
		return 1
	}

	report := NewReport(false)
	report.Audit = results
	report.Finish()
	if err := report.PrintJSON(os.Stdout); err != nil {
		r.Log.Error().Err(err).Msg("failed to print audit report")
		return 1
	}
	return 0
}

// Validate runs the comprehensive post-migration validation suite.
// Exit codes: 0 = all pass, 4 = any check failed.
func (r *Runner) Validate() (exitCode int) {
	r.Log.Info().Msg("running comprehensive post-migration validation")
	checks, err := r.runValidation()
	if err != nil {
		r.Log.Error().Err(err).Msg("validation execution failed")
		return 1
	}

	report := NewReport(false)
	report.Validation = checks
	report.Finish()
	if err := report.PrintJSON(os.Stdout); err != nil {
		r.Log.Error().Err(err).Msg("failed to print validation report")
		return 1
	}
	report.PrintSummary(os.Stderr)

	if report.Overall == "FAIL" {
		return 4
	}
	return 0
}

func (r *Runner) printReport(report *Report) {
	if err := report.PrintJSON(os.Stdout); err != nil {
		r.Log.Error().Err(err).Msg("failed to print JSON report")
	}
	report.PrintSummary(os.Stderr)
}

// runSchemaChanges applies all schema alterations (column widening, table creation).
func (r *Runner) runSchemaChanges() error {
	// US0: Create users and user_preferences tables first (needed for data copy)
	r.Log.Info().Msg("creating users and user_preferences tables")
	if err := createUsersAndPreferencesTables(r.TargetDB, r.DryRun); err != nil {
		return fmt.Errorf("create users/user_preferences tables: %w", err)
	}

	// US1: Ensure users table PK and user_preferences FK are varchar(255)
	r.Log.Info().Msg("ensuring users.id is varchar(255)")
	if err := ensureUsersTableVarchar(r.TargetDB, r.DryRun); err != nil {
		return fmt.Errorf("users table schema: %w", err)
	}

	r.Log.Info().Msg("ensuring user_preferences.user_id is varchar(255)")
	if err := ensureUserPreferencesVarchar(r.TargetDB, r.DryRun); err != nil {
		return fmt.Errorf("user_preferences schema: %w", err)
	}

	// US2: Create gamification tables from gamification-db schema
	r.Log.Info().Msg("creating gamification tables via AutoMigrate")
	if err := createGamificationSchema(r.TargetDB, r.DryRun); err != nil {
		return fmt.Errorf("gamification schema: %w", err)
	}

	// US3: Standardize all user_id columns to VARCHAR(255)
	r.Log.Info().Msg("standardizing user_id columns to varchar(255)")
	if err := standardizeUserIDColumns(r.TargetDB, r.DryRun); err != nil {
		return fmt.Errorf("standardize user_id columns: %w", err)
	}

	// US4: Add soft delete columns and partial indexes
	r.Log.Info().Msg("adding deleted_at columns to all entity tables")
	if err := addSoftDeleteColumns(r.TargetDB, r.DryRun); err != nil {
		return fmt.Errorf("add soft delete columns: %w", err)
	}

	r.Log.Info().Msg("creating partial indexes for key query paths")
	if err := createPartialIndexes(r.TargetDB, r.DryRun); err != nil {
		return fmt.Errorf("create partial indexes: %w", err)
	}

	return nil
}

// runDataMigration copies data from source databases into the target.
func (r *Runner) runDataMigration(report *Report) error {
	// US1: Copy users from users_db
	if r.UsersDB != nil {
		r.Log.Info().Msg("copying users from users_db")
		copied, skipped, err := CopyUsers(r.UsersDB, r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("copy users: %w", err)
		}
		report.DataCopied["users"] = copied
		report.DataSkipped["users"] = skipped

		r.Log.Info().Msg("copying user_preferences from users_db")
		copied, skipped, err = CopyUserPreferences(r.UsersDB, r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("copy user_preferences: %w", err)
		}
		report.DataCopied["user_preferences"] = copied
		report.DataSkipped["user_preferences"] = skipped
	} else {
		r.Log.Warn().Msg("users_db not connected, skipping user data migration")
	}

	// US2: Copy gamification data from gamification-db
	if r.GamificationDB != nil {
		r.Log.Info().Msg("copying gamification core data (user_gamification, achievements, user_actions)")
		coreCopied, coreSkipped, err := CopyGamificationCore(r.GamificationDB, r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("copy gamification core: %w", err)
		}
		for k, v := range coreCopied {
			report.DataCopied[k] = v
		}
		for k, v := range coreSkipped {
			report.DataSkipped[k] = v
		}

		r.Log.Info().Msg("copying gamification tables (challenges, user_challenges, challenge_progress_tracking)")
		tblCopied, tblSkipped, err := CopyGamificationTables(r.GamificationDB, r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("copy gamification tables: %w", err)
		}
		for k, v := range tblCopied {
			report.DataCopied[k] = v
		}
		for k, v := range tblSkipped {
			report.DataSkipped[k] = v
		}

		// US2: Deduplication
		r.Log.Info().Msg("deduplicating gamification data")
		ugDedup, err := DeduplicateUserGamification(r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("dedup user_gamification: %w", err)
		}
		report.Dedup = append(report.Dedup, ugDedup)

		achDedup, err := DeduplicateAchievements(r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("dedup achievements: %w", err)
		}
		report.Dedup = append(report.Dedup, achDedup)

		uaDedup, err := DeduplicateUserActions(r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("dedup user_actions: %w", err)
		}
		report.Dedup = append(report.Dedup, uaDedup)

		// Phase 4: Copy financial data (expenses, incomes, budgets, etc.)
		r.Log.Info().Msg("copying financial data (expenses, incomes, categories, budgets, savings)")
		finCopied, finSkipped, err := CopyFinancialData(r.GamificationDB, r.TargetDB, r.Log, r.DryRun)
		if err != nil {
			return fmt.Errorf("copy financial data: %w", err)
		}
		for k, v := range finCopied {
			report.DataCopied[k] = v
		}
		for k, v := range finSkipped {
			report.DataSkipped[k] = v
		}
	} else {
		r.Log.Warn().Msg("gamification_db not connected, skipping gamification and financial migration")
	}

	return nil
}

// runValidation runs all post-migration validation checks.
func (r *Runner) runValidation() ([]ValidationCheck, error) {
	var checks []ValidationCheck

	// Determine expected user count from source.
	var expectedUsers int64
	if r.UsersDB != nil {
		if err := r.UsersDB.Raw("SELECT COUNT(*) FROM users").Scan(&expectedUsers).Error; err != nil {
			r.Log.Warn().Err(err).Msg("could not get source user count, using 0")
		}
	}

	// US1: Validate user migration
	userChecks := ValidateUsers(r.TargetDB, expectedUsers)
	checks = append(checks, userChecks...)

	// US2: Validate gamification migration
	if r.GamificationDB != nil {
		expectedCounts := make(map[string]int64)
		gamTables := []string{"user_gamification", "achievements", "user_actions", "challenges", "user_challenges", "challenge_progress_tracking"}
		for _, table := range gamTables {
			var count int64
			if err := r.GamificationDB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %q", table)).Scan(&count).Error; err != nil {
				r.Log.Warn().Err(err).Str("table", table).Msg("could not get source count")
				continue
			}
			expectedCounts[table] = count
		}
		gamChecks := ValidateGamification(r.TargetDB, expectedCounts)
		checks = append(checks, gamChecks...)
	}

	// US3: Validate user_id column types
	typeChecks := ValidateUserIDTypes(r.TargetDB)
	checks = append(checks, typeChecks...)

	// US4: Validate soft delete columns and indexes
	softDeleteChecks := ValidateSoftDelete(r.TargetDB)
	checks = append(checks, softDeleteChecks...)

	return checks, nil
}
