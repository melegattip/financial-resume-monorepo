package migration

import (
	"fmt"

	"gorm.io/gorm"
)

// runAudit performs a pre-migration audit across all three databases.
// It collects table counts, identifies orphaned references, duplicate
// gamification rows, and user_id column types.
func (r *Runner) runAudit() ([]AuditResult, error) {
	var results []AuditResult

	// 1. Audit target database (financial_resume / main-db)
	targetResult, err := auditDatabase(r.TargetDB, "target (financial_resume)")
	if err != nil {
		return nil, fmt.Errorf("audit target db: %w", err)
	}
	r.Log.Info().Str("database", targetResult.Database).Int("tables", len(targetResult.Tables)).Msg("audit complete")
	results = append(results, targetResult)

	// 2. Audit users-db (if connected)
	if r.UsersDB != nil {
		usersResult, err := auditDatabase(r.UsersDB, "source (users_db)")
		if err != nil {
			return nil, fmt.Errorf("audit users db: %w", err)
		}
		r.Log.Info().Str("database", usersResult.Database).Int("tables", len(usersResult.Tables)).Msg("audit complete")
		results = append(results, usersResult)
	}

	// 3. Audit gamification-db (if connected)
	if r.GamificationDB != nil {
		gamResult, err := auditDatabase(r.GamificationDB, "source (gamification_db)")
		if err != nil {
			return nil, fmt.Errorf("audit gamification db: %w", err)
		}
		r.Log.Info().Str("database", gamResult.Database).Int("tables", len(gamResult.Tables)).Msg("audit complete")
		results = append(results, gamResult)
	}

	// 4. Check user_id column types in target database
	colTypeResult := auditColumnTypes(r.TargetDB)
	results = append(results, colTypeResult)

	// 5. Check gamification duplicates in target database
	dupResult, err := auditGamificationDuplicates(r.TargetDB)
	if err != nil {
		r.Log.Warn().Err(err).Msg("skipping gamification duplicate check (tables may not exist yet)")
	} else {
		results = append(results, dupResult)
	}

	return results, nil
}

// auditDatabase queries all public tables and their row counts.
func auditDatabase(db *gorm.DB, name string) (AuditResult, error) {
	result := AuditResult{
		Database: name,
		Counts:   make(map[string]int64),
		Issues:   make(map[string]string),
	}

	// List all user-facing tables in the public schema.
	type tableRow struct {
		TableName string
	}
	var tables []tableRow
	if err := db.Raw(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		  AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`).Scan(&tables).Error; err != nil {
		return result, fmt.Errorf("list tables: %w", err)
	}

	for _, t := range tables {
		result.Tables = append(result.Tables, t.TableName)
		var count int64
		if err := db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %q", t.TableName)).Scan(&count).Error; err != nil {
			result.Issues[t.TableName] = fmt.Sprintf("count failed: %v", err)
			continue
		}
		result.Counts[t.TableName] = count
	}

	return result, nil
}

// auditColumnTypes checks id and user_id column types across all target tables.
func auditColumnTypes(db *gorm.DB) AuditResult {
	result := AuditResult{
		Database: "column_types",
		Counts:   make(map[string]int64),
		Issues:   make(map[string]string),
	}

	type colInfo struct {
		TableName              string
		ColumnName             string
		DataType               string
		CharacterMaximumLength *int64
	}

	var cols []colInfo
	if err := db.Raw(`
		SELECT table_name, column_name, data_type,
		       character_maximum_length
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND column_name IN ('id', 'user_id')
		ORDER BY table_name, column_name
	`).Scan(&cols).Error; err != nil {
		result.Issues["query"] = fmt.Sprintf("failed: %v", err)
		return result
	}

	for _, c := range cols {
		key := fmt.Sprintf("%s.%s", c.TableName, c.ColumnName)
		if c.DataType == "uuid" {
			result.Issues[key] = fmt.Sprintf("type=uuid (needs migration to varchar)")
		} else if c.DataType == "character varying" && c.CharacterMaximumLength != nil && *c.CharacterMaximumLength < 255 {
			result.Issues[key] = fmt.Sprintf("type=varchar(%d) (needs widening to 255)", *c.CharacterMaximumLength)
		}
		if c.CharacterMaximumLength != nil {
			result.Counts[key] = *c.CharacterMaximumLength
		}
	}

	return result
}

// auditGamificationDuplicates checks for duplicate rows in gamification tables.
func auditGamificationDuplicates(db *gorm.DB) (AuditResult, error) {
	result := AuditResult{
		Database: "gamification_duplicates",
		Counts:   make(map[string]int64),
		Issues:   make(map[string]string),
	}

	// Check if user_gamification table exists before querying.
	var tableExists int64
	if err := db.Raw(`
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'user_gamification'
	`).Scan(&tableExists).Error; err != nil {
		return result, fmt.Errorf("check table existence: %w", err)
	}
	if tableExists == 0 {
		result.Issues["user_gamification"] = "table does not exist"
		return result, nil
	}

	// Check user_gamification duplicates
	var ugDups int64
	if err := db.Raw(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM user_gamification
			GROUP BY user_id HAVING COUNT(*) > 1
		) sub
	`).Scan(&ugDups).Error; err != nil {
		return result, fmt.Errorf("check user_gamification dups: %w", err)
	}
	result.Counts["user_gamification_dup_users"] = ugDups
	if ugDups > 0 {
		result.Issues["user_gamification"] = fmt.Sprintf("%d users with duplicate rows", ugDups)
	}

	// Check achievements duplicates (same user_id + type)
	var achTableExists int64
	if err := db.Raw(`
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'achievements'
	`).Scan(&achTableExists).Error; err == nil && achTableExists > 0 {
		var achDups int64
		if err := db.Raw(`
			SELECT COUNT(*) FROM (
				SELECT user_id, type FROM achievements
				GROUP BY user_id, type HAVING COUNT(*) > 1
			) sub
		`).Scan(&achDups).Error; err == nil {
			result.Counts["achievements_dup_pairs"] = achDups
			if achDups > 0 {
				result.Issues["achievements"] = fmt.Sprintf("%d duplicate user+type pairs", achDups)
			}
		}
	}

	return result, nil
}
