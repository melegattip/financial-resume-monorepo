package database

import (
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// tenantScopedTables are the tables that hold per-tenant user data and must be
// protected by Row-Level Security.
var tenantScopedTables = []string{
	"expenses",
	"incomes",
	"categories",
	"budgets",
	"savings_goals",
	"savings_transactions",
	"recurring_transactions",
	"audit_logs",
}

// MigrateRLSPolicies enables PostgreSQL Row-Level Security on all tenant-scoped
// tables and creates PERMISSIVE policies that:
//   - Allow all rows when app.current_tenant is not set (empty/unset), so that
//     migrations and admin operations continue to work normally.
//   - Restrict rows to tenant_id = app.current_tenant when the setting is present,
//     providing DB-level defence-in-depth on top of the application-level filtering.
//
// The function is idempotent — safe to call on every startup.
//
// Usage:
//
//	database.MigrateRLSPolicies(db, logger)
//
// To enforce the policy in a handler, wrap DB operations with WithTenant:
//
//	err := database.WithTenant(db, tenantID, func(tx *gorm.DB) error {
//	    return tx.Find(&expenses).Error
//	})
func MigrateRLSPolicies(db *gorm.DB, logger zerolog.Logger) {
	for _, table := range tenantScopedTables {
		// Enable RLS — idempotent if already enabled.
		enableSQL := fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", table) //nolint:gosec
		if err := db.Exec(enableSQL).Error; err != nil {
			logger.Warn().Err(err).Str("table", table).Msg("rls: enable RLS warning")
		}

		// Create the policy only if it doesn't exist yet.
		// PostgreSQL has no "CREATE POLICY IF NOT EXISTS", so we check pg_policies.
		//
		// Policy logic:
		//   - empty/unset setting  → allow all rows (migrations, admin tools)
		//   - setting is present   → restrict to matching tenant_id
		//nolint:gosec
		policySQL := fmt.Sprintf(`
			DO $$ BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_policies
					WHERE schemaname = 'public'
					  AND tablename  = '%s'
					  AND policyname = 'tenant_isolation'
				) THEN
					CREATE POLICY tenant_isolation ON %s
						AS PERMISSIVE FOR ALL
						USING (
							current_setting('app.current_tenant', true) = ''
							OR tenant_id = current_setting('app.current_tenant', true)
						);
				END IF;
			END $$;
		`, table, table)

		if err := db.Exec(policySQL).Error; err != nil {
			logger.Warn().Err(err).Str("table", table).Msg("rls: create policy warning")
		} else {
			logger.Info().Str("table", table).Msg("rls: policy ensured")
		}
	}

	logger.Info().Msg("rls: row-level security policies ensured")
}

// WithTenant executes fn inside a database transaction where the PostgreSQL
// session variable app.current_tenant is set to tenantID via SET LOCAL.
//
// Because SET LOCAL scopes the variable to the current transaction, this
// guarantees that the RLS policy sees the correct tenant on the same DB
// connection as the queries in fn, avoiding connection-pool leakage issues.
//
// Usage:
//
//	err := database.WithTenant(db, tenantID, func(tx *gorm.DB) error {
//	    return tx.Where("deleted_at IS NULL").Find(&expenses).Error
//	})
func WithTenant(db *gorm.DB, tenantID string, fn func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SET LOCAL app.current_tenant = ?", tenantID).Error; err != nil {
			return fmt.Errorf("rls: set tenant context: %w", err)
		}
		return fn(tx)
	})
}
