package domain

import "time"

// TenantMember represents a user's membership in a tenant.
type TenantMember struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id"`
	UserEmail string    `json:"user_email,omitempty"` // populated via JOIN on read
	UserName  string    `json:"user_name,omitempty"`  // populated via JOIN on read
	Role      string    `json:"role"`
	InvitedBy *string   `json:"invited_by,omitempty"`
	JoinedAt  time.Time `json:"joined_at"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateMemberRoleRequest holds the new role to assign.
type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member viewer"`
}
