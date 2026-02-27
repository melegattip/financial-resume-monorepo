package domain

import "time"

// TenantMember represents a user's membership in a tenant.
type TenantMember struct {
	ID        string
	TenantID  string
	UserID    string
	Role      string
	InvitedBy *string
	JoinedAt  time.Time
	CreatedAt time.Time
}

// UpdateMemberRoleRequest holds the new role to assign.
type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member viewer"`
}
