package ports

import (
	"context"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
)

// TenantRepository defines persistence operations for tenants.
type TenantRepository interface {
	FindTenantByID(ctx context.Context, id string) (*domain.Tenant, error)
	FindTenantForUser(ctx context.Context, userID string) (*domain.Tenant, error)
	FindTenantsByUserID(ctx context.Context, userID string) ([]domain.TenantWithRole, error)
	UpdateTenant(ctx context.Context, id string, req domain.UpdateTenantRequest) error
	DeleteTenant(ctx context.Context, id string) error
}

// MemberRepository defines persistence operations for tenant members.
type MemberRepository interface {
	ListMembers(ctx context.Context, tenantID string) ([]domain.TenantMember, error)
	FindMember(ctx context.Context, tenantID, userID string) (*domain.TenantMember, error)
	AddMember(ctx context.Context, tenantID, userID, role string, invitedBy *string) (string, error)
	UpdateMemberRole(ctx context.Context, tenantID, userID, role string) error
	RemoveMember(ctx context.Context, tenantID, userID string) error
}

// InvitationRepository defines persistence operations for tenant invitations.
type InvitationRepository interface {
	CreateInvitation(ctx context.Context, inv domain.Invitation) error
	FindInvitationByCode(ctx context.Context, code string) (*domain.Invitation, error)
	ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error)
	IncrementInvitationUsed(ctx context.Context, code string) error
	RevokeInvitation(ctx context.Context, tenantID, code string) error
}

// PermissionRepository defines read operations for RBAC permissions.
type PermissionRepository interface {
	ListPermissionsByRole(ctx context.Context, tenantID, role string) ([]string, error)
}

// AuditRepository defines persistence operations for audit logs.
type AuditRepository interface {
	CreateAuditLog(ctx context.Context, log domain.AuditLog) error
	ListAuditLogs(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditLog, error)
}

// TenantService defines all business operations for the tenants module.
type TenantService interface {
	GetMyTenant(ctx context.Context, tenantID string) (*domain.Tenant, error)
	ListMyTenants(ctx context.Context, userID string) ([]domain.TenantWithRole, error)
	UpdateMyTenant(ctx context.Context, tenantID, userID string, req domain.UpdateTenantRequest) error
	DeleteMyTenant(ctx context.Context, tenantID, userID string) error

	ListMembers(ctx context.Context, tenantID string) ([]domain.TenantMember, error)
	UpdateMemberRole(ctx context.Context, tenantID, requesterRole, targetUserID, newRole string) error
	RemoveMember(ctx context.Context, tenantID, requesterRole, requesterID, targetUserID string) error

	CreateInvitation(ctx context.Context, tenantID, userID string, req domain.CreateInvitationRequest) (*domain.Invitation, error)
	ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error)
	RevokeInvitation(ctx context.Context, tenantID, code string) error
	JoinTenant(ctx context.Context, userID, code string) (*domain.Tenant, error)

	ListAuditLogs(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditLog, error)
	GetMyPermissions(ctx context.Context, tenantID, role string) ([]string, error)
}
