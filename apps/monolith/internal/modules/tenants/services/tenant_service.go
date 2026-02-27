package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/ports"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/repository"
)

// tenantService implements ports.TenantService.
type tenantService struct {
	tenantRepo ports.TenantRepository
	memberRepo ports.MemberRepository
	invRepo    ports.InvitationRepository
	permRepo   ports.PermissionRepository
	auditRepo  ports.AuditRepository
	logger     zerolog.Logger
}

// NewTenantService creates a new TenantService, wiring the GormRepository
// as the implementation for all repository interfaces.
func NewTenantService(repo *repository.GormRepository, logger zerolog.Logger) ports.TenantService {
	return &tenantService{
		tenantRepo: repo,
		memberRepo: repo,
		invRepo:    repo,
		permRepo:   repo,
		auditRepo:  repo,
		logger:     logger,
	}
}

// ─── Tenant ───────────────────────────────────────────────────────────────────

// GetMyTenant returns the tenant of the calling user.
func (s *tenantService) GetMyTenant(ctx context.Context, userID string) (*domain.Tenant, error) {
	tenant, err := s.tenantRepo.FindTenantForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, errors.New("tenant not found")
	}
	return tenant, nil
}

// UpdateMyTenant updates the tenant name or settings.
// The caller must already have manage_tenant permission (enforced by middleware).
func (s *tenantService) UpdateMyTenant(ctx context.Context, tenantID, _ string, req domain.UpdateTenantRequest) error {
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if len(trimmed) == 0 {
			return errors.New("tenant name cannot be empty")
		}
		req.Name = &trimmed
	}
	return s.tenantRepo.UpdateTenant(ctx, tenantID, req)
}

// DeleteMyTenant soft-deletes the tenant.
// Only the owner can delete; enforced by the delete_tenant permission middleware.
func (s *tenantService) DeleteMyTenant(ctx context.Context, tenantID, _ string) error {
	return s.tenantRepo.DeleteTenant(ctx, tenantID)
}

// ─── Members ─────────────────────────────────────────────────────────────────

// ListMembers returns all members of the tenant.
func (s *tenantService) ListMembers(ctx context.Context, tenantID string) ([]domain.TenantMember, error) {
	return s.memberRepo.ListMembers(ctx, tenantID)
}

// UpdateMemberRole changes a member's role.
// Callers cannot change the owner's role or promote someone to owner.
func (s *tenantService) UpdateMemberRole(ctx context.Context, tenantID, requesterRole, targetUserID, newRole string) error {
	if newRole == "owner" {
		return errors.New("cannot assign owner role; use transfer_ownership instead")
	}

	target, err := s.memberRepo.FindMember(ctx, tenantID, targetUserID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("member not found")
	}
	if target.Role == "owner" {
		return errors.New("cannot change owner's role")
	}

	return s.memberRepo.UpdateMemberRole(ctx, tenantID, targetUserID, newRole)
}

// RemoveMember removes a user from the tenant.
// Cannot remove the owner; a member cannot remove themselves (use leave instead).
func (s *tenantService) RemoveMember(ctx context.Context, tenantID, requesterRole, requesterID, targetUserID string) error {
	if requesterID == targetUserID {
		return errors.New("cannot remove yourself; leave the tenant instead")
	}

	target, err := s.memberRepo.FindMember(ctx, tenantID, targetUserID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("member not found")
	}
	if target.Role == "owner" {
		return errors.New("cannot remove the tenant owner")
	}

	return s.memberRepo.RemoveMember(ctx, tenantID, targetUserID)
}

// ─── Invitations ─────────────────────────────────────────────────────────────

// CreateInvitation generates a new invitation code for the tenant.
func (s *tenantService) CreateInvitation(ctx context.Context, tenantID, userID string, req domain.CreateInvitationRequest) (*domain.Invitation, error) {
	maxUses := req.MaxUses
	if maxUses <= 0 {
		maxUses = 10 // default
	}

	inv := domain.Invitation{
		ID:        "inv_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:8],
		TenantID:  tenantID,
		Code:      domain.GenerateInviteCode(),
		Role:      req.Role,
		CreatedBy: userID,
		ExpiresAt: req.ExpiresAt,
		MaxUses:   maxUses,
		UsedCount: 0,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.invRepo.CreateInvitation(ctx, inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

// ListInvitations returns all active invitations for the tenant.
func (s *tenantService) ListInvitations(ctx context.Context, tenantID string) ([]domain.Invitation, error) {
	return s.invRepo.ListInvitations(ctx, tenantID)
}

// RevokeInvitation deactivates an invitation so it can no longer be used.
func (s *tenantService) RevokeInvitation(ctx context.Context, tenantID, code string) error {
	return s.invRepo.RevokeInvitation(ctx, tenantID, code)
}

// JoinTenant adds the calling user to a tenant using a valid invitation code.
// After joining, the user must re-login to receive updated JWT claims.
func (s *tenantService) JoinTenant(ctx context.Context, userID, code string) (*domain.Tenant, error) {
	inv, err := s.invRepo.FindInvitationByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if inv == nil || !inv.IsUsable() {
		return nil, errors.New("invitation is invalid or expired")
	}

	// Check if user is already a member
	existing, err := s.memberRepo.FindMember(ctx, inv.TenantID, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("you are already a member of this tenant")
	}

	invitedBy := inv.CreatedBy
	if _, err := s.memberRepo.AddMember(ctx, inv.TenantID, userID, inv.Role, &invitedBy); err != nil {
		return nil, err
	}

	if err := s.invRepo.IncrementInvitationUsed(ctx, code); err != nil {
		s.logger.Warn().Err(err).Str("code", code).Msg("failed to increment invitation used count")
	}

	tenant, err := s.tenantRepo.FindTenantByID(ctx, inv.TenantID)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

// ─── Audit & Permissions ─────────────────────────────────────────────────────

// ListAuditLogs returns paginated audit log entries for the tenant.
func (s *tenantService) ListAuditLogs(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditLog, error) {
	return s.auditRepo.ListAuditLogs(ctx, tenantID, limit, offset)
}

// GetMyPermissions returns the permission keys assigned to the caller's role.
func (s *tenantService) GetMyPermissions(ctx context.Context, tenantID, role string) ([]string, error) {
	return s.permRepo.ListPermissionsByRole(ctx, tenantID, role)
}
