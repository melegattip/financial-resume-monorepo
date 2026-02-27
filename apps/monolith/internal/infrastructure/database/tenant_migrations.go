package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MigrateTenantColumns adds the tenant_id column (idempotent via ADD COLUMN IF NOT EXISTS)
// to all tables that require tenant-scoping in the multi-tenant model.
func MigrateTenantColumns(db *gorm.DB, logger zerolog.Logger) {
	tables := []string{
		"expenses", "incomes", "categories",
		"budgets", "savings_goals", "savings_transactions",
		"recurring_transactions", "user_gamification",
		"achievements", "user_actions",
	}

	for _, table := range tables {
		sql := fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50)",
			table,
		)
		if err := db.Exec(sql).Error; err != nil {
			logger.Warn().Err(err).Str("table", table).Msg("tenant_id column migration warning")
		} else {
			logger.Info().Str("table", table).Msg("tenant_id column ensured")
		}
	}

	// Indexes for common tenant-scoped queries
	indexes := []struct{ table, name, cols string }{
		{"expenses", "idx_expenses_tenant", "tenant_id"},
		{"incomes", "idx_incomes_tenant", "tenant_id"},
		{"categories", "idx_categories_tenant", "tenant_id"},
		{"budgets", "idx_budgets_tenant", "tenant_id"},
		{"savings_goals", "idx_savings_goals_tenant", "tenant_id"},
		{"savings_transactions", "idx_savings_transactions_tenant", "tenant_id"},
		{"recurring_transactions", "idx_recurring_tenant", "tenant_id"},
		{"user_gamification", "idx_gamification_tenant", "tenant_id"},
		{"achievements", "idx_achievements_tenant", "tenant_id"},
		{"user_actions", "idx_user_actions_tenant", "tenant_id"},
	}

	for _, idx := range indexes {
		sql := fmt.Sprintf(
			"CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
			idx.name, idx.table, idx.cols,
		)
		if err := db.Exec(sql).Error; err != nil {
			logger.Warn().Err(err).Str("index", idx.name).Msg("index creation warning")
		}
	}

	logger.Info().Msg("tenant columns and indexes ensured")
}

// defaultPermissions is the global permissions catalog.
var defaultPermissions = []struct{ Key, Description, Category string }{
	{"view_data", "View all tenant data", "data"},
	{"create_transaction", "Create expenses and incomes", "data"},
	{"edit_any_transaction", "Edit any transaction in the tenant", "data"},
	{"delete_any_transaction", "Delete any transaction in the tenant", "data"},
	{"manage_budgets", "Create, edit and delete budgets", "data"},
	{"manage_savings", "Create, edit and delete savings goals", "data"},
	{"manage_recurring", "Create, edit and delete recurring transactions", "data"},
	{"invite_members", "Generate and revoke invitation codes", "member_management"},
	{"manage_roles", "Change member roles", "member_management"},
	{"remove_members", "Remove members from the tenant", "member_management"},
	{"view_audit_logs", "View audit logs", "admin"},
	{"manage_tenant", "Edit tenant settings", "settings"},
	{"delete_tenant", "Delete the tenant", "admin"},
	{"transfer_ownership", "Transfer tenant ownership", "admin"},
}

// roleDefaultPermissions maps each role to its default permission keys.
var roleDefaultPermissions = map[string][]string{
	"owner": {
		"view_data", "create_transaction", "edit_any_transaction", "delete_any_transaction",
		"manage_budgets", "manage_savings", "manage_recurring",
		"invite_members", "manage_roles", "remove_members",
		"view_audit_logs", "manage_tenant", "delete_tenant", "transfer_ownership",
	},
	"admin": {
		"view_data", "create_transaction", "edit_any_transaction", "delete_any_transaction",
		"manage_budgets", "manage_savings", "manage_recurring",
		"invite_members", "manage_roles", "remove_members",
		"view_audit_logs", "manage_tenant",
	},
	"member": {
		"view_data", "create_transaction", "edit_any_transaction", "delete_any_transaction",
		"manage_budgets", "manage_savings", "manage_recurring",
	},
	"viewer": {
		"view_data",
	},
}

// MigrateExistingDataToTenants creates a personal tenant for each existing user,
// seeds the permissions catalog, and populates tenant_id on all existing records.
// This function is idempotent — safe to run multiple times.
func MigrateExistingDataToTenants(db *gorm.DB, logger zerolog.Logger) {
	logger.Info().Msg("starting tenant data migration")

	// 1. Seed global permissions catalog (idempotent)
	seedPermissions(db, logger)

	// 2. Create a personal tenant for each user that doesn't have one yet
	type userRow struct {
		ID        string
		Email     string
		FirstName string
		LastName  string
	}

	var users []userRow
	if err := db.Raw(`
		SELECT u.id, u.email, u.first_name, u.last_name
		FROM users u
		WHERE u.deleted_at IS NULL
		  AND NOT EXISTS (
		    SELECT 1 FROM tenants t WHERE t.owner_id = u.id AND t.deleted_at IS NULL
		  )
	`).Scan(&users).Error; err != nil {
		logger.Error().Err(err).Msg("failed to query users for tenant migration")
		return
	}

	logger.Info().Int("count", len(users)).Msg("users without personal tenant found")

	for _, u := range users {
		tenantID := "tnt_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
		slug := buildTenantSlug(u.Email)

		name := strings.TrimSpace(u.FirstName + " " + u.LastName)
		if name == "" {
			name = u.Email
		}
		name = name + " (personal)"

		// Create tenant
		if err := db.Exec(`
			INSERT INTO tenants (id, name, slug, owner_id, is_active, plan, created_at, updated_at)
			VALUES (?, ?, ?, ?, true, 'free', NOW(), NOW())
			ON CONFLICT DO NOTHING
		`, tenantID, name, slug, u.ID).Error; err != nil {
			logger.Warn().Err(err).Str("user_id", u.ID).Msg("failed to create personal tenant")
			continue
		}

		// Register as owner in tenant_members
		memberID := "tmb_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
		if err := db.Exec(`
			INSERT INTO tenant_members (id, tenant_id, user_id, role, joined_at, created_at)
			VALUES (?, ?, ?, 'owner', NOW(), NOW())
			ON CONFLICT (tenant_id, user_id) DO NOTHING
		`, memberID, tenantID, u.ID).Error; err != nil {
			logger.Warn().Err(err).Str("user_id", u.ID).Msg("failed to create tenant membership")
		}

		// Seed default role_permissions for this new tenant
		seedTenantRolePermissions(db, tenantID, logger)

		logger.Info().Str("user_id", u.ID).Str("tenant_id", tenantID).Msg("personal tenant created")
	}

	// 3. Backfill tenant_id on all existing records (set to the owner's personal tenant)
	tenantedTables := []struct{ table, userCol string }{
		{"expenses", "user_id"},
		{"incomes", "user_id"},
		{"categories", "user_id"},
		{"budgets", "user_id"},
		{"savings_goals", "user_id"},
		{"savings_transactions", "user_id"},
		{"recurring_transactions", "user_id"},
		{"user_gamification", "user_id"},
		{"achievements", "user_id"},
		{"user_actions", "user_id"},
	}

	for _, t := range tenantedTables {
		sql := fmt.Sprintf(`
			UPDATE %s rec
			SET tenant_id = ten.id
			FROM tenants ten
			WHERE ten.owner_id = rec.%s
			  AND rec.tenant_id IS NULL
			  AND ten.deleted_at IS NULL
		`, t.table, t.userCol)

		result := db.Exec(sql)
		if result.Error != nil {
			logger.Warn().Err(result.Error).Str("table", t.table).Msg("tenant_id backfill warning")
		} else {
			logger.Info().Str("table", t.table).Int64("rows", result.RowsAffected).Msg("tenant_id backfilled")
		}
	}

	logger.Info().Msg("tenant data migration complete")
}

func seedPermissions(db *gorm.DB, logger zerolog.Logger) {
	for _, p := range defaultPermissions {
		if err := db.Exec(`
			INSERT INTO permissions (key, description, category)
			VALUES (?, ?, ?)
			ON CONFLICT (key) DO NOTHING
		`, p.Key, p.Description, p.Category).Error; err != nil {
			logger.Warn().Err(err).Str("key", p.Key).Msg("permission seed warning")
		}
	}
	logger.Info().Msg("permissions catalog seeded")
}

func seedTenantRolePermissions(db *gorm.DB, tenantID string, logger zerolog.Logger) {
	for role, perms := range roleDefaultPermissions {
		for _, perm := range perms {
			id := "rp_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
			if err := db.Exec(`
				INSERT INTO role_permissions (id, tenant_id, role, permission_key, created_at)
				VALUES (?, ?, ?, ?, NOW())
				ON CONFLICT (tenant_id, role, permission_key) DO NOTHING
			`, id, tenantID, role, perm).Error; err != nil {
				logger.Warn().Err(err).
					Str("tenant_id", tenantID).
					Str("role", role).
					Str("perm", perm).
					Msg("role_permission seed warning")
			}
		}
	}
}

func buildTenantSlug(email string) string {
	slug := strings.ToLower(email)
	slug = strings.ReplaceAll(slug, "@", "-at-")
	slug = strings.ReplaceAll(slug, ".", "-")
	slug = strings.ReplaceAll(slug, "+", "-")
	if len(slug) > 70 {
		slug = slug[:70]
	}
	// Add random suffix to guarantee uniqueness
	suffix := strings.ReplaceAll(uuid.New().String(), "-", "")[:6]
	return slug + "-" + suffix
}

// SeedTenantRolePermissionsPublic exposes seedTenantRolePermissions for use by the tenants module
// when creating new tenants after initial migration.
func SeedTenantRolePermissionsPublic(db *gorm.DB, tenantID string, logger zerolog.Logger) {
	seedTenantRolePermissions(db, tenantID, logger)
}

// GetRoleDefaultPermissions returns a copy of the default permission map.
// Used by the tenants module when creating new tenants.
func GetRoleDefaultPermissions() map[string][]string {
	result := make(map[string][]string, len(roleDefaultPermissions))
	for role, perms := range roleDefaultPermissions {
		cp := make([]string, len(perms))
		copy(cp, perms)
		result[role] = cp
	}
	return result
}

// GetCurrentTimestamp returns the current UTC time. Helper for migration scripts.
func GetCurrentTimestamp() time.Time {
	return time.Now().UTC()
}
