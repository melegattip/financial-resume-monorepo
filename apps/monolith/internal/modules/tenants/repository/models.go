package repository

import "time"

// TenantModel is the GORM model for the tenants table.
type TenantModel struct {
	ID        string     `gorm:"column:id;type:varchar(50);primaryKey"`
	Name      string     `gorm:"column:name;type:varchar(255);not null"`
	Slug      string     `gorm:"column:slug;type:varchar(100);uniqueIndex;not null"`
	OwnerID   string     `gorm:"column:owner_id;type:varchar(255);not null;index"`
	IsActive  bool       `gorm:"column:is_active;not null;default:true"`
	Plan      string     `gorm:"column:plan;type:varchar(20);not null;default:free"`
	Settings  string     `gorm:"column:settings;type:jsonb"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (TenantModel) TableName() string { return "tenants" }

// TenantMemberModel is the GORM model for the tenant_members table.
type TenantMemberModel struct {
	ID        string    `gorm:"column:id;type:varchar(50);primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;type:varchar(50);not null;uniqueIndex:idx_unique_membership"`
	UserID    string    `gorm:"column:user_id;type:varchar(255);not null;uniqueIndex:idx_unique_membership"`
	Role      string    `gorm:"column:role;type:varchar(20);not null;default:member"`
	InvitedBy *string   `gorm:"column:invited_by;type:varchar(255)"`
	JoinedAt  time.Time `gorm:"column:joined_at;autoCreateTime"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (TenantMemberModel) TableName() string { return "tenant_members" }

// PermissionModel is the GORM model for the permissions catalog table.
type PermissionModel struct {
	Key         string `gorm:"column:key;type:varchar(100);primaryKey"`
	Description string `gorm:"column:description;type:text"`
	Category    string `gorm:"column:category;type:varchar(50)"`
}

func (PermissionModel) TableName() string { return "permissions" }

// RolePermissionModel is the GORM model for the role_permissions table.
type RolePermissionModel struct {
	ID            string    `gorm:"column:id;type:varchar(50);primaryKey"`
	TenantID      string    `gorm:"column:tenant_id;type:varchar(50);not null;uniqueIndex:idx_unique_role_perm"`
	Role          string    `gorm:"column:role;type:varchar(20);not null;uniqueIndex:idx_unique_role_perm"`
	PermissionKey string    `gorm:"column:permission_key;type:varchar(100);not null;uniqueIndex:idx_unique_role_perm"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (RolePermissionModel) TableName() string { return "role_permissions" }

// TenantInvitationModel is the GORM model for the tenant_invitations table.
type TenantInvitationModel struct {
	ID        string     `gorm:"column:id;type:varchar(50);primaryKey"`
	TenantID  string     `gorm:"column:tenant_id;type:varchar(50);not null;index"`
	Code      string     `gorm:"column:code;type:varchar(20);uniqueIndex;not null"`
	Role      string     `gorm:"column:role;type:varchar(20);not null;default:member"`
	CreatedBy string     `gorm:"column:created_by;type:varchar(255);not null"`
	ExpiresAt *time.Time `gorm:"column:expires_at"`
	MaxUses   int        `gorm:"column:max_uses;not null;default:10"`
	UsedCount int        `gorm:"column:used_count;not null;default:0"`
	IsActive  bool       `gorm:"column:is_active;not null;default:true"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (TenantInvitationModel) TableName() string { return "tenant_invitations" }

// AuditLogModel is the GORM model for the audit_logs table.
// No soft delete — audit logs are immutable.
type AuditLogModel struct {
	ID         string    `gorm:"column:id;type:varchar(50);primaryKey"`
	TenantID   string    `gorm:"column:tenant_id;type:varchar(50);not null;index:idx_audit_tenant_time"`
	UserID     string    `gorm:"column:user_id;type:varchar(255);not null"`
	Action     string    `gorm:"column:action;type:varchar(50);not null"`
	EntityType string    `gorm:"column:entity_type;type:varchar(50)"`
	EntityID   string    `gorm:"column:entity_id;type:varchar(255)"`
	OldValues  string    `gorm:"column:old_values;type:jsonb"`
	NewValues  string    `gorm:"column:new_values;type:jsonb"`
	IPAddress  string    `gorm:"column:ip_address;type:varchar(45)"`
	UserAgent  string    `gorm:"column:user_agent;type:varchar(500)"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_audit_tenant_time"`
}

func (AuditLogModel) TableName() string { return "audit_logs" }
