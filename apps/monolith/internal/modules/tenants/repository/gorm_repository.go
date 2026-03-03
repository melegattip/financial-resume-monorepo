package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
)

// GormRepository provides all data access for the tenants module.
// It implements the following port interfaces:
//   - ports.TenantRepository
//   - ports.MemberRepository
//   - ports.InvitationRepository
//   - ports.PermissionRepository
//   - ports.AuditRepository
//   - auth/ports.TenantCreator
//   - auth/ports.TenantMemberFinder
//   - middleware.PermissionChecker
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new tenants GORM repository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// ─── Auth-facing interfaces ───────────────────────────────────────────────────

// CreatePersonalTenant creates a personal tenant for a newly registered user,
// registers them as owner, and seeds default role permissions.
// Implements auth/ports.TenantCreator.
func (r *GormRepository) CreatePersonalTenant(ctx context.Context, userID, email string) (string, error) {
	tenantID := "tnt_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8]

	slug := strings.ToLower(email)
	slug = strings.ReplaceAll(slug, "@", "-at-")
	slug = strings.ReplaceAll(slug, ".", "-")
	if len(slug) > 70 {
		slug = slug[:70]
	}
	slug = slug + "-" + strings.ReplaceAll(uuid.New().String(), "-", "")[:6]

	name := email + " (personal)"

	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO tenants (id, name, slug, owner_id, is_active, plan, created_at, updated_at)
		VALUES (?, ?, ?, ?, true, 'free', NOW(), NOW())
	`, tenantID, name, slug, userID).Error; err != nil {
		return "", fmt.Errorf("failed to create tenant: %w", err)
	}

	memberID := "tmb_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO tenant_members (id, tenant_id, user_id, role, joined_at, created_at)
		VALUES (?, ?, ?, 'owner', NOW(), NOW())
		ON CONFLICT (tenant_id, user_id) DO NOTHING
	`, memberID, tenantID, userID).Error; err != nil {
		return "", fmt.Errorf("failed to create tenant membership: %w", err)
	}

	r.seedDefaultRolePermissions(ctx, tenantID)

	return tenantID, nil
}

// FindTenantByUserID returns the tenant ID and role for a given user.
// Returns the most recently joined tenant when the user has multiple memberships.
// Implements auth/ports.TenantMemberFinder.
func (r *GormRepository) FindTenantByUserID(ctx context.Context, userID string) (tenantID, role string, err error) {
	var member TenantMemberModel
	if dbErr := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("joined_at ASC").
		First(&member).Error; dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			return "", "owner", nil
		}
		return "", "", fmt.Errorf("failed to find tenant membership: %w", dbErr)
	}
	return member.TenantID, member.Role, nil
}

// FindMemberInTenant returns the role of a user within a specific tenant.
// Returns an empty string and no error when the user is not a member.
// Implements auth/ports.TenantMemberFinder (extended).
func (r *GormRepository) FindMemberInTenant(ctx context.Context, userID, tenantID string) (role string, err error) {
	var member TenantMemberModel
	if dbErr := r.db.WithContext(ctx).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		First(&member).Error; dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to find member in tenant: %w", dbErr)
	}
	return member.Role, nil
}

// FindTenantsByUserID returns all tenants the user is a member of, with their role.
func (r *GormRepository) FindTenantsByUserID(ctx context.Context, userID string) ([]domain.TenantWithRole, error) {
	type row struct {
		TenantID string    `gorm:"column:tenant_id"`
		Name     string    `gorm:"column:name"`
		Slug     string    `gorm:"column:slug"`
		OwnerID  string    `gorm:"column:owner_id"`
		IsActive bool      `gorm:"column:is_active"`
		Plan     string    `gorm:"column:plan"`
		Role     string    `gorm:"column:role"`
		JoinedAt time.Time `gorm:"column:joined_at"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT tm.tenant_id, t.name, t.slug, t.owner_id, t.is_active, t.plan,
		       tm.role, tm.joined_at
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id AND t.deleted_at IS NULL
		WHERE tm.user_id = ?
		ORDER BY tm.joined_at ASC
	`, userID).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("find tenants by user: %w", err)
	}
	result := make([]domain.TenantWithRole, len(rows))
	for i, r := range rows {
		result[i] = domain.TenantWithRole{
			ID:       r.TenantID,
			Name:     r.Name,
			Slug:     r.Slug,
			OwnerID:  r.OwnerID,
			IsActive: r.IsActive,
			Plan:     r.Plan,
			Role:     r.Role,
			JoinedAt: r.JoinedAt,
		}
	}
	return result, nil
}

// HasPermission checks if a role has a specific permission within a tenant.
// Implements middleware.PermissionChecker.
func (r *GormRepository) HasPermission(ctx context.Context, tenantID, role, permission string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&RolePermissionModel{}).
		Where("tenant_id = ? AND role = ? AND permission_key = ?", tenantID, role, permission).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return count > 0, nil
}

// ─── TenantRepository ─────────────────────────────────────────────────────────

// FindTenantByID returns a tenant by its ID.
func (r *GormRepository) FindTenantByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var m TenantModel
	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find tenant by id: %w", err)
	}
	return tenantModelToDomain(&m), nil
}

// FindTenantForUser returns the tenant of which the given user is a member.
func (r *GormRepository) FindTenantForUser(ctx context.Context, userID string) (*domain.Tenant, error) {
	var member TenantMemberModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find member for user: %w", err)
	}
	return r.FindTenantByID(ctx, member.TenantID)
}

// UpdateTenant updates the mutable fields of a tenant.
func (r *GormRepository) UpdateTenant(ctx context.Context, id string, req domain.UpdateTenantRequest) error {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Settings != nil {
		updates["settings"] = *req.Settings
	}
	return r.db.WithContext(ctx).
		Model(&TenantModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

// DeleteTenant soft-deletes a tenant.
func (r *GormRepository) DeleteTenant(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&TenantModel{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// ─── MemberRepository ─────────────────────────────────────────────────────────

// ListMembers returns all members of a tenant ordered by join date,
// enriched with the user's email and display name via a LEFT JOIN on users.
func (r *GormRepository) ListMembers(ctx context.Context, tenantID string) ([]domain.TenantMember, error) {
	type memberRow struct {
		ID        string    `gorm:"column:id"`
		TenantID  string    `gorm:"column:tenant_id"`
		UserID    string    `gorm:"column:user_id"`
		UserEmail string    `gorm:"column:user_email"`
		UserName  string    `gorm:"column:user_name"`
		Role      string    `gorm:"column:role"`
		InvitedBy *string   `gorm:"column:invited_by"`
		JoinedAt  time.Time `gorm:"column:joined_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}
	var rows []memberRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT tm.id, tm.tenant_id, tm.user_id,
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.first_name || ' ' || u.last_name), ''), u.email) AS user_name,
		       tm.role, tm.invited_by, tm.joined_at, tm.created_at
		FROM tenant_members tm
		LEFT JOIN users u ON u.id = tm.user_id
		WHERE tm.tenant_id = ?
		ORDER BY tm.joined_at ASC
	`, tenantID).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	members := make([]domain.TenantMember, len(rows))
	for i, row := range rows {
		members[i] = domain.TenantMember{
			ID:        row.ID,
			TenantID:  row.TenantID,
			UserID:    row.UserID,
			UserEmail: row.UserEmail,
			UserName:  row.UserName,
			Role:      row.Role,
			InvitedBy: row.InvitedBy,
			JoinedAt:  row.JoinedAt,
			CreatedAt: row.CreatedAt,
		}
	}
	return members, nil
}

// FindMember returns a specific member of a tenant, or nil if not found.
func (r *GormRepository) FindMember(ctx context.Context, tenantID, userID string) (*domain.TenantMember, error) {
	var m TenantMemberModel
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find member: %w", err)
	}
	result := memberModelToDomain(&m)
	return &result, nil
}

// AddMember inserts a new tenant membership.
func (r *GormRepository) AddMember(ctx context.Context, tenantID, userID, role string, invitedBy *string) (string, error) {
	id := "tmb_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
	m := TenantMemberModel{
		ID:        id,
		TenantID:  tenantID,
		UserID:    userID,
		Role:      role,
		InvitedBy: invitedBy,
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return "", fmt.Errorf("add member: %w", err)
	}
	return id, nil
}

// UpdateMemberRole changes the role of an existing tenant member.
func (r *GormRepository) UpdateMemberRole(ctx context.Context, tenantID, userID, role string) error {
	return r.db.WithContext(ctx).
		Model(&TenantMemberModel{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Update("role", role).Error
}

// RemoveMember removes a user from a tenant.
func (r *GormRepository) RemoveMember(ctx context.Context, tenantID, userID string) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Delete(&TenantMemberModel{}).Error
}

// ─── InvitationRepository ─────────────────────────────────────────────────────

// CreateInvitation persists a new tenant invitation.
func (r *GormRepository) CreateInvitation(ctx context.Context, inv domain.Invitation) error {
	m := TenantInvitationModel{
		ID:        inv.ID,
		TenantID:  inv.TenantID,
		Code:      inv.Code,
		Role:      inv.Role,
		CreatedBy: inv.CreatedBy,
		ExpiresAt: inv.ExpiresAt,
		MaxUses:   inv.MaxUses,
		UsedCount: inv.UsedCount,
		IsActive:  inv.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

// FindInvitationByCode returns an invitation by its code, or nil if not found.
func (r *GormRepository) FindInvitationByCode(ctx context.Context, code string) (*domain.Invitation, error) {
	var m TenantInvitationModel
	if err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find invitation by code: %w", err)
	}
	result := invitationModelToDomain(&m)
	return &result, nil
}

// ListInvitations returns all active invitations for a tenant, newest first.
func (r *GormRepository) ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error) {
	var models []TenantInvitationModel
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = true", tenantID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list invitations: %w", err)
	}
	invs := make([]domain.Invitation, len(models))
	for i, m := range models {
		invs[i] = invitationModelToDomain(&m)
	}
	return invs, nil
}

// IncrementInvitationUsed atomically increments the used_count of an invitation.
func (r *GormRepository) IncrementInvitationUsed(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).
		Model(&TenantInvitationModel{}).
		Where("code = ?", code).
		Updates(map[string]interface{}{
			"used_count": gorm.Expr("used_count + 1"),
			"updated_at": time.Now(),
		}).Error
}

// RevokeInvitation marks an invitation as inactive.
func (r *GormRepository) RevokeInvitation(ctx context.Context, tenantID, code string) error {
	return r.db.WithContext(ctx).
		Model(&TenantInvitationModel{}).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

// ─── PermissionRepository ─────────────────────────────────────────────────────

// ListPermissionsByRole returns all permission keys for a role in a tenant.
func (r *GormRepository) ListPermissionsByRole(ctx context.Context, tenantID, role string) ([]string, error) {
	var models []RolePermissionModel
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND role = ?", tenantID, role).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list permissions by role: %w", err)
	}
	perms := make([]string, len(models))
	for i, m := range models {
		perms[i] = m.PermissionKey
	}
	return perms, nil
}

// ─── AuditRepository ─────────────────────────────────────────────────────────

// CreateAuditLog persists a new audit log entry.
func (r *GormRepository) CreateAuditLog(ctx context.Context, log domain.AuditLog) error {
	oldValues := log.OldValues
	if oldValues == "" {
		oldValues = "null"
	}
	newValues := log.NewValues
	if newValues == "" {
		newValues = "null"
	}
	m := AuditLogModel{
		ID:          log.ID,
		TenantID:    log.TenantID,
		UserID:      log.UserID,
		Action:      log.Action,
		Description: log.Description,
		EntityType:  log.EntityType,
		EntityID:    log.EntityID,
		OldValues:   oldValues,
		NewValues:   newValues,
		IPAddress:   log.IPAddress,
		UserAgent:   log.UserAgent,
		CreatedAt:   time.Now(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

// ListAuditLogs returns audit log entries for a tenant, newest first.
// Descriptions are computed on-the-fly via LEFT JOINs on entity tables, so that
// both old records (no stored description) and new records show context.
func (r *GormRepository) ListAuditLogs(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	type auditLogRow struct {
		ID          string    `gorm:"column:id"`
		TenantID    string    `gorm:"column:tenant_id"`
		UserID      string    `gorm:"column:user_id"`
		UserName    string    `gorm:"column:user_name"`
		Action      string    `gorm:"column:action"`
		Description string    `gorm:"column:description"`
		EntityType  string    `gorm:"column:entity_type"`
		EntityID    string    `gorm:"column:entity_id"`
		OldValues   string    `gorm:"column:old_values"`
		NewValues   string    `gorm:"column:new_values"`
		IPAddress   string    `gorm:"column:ip_address"`
		UserAgent   string    `gorm:"column:user_agent"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}
	var rows []auditLogRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT al.id, al.tenant_id, al.user_id,
		       COALESCE(NULLIF(TRIM(u.first_name || ' ' || u.last_name), ''), u.email) AS user_name,
		       al.action,
		       CASE
		         WHEN al.entity_type = 'expense' AND e.id IS NOT NULL THEN
		           e.description || ' · $' || TO_CHAR(e.amount, 'FM9999999990.00')
		         WHEN al.entity_type = 'income' AND i.id IS NOT NULL THEN
		           i.source || ' · $' || TO_CHAR(i.amount, 'FM9999999990.00')
		         WHEN al.entity_type = 'recurring_transaction' AND r.id IS NOT NULL THEN
		           r.description || ' · $' || TO_CHAR(r.amount, 'FM9999999990.00')
		         WHEN al.entity_type = 'budget' AND b.id IS NOT NULL THEN
		           '$' || TO_CHAR(b.amount, 'FM9999999990.00') || ' (' || b.period || ')'
		         WHEN al.entity_type = 'savings_goal' AND sg.id IS NOT NULL THEN
		           sg.name || ' · $' || TO_CHAR(sg.target_amount, 'FM9999999990.00')
		         ELSE NULL
		       END AS description,
		       al.entity_type, al.entity_id,
		       al.old_values, al.new_values, al.ip_address, al.user_agent, al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.user_id
		LEFT JOIN expenses e ON al.entity_type = 'expense' AND e.id = al.entity_id
		LEFT JOIN incomes i ON al.entity_type = 'income' AND i.id = al.entity_id
		LEFT JOIN recurring_transactions r ON al.entity_type = 'recurring_transaction' AND r.id = al.entity_id
		LEFT JOIN budgets b ON al.entity_type = 'budget' AND b.id = al.entity_id
		LEFT JOIN savings_goals sg ON al.entity_type = 'savings_goal' AND sg.id = al.entity_id
		WHERE al.tenant_id = ?
		ORDER BY al.created_at DESC
		LIMIT ? OFFSET ?
	`, tenantID, limit, offset).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	logs := make([]domain.AuditLog, len(rows))
	for i, row := range rows {
		logs[i] = domain.AuditLog{
			ID:          row.ID,
			TenantID:    row.TenantID,
			UserID:      row.UserID,
			UserName:    row.UserName,
			Action:      row.Action,
			Description: row.Description,
			EntityType:  row.EntityType,
			EntityID:    row.EntityID,
			OldValues:   row.OldValues,
			NewValues:   row.NewValues,
			IPAddress:   row.IPAddress,
			UserAgent:   row.UserAgent,
			CreatedAt:   row.CreatedAt,
		}
	}
	return logs, nil
}

// ─── Seeding ─────────────────────────────────────────────────────────────────

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

func (r *GormRepository) seedDefaultRolePermissions(ctx context.Context, tenantID string) {
	for role, perms := range roleDefaultPermissions {
		for _, perm := range perms {
			id := "rp_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
			_ = r.db.WithContext(ctx).Exec(`
				INSERT INTO role_permissions (id, tenant_id, role, permission_key, created_at)
				VALUES (?, ?, ?, ?, NOW())
				ON CONFLICT (tenant_id, role, permission_key) DO NOTHING
			`, id, tenantID, role, perm).Error
		}
	}
}

// ─── Mappers ─────────────────────────────────────────────────────────────────

func tenantModelToDomain(m *TenantModel) *domain.Tenant {
	return &domain.Tenant{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		OwnerID:   m.OwnerID,
		IsActive:  m.IsActive,
		Plan:      m.Plan,
		Settings:  m.Settings,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

func memberModelToDomain(m *TenantMemberModel) domain.TenantMember {
	return domain.TenantMember{
		ID:        m.ID,
		TenantID:  m.TenantID,
		UserID:    m.UserID,
		Role:      m.Role,
		InvitedBy: m.InvitedBy,
		JoinedAt:  m.JoinedAt,
		CreatedAt: m.CreatedAt,
	}
}

func invitationModelToDomain(m *TenantInvitationModel) domain.Invitation {
	return domain.Invitation{
		ID:        m.ID,
		TenantID:  m.TenantID,
		Code:      m.Code,
		Role:      m.Role,
		CreatedBy: m.CreatedBy,
		ExpiresAt: m.ExpiresAt,
		MaxUses:   m.MaxUses,
		UsedCount: m.UsedCount,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func auditModelToDomain(m *AuditLogModel) domain.AuditLog {
	return domain.AuditLog{
		ID:         m.ID,
		TenantID:   m.TenantID,
		UserID:     m.UserID,
		Action:     m.Action,
		EntityType: m.EntityType,
		EntityID:   m.EntityID,
		OldValues:  m.OldValues,
		NewValues:  m.NewValues,
		IPAddress:  m.IPAddress,
		UserAgent:  m.UserAgent,
		CreatedAt:  m.CreatedAt,
	}
}
